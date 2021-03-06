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
// Package netlink provides a simple library for netlink. Netlink is
// the interface a user-space program in linux uses to communicate with
// the kernel. It can be used to add and remove interfaces, set up ip
// addresses and routes, and confiugre ipsec. Netlink communication
// requires elevated privileges, so in most cases this code needs to
// be run as root. The low level primitives for netlink are contained
// in the nl subpackage. This package attempts to provide a high-level
// interface that is loosly modeled on the iproute2 cli.
package netlink

import "net"

// ParseIPNet parses a string in ip/net format and returns a net.IPNet.
// This is valuable because addresses in netlink are often IPNets and
// ParseCIDR returns an IPNet with the IP part set to the base IP of the
// range.
func ParseIPNet(s string) (*net.IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &net.IPNet{IP: ip, Mask: ipNet.Mask}, nil
}

// NewIPNet generates an IPNet from an ip address using a netmask of 32 or 128.
func NewIPNet(ip net.IP) *net.IPNet {
	if ip.To4() != nil {
		return &net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}
	}
	return &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
}
