/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package cluster

import (
	"errors"
	"fmt"
	"net"
)

var (
	errNoSuchInterface         = errors.New("no such interface")
	errNoIP                    = errors.New("could not find the system's IP address")
	errMustSpecifyListenAddr   = errors.New("must specify a listening address because the address to advertise is not recognized as a system address, and a system's IP address to use could not be uniquely identified")
	errBadListenAddr           = errors.New("listen address must be an IP address or network interface (with optional port number)")
	errBadAdvertiseAddr        = errors.New("advertise address must be a non-zero IP address or network interface (with optional port number)")
	errBadDefaultAdvertiseAddr = errors.New("default advertise address must be a non-zero IP address or network interface (without a port number)")
)

func resolveListenAddr(specifiedAddr string) (string, string, error) {
	specifiedHost, specifiedPort, err := net.SplitHostPort(specifiedAddr)
	if err != nil {
		return "", "", fmt.Errorf("could not parse listen address %s", specifiedAddr)
	}

	// Does the host component match any of the interface names on the
	// system? If so, use the address from that interface.
	interfaceAddr, err := resolveInterfaceAddr(specifiedHost)
	if err == nil {
		return interfaceAddr.String(), specifiedPort, nil
	}
	if err != errNoSuchInterface {
		return "", "", err
	}

	// If it's not an interface, it must be an IP (for now)
	if net.ParseIP(specifiedHost) == nil {
		return "", "", errBadListenAddr
	}

	return specifiedHost, specifiedPort, nil
}

func (c *Cluster) resolveAdvertiseAddr(advertiseAddr, listenAddrPort string) (string, string, error) {
	// Approach:
	// - If an advertise address is specified, use that. Resolve the
	//   interface's address if an interface was specified in
	//   advertiseAddr. Fill in the port from listenAddrPort if necessary.
	// - If DefaultAdvertiseAddr is not empty, use that with the port from
	//   listenAddrPort. Resolve the interface's address from
	//   if an interface name was specified in DefaultAdvertiseAddr.
	// - Otherwise, try to autodetect the system's address. Use the port in
	//   listenAddrPort with this address if autodetection succeeds.

	if advertiseAddr != "" {
		advertiseHost, advertisePort, err := net.SplitHostPort(advertiseAddr)
		if err != nil {
			// Not a host:port specification
			advertiseHost = advertiseAddr
			advertisePort = listenAddrPort
		}

		// Does the host component match any of the interface names on the
		// system? If so, use the address from that interface.
		interfaceAddr, err := resolveInterfaceAddr(advertiseHost)
		if err == nil {
			return interfaceAddr.String(), advertisePort, nil
		}
		if err != errNoSuchInterface {
			return "", "", err
		}

		// If it's not an interface, it must be an IP (for now)
		if ip := net.ParseIP(advertiseHost); ip == nil || ip.IsUnspecified() {
			return "", "", errBadAdvertiseAddr
		}

		return advertiseHost, advertisePort, nil
	}

	if c.config.DefaultAdvertiseAddr != "" {
		// Does the default advertise address component match any of the
		// interface names on the system? If so, use the address from
		// that interface.
		interfaceAddr, err := resolveInterfaceAddr(c.config.DefaultAdvertiseAddr)
		if err == nil {
			return interfaceAddr.String(), listenAddrPort, nil
		}
		if err != errNoSuchInterface {
			return "", "", err
		}

		// If it's not an interface, it must be an IP (for now)
		if ip := net.ParseIP(c.config.DefaultAdvertiseAddr); ip == nil || ip.IsUnspecified() {
			return "", "", errBadDefaultAdvertiseAddr
		}

		return c.config.DefaultAdvertiseAddr, listenAddrPort, nil
	}

	systemAddr, err := c.resolveSystemAddr()
	if err != nil {
		return "", "", err
	}
	return systemAddr.String(), listenAddrPort, nil
}

