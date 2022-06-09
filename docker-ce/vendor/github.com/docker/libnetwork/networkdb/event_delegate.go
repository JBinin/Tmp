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
package networkdb

import (
	"encoding/json"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/memberlist"
)

type eventDelegate struct {
	nDB *NetworkDB
}

func (e *eventDelegate) broadcastNodeEvent(addr net.IP, op opType) {
	value, err := json.Marshal(&NodeAddr{addr})
	if err == nil {
		e.nDB.broadcaster.Write(makeEvent(op, NodeTable, "", "", value))
	} else {
		logrus.Errorf("Error marshalling node broadcast event %s", addr.String())
	}
}

func (e *eventDelegate) NotifyJoin(mn *memberlist.Node) {
	e.broadcastNodeEvent(mn.Addr, opCreate)
	e.nDB.Lock()
	// In case the node is rejoining after a failure or leave,
	// wait until an explicit join message arrives before adding
	// it to the nodes just to make sure this is not a stale
	// join. If you don't know about this node add it immediately.
	_, fOk := e.nDB.failedNodes[mn.Name]
	_, lOk := e.nDB.leftNodes[mn.Name]
	if fOk || lOk {
		e.nDB.Unlock()
		return
	}

	e.nDB.nodes[mn.Name] = &node{Node: *mn}
	e.nDB.Unlock()
}

func (e *eventDelegate) NotifyLeave(mn *memberlist.Node) {
	e.broadcastNodeEvent(mn.Addr, opDelete)
	e.nDB.deleteNodeTableEntries(mn.Name)
	e.nDB.deleteNetworkEntriesForNode(mn.Name)
	e.nDB.Lock()
	if n, ok := e.nDB.nodes[mn.Name]; ok {
		delete(e.nDB.nodes, mn.Name)

		n.reapTime = reapInterval
		e.nDB.failedNodes[mn.Name] = n
	}
	e.nDB.Unlock()
}

func (e *eventDelegate) NotifyUpdate(n *memberlist.Node) {
}
