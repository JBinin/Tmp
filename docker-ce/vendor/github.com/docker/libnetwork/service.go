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
package libnetwork

import (
	"fmt"
	"net"
	"sync"
)

var (
	// A global monotonic counter to assign firewall marks to
	// services.
	fwMarkCtr   uint32 = 256
	fwMarkCtrMu sync.Mutex
)

type portConfigs []*PortConfig

func (p portConfigs) String() string {
	if len(p) == 0 {
		return ""
	}

	pc := p[0]
	str := fmt.Sprintf("%d:%d/%s", pc.PublishedPort, pc.TargetPort, PortConfig_Protocol_name[int32(pc.Protocol)])
	for _, pc := range p[1:] {
		str = str + fmt.Sprintf(",%d:%d/%s", pc.PublishedPort, pc.TargetPort, PortConfig_Protocol_name[int32(pc.Protocol)])
	}

	return str
}

type serviceKey struct {
	id    string
	ports string
}

type service struct {
	name string // Service Name
	id   string // Service ID

	// Map of loadbalancers for the service one-per attached
	// network. It is keyed with network ID.
	loadBalancers map[string]*loadBalancer

	// List of ingress ports exposed by the service
	ingressPorts portConfigs

	// Service aliases
	aliases []string

	sync.Mutex
}

type loadBalancer struct {
	vip    net.IP
	fwMark uint32

	// Map of backend IPs backing this loadbalancer on this
	// network. It is keyed with endpoint ID.
	backEnds map[string]net.IP

	// Back pointer to service to which the loadbalancer belongs.
	service *service
}
