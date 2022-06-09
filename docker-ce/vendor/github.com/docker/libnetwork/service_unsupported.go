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
// +build !linux,!windows

package libnetwork

import (
	"fmt"
	"net"
)

func (c *controller) cleanupServiceBindings(nid string) {
}

func (c *controller) addServiceBinding(name, sid, nid, eid string, vip net.IP, ingressPorts []*PortConfig, aliases []string, ip net.IP) error {
	return fmt.Errorf("not supported")
}

func (c *controller) rmServiceBinding(name, sid, nid, eid string, vip net.IP, ingressPorts []*PortConfig, aliases []string, ip net.IP) error {
	return fmt.Errorf("not supported")
}

func (sb *sandbox) populateLoadbalancers(ep *endpoint) {
}

func arrangeIngressFilterRule() {
}
