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
package serf

import (
	"github.com/hashicorp/memberlist"
)

type eventDelegate struct {
	serf *Serf
}

func (e *eventDelegate) NotifyJoin(n *memberlist.Node) {
	e.serf.handleNodeJoin(n)
}

func (e *eventDelegate) NotifyLeave(n *memberlist.Node) {
	e.serf.handleNodeLeave(n)
}

func (e *eventDelegate) NotifyUpdate(n *memberlist.Node) {
	e.serf.handleNodeUpdate(n)
}