func resolveInterfaceAddr(specifiedInterface string) (net.IP, error) {
	// Use a specific interface's IP address.
	intf, err := net.InterfaceByName(specifiedInterface)
	if err != nil {
		return nil, errNoSuchInterface
	}

	addrs, err := intf.Addrs()
	if err != nil {
		return nil, err
	}

	var interfaceAddr4, interfaceAddr6 net.IP

	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)

		if ok {
			if ipAddr.IP.To4() != nil {
				// IPv4
				if interfaceAddr4 != nil {
					return nil, fmt.Errorf("interface %s has more than one IPv4 address (%s and %s)", specifiedInterface, interfaceAddr4, ipAddr.IP)
				}
				interfaceAddr4 = ipAddr.IP
			} else {
				// IPv6
				if interfaceAddr6 != nil {
					return nil, fmt.Errorf("interface %s has more than one IPv6 address (%s and %s)", specifiedInterface, interfaceAddr6, ipAddr.IP)
				}
				interfaceAddr6 = ipAddr.IP
			}
		}
	}

	if interfaceAddr4 == nil && interfaceAddr6 == nil {
		return nil, fmt.Errorf("interface %s has no usable IPv4 or IPv6 address", specifiedInterface)
	}

	// In the case that there's exactly one IPv4 address
	// and exactly one IPv6 address, favor IPv4 over IPv6.
	if interfaceAddr4 != nil {
		return interfaceAddr4, nil
	}
	return interfaceAddr6, nil
}

func (c *Cluster) resolveSystemAddrViaSubnetCheck() (net.IP, error) {
	// Use the system's only IP address, or fail if there are
	// multiple addresses to choose from. Skip interfaces which
	// are managed by docker via subnet check.
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var systemAddr net.IP
	var systemInterface string

	// List Docker-managed subnets
	v4Subnets := c.config.NetworkSubnetsProvider.V4Subnets()
	v6Subnets := c.config.NetworkSubnetsProvider.V6Subnets()

ifaceLoop:
	for _, intf := range interfaces {
		// Skip inactive interfaces and loopback interfaces
		if (intf.Flags&net.FlagUp == 0) || (intf.Flags&net.FlagLoopback) != 0 {
			continue
		}

		addrs, err := intf.Addrs()
		if err != nil {
			continue
		}

		var interfaceAddr4, interfaceAddr6 net.IP

		for _, addr := range addrs {
			ipAddr, ok := addr.(*net.IPNet)

			// Skip loopback and link-local addresses
			if !ok || !ipAddr.IP.IsGlobalUnicast() {
				continue
			}

			if ipAddr.IP.To4() != nil {
				// IPv4

				// Ignore addresses in subnets that are managed by Docker.
				for _, subnet := range v4Subnets {
					if subnet.Contains(ipAddr.IP) {
						continue ifaceLoop
					}
				}

				if interfaceAddr4 != nil {
					return nil, errMultipleIPs(intf.Name, intf.Name, interfaceAddr4, ipAddr.IP)
				}

				interfaceAddr4 = ipAddr.IP
			} else {
				// IPv6

				// Ignore addresses in subnets that are managed by Docker.
				for _, subnet := range v6Subnets {
					if subnet.Contains(ipAddr.IP) {
						continue ifaceLoop
					}
				}

				if interfaceAddr6 != nil {
					return nil, errMultipleIPs(intf.Name, intf.Name, interfaceAddr6, ipAddr.IP)
				}

				interfaceAddr6 = ipAddr.IP
			}
		}

		// In the case that this interface has exactly one IPv4 address
		// and exactly one IPv6 address, favor IPv4 over IPv6.
		if interfaceAddr4 != nil {
			if systemAddr != nil {
				return nil, errMultipleIPs(systemInterface, intf.Name, systemAddr, interfaceAddr4)
			}
			systemAddr = interfaceAddr4
			systemInterface = intf.Name
		} else if interfaceAddr6 != nil {
			if systemAddr != nil {
				return nil, errMultipleIPs(systemInterface, intf.Name, systemAddr, interfaceAddr6)
			}
			systemAddr = interfaceAddr6
			systemInterface = intf.Name
		}
	}

	if systemAddr == nil {
		return nil, errNoIP
	}

	return systemAddr, nil
}

func listSystemIPs() []net.IP {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var systemAddrs []net.IP

	for _, intf := range interfaces {
		addrs, err := intf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipAddr, ok := addr.(*net.IPNet)

			if ok {
				systemAddrs = append(systemAddrs, ipAddr.IP)
			}
		}
	}

	return systemAddrs
}

func errMultipleIPs(interfaceA, interfaceB string, addrA, addrB net.IP) error {
	if interfaceA == interfaceB {
		return fmt.Errorf("could not choose an IP address to advertise since this system has multiple addresses on interface %s (%s and %s)", interfaceA, addrA, addrB)
	}
	return fmt.Errorf("could not choose an IP address to advertise since this system has multiple addresses on different interfaces (%s on %s and %s on %s)", addrA, interfaceA, addrB, interfaceB)
}
