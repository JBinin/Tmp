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
	"net"

	"github.com/hashicorp/memberlist"
)

type MergeDelegate interface {
	NotifyMerge([]*Member) error
}

type mergeDelegate struct {
	serf *Serf
}

func (m *mergeDelegate) NotifyMerge(nodes []*memberlist.Node) error {
	members := make([]*Member, len(nodes))
	for idx, n := range nodes {
		members[idx] = m.nodeToMember(n)
	}
	return m.serf.config.Merge.NotifyMerge(members)
}

func (m *mergeDelegate) NotifyAlive(peer *memberlist.Node) error {
	member := m.nodeToMember(peer)
	return m.serf.config.Merge.NotifyMerge([]*Member{member})
}

func (m *mergeDelegate) nodeToMember(n *memberlist.Node) *Member {
	return &Member{
		Name:        n.Name,
		Addr:        net.IP(n.Addr),
		Port:        n.Port,
		Tags:        m.serf.decodeTags(n.Meta),
		Status:      StatusNone,
		ProtocolMin: n.PMin,
		ProtocolMax: n.PMax,
		ProtocolCur: n.PCur,
		DelegateMin: n.DMin,
		DelegateMax: n.DMax,
		DelegateCur: n.DCur,
	}
}
