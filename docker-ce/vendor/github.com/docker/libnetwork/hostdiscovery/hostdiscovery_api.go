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
package hostdiscovery

import "net"

// JoinCallback provides a callback event for new node joining the cluster
type JoinCallback func(entries []net.IP)

// ActiveCallback provides a callback event for active discovery event
type ActiveCallback func()

// LeaveCallback provides a callback event for node leaving the cluster
type LeaveCallback func(entries []net.IP)

// HostDiscovery primary interface
type HostDiscovery interface {
	//Watch Node join and leave cluster events
	Watch(activeCallback ActiveCallback, joinCallback JoinCallback, leaveCallback LeaveCallback) error
	// StopDiscovery stops the discovery process
	StopDiscovery() error
	// Fetch returns a list of host IPs that are currently discovered
	Fetch() []net.IP
}
