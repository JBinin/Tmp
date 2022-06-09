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
package simple

import (
	"k8s.io/kubernetes/third_party/forked/gonum/graph"
)

// DirectedAcyclicGraph implements graph.Directed using UndirectedGraph,
// which only stores one edge for any node pair.
type DirectedAcyclicGraph struct {
	*UndirectedGraph
}

func NewDirectedAcyclicGraph(self, absent float64) *DirectedAcyclicGraph {
	return &DirectedAcyclicGraph{
		UndirectedGraph: NewUndirectedGraph(self, absent),
	}
}

func (g *DirectedAcyclicGraph) HasEdgeFromTo(u, v graph.Node) bool {
	edge := g.UndirectedGraph.EdgeBetween(u, v)
	if edge == nil {
		return false
	}
	return (edge.From().ID() == u.ID())
}

func (g *DirectedAcyclicGraph) From(n graph.Node) []graph.Node {
	if !g.Has(n) {
		return nil
	}

	fid := n.ID()
	nodes := make([]graph.Node, 0, g.UndirectedGraph.edges[n.ID()].Len())
	g.UndirectedGraph.edges[n.ID()].Visit(func(neighbor int, edge graph.Edge) {
		if edge.From().ID() == fid {
			nodes = append(nodes, g.UndirectedGraph.nodes[edge.To().ID()])
		}
	})
	return nodes
}

func (g *DirectedAcyclicGraph) VisitFrom(n graph.Node, visitor func(neighbor graph.Node) (shouldContinue bool)) {
	if !g.Has(n) {
		return
	}
	fid := n.ID()
	g.UndirectedGraph.edges[n.ID()].Visit(func(neighbor int, edge graph.Edge) {
		if edge.From().ID() == fid {
			if !visitor(g.UndirectedGraph.nodes[edge.To().ID()]) {
				return
			}
		}
	})
}

func (g *DirectedAcyclicGraph) To(n graph.Node) []graph.Node {
	if !g.Has(n) {
		return nil
	}

	tid := n.ID()
	nodes := make([]graph.Node, 0, g.UndirectedGraph.edges[n.ID()].Len())
	g.UndirectedGraph.edges[n.ID()].Visit(func(neighbor int, edge graph.Edge) {
		if edge.To().ID() == tid {
			nodes = append(nodes, g.UndirectedGraph.nodes[edge.From().ID()])
		}
	})
	return nodes
}

func (g *DirectedAcyclicGraph) VisitTo(n graph.Node, visitor func(neighbor graph.Node) (shouldContinue bool)) {
	if !g.Has(n) {
		return
	}
	tid := n.ID()
	g.UndirectedGraph.edges[n.ID()].Visit(func(neighbor int, edge graph.Edge) {
		if edge.To().ID() == tid {
			if !visitor(g.UndirectedGraph.nodes[edge.From().ID()]) {
				return
			}
		}
	})
}
