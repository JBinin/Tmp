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
package edge

import . "github.com/onsi/gomega/matchers/support/goraph/node"

type Edge struct {
	Node1 Node
	Node2 Node
}

type EdgeSet []Edge

func (ec EdgeSet) Free(node Node) bool {
	for _, e := range ec {
		if e.Node1 == node || e.Node2 == node {
			return false
		}
	}

	return true
}

func (ec EdgeSet) Contains(edge Edge) bool {
	for _, e := range ec {
		if e == edge {
			return true
		}
	}

	return false
}

func (ec EdgeSet) FindByNodes(node1, node2 Node) (Edge, bool) {
	for _, e := range ec {
		if (e.Node1 == node1 && e.Node2 == node2) || (e.Node1 == node2 && e.Node2 == node1) {
			return e, true
		}
	}

	return Edge{}, false
}

func (ec EdgeSet) SymmetricDifference(ec2 EdgeSet) EdgeSet {
	edgesToInclude := make(map[Edge]bool)

	for _, e := range ec {
		edgesToInclude[e] = true
	}

	for _, e := range ec2 {
		edgesToInclude[e] = !edgesToInclude[e]
	}

	result := EdgeSet{}
	for e, include := range edgesToInclude {
		if include {
			result = append(result, e)
		}
	}

	return result
}
