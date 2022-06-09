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
package netutils

import (
	"net"

	"github.com/docker/libnetwork/types"
)

// ElectInterfaceAddresses looks for an interface on the OS with the specified name
// and returns returns all its IPv4 and IPv6 addresses in CIDR notation.
// If a failure in retrieving the addresses or no IPv4 address is found, an error is returned.
// If the interface does not exist, it chooses from a predefined
// list the first IPv4 address which does not conflict with other
// interfaces on the system.
func ElectInterfaceAddresses(name string) ([]*net.IPNet, []*net.IPNet, error) {
	return nil, nil, types.NotImplementedErrorf("not supported on windows")
}

// FindAvailableNetwork returns a network from the passed list which does not
// overlap with existing interfaces in the system

// TODO : Use appropriate windows APIs to identify non-overlapping subnets
func FindAvailableNetwork(list []*net.IPNet) (*net.IPNet, error) {
	return nil, nil
}
