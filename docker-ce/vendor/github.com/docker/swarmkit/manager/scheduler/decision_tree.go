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
package scheduler

import (
	"container/heap"
)

type decisionTree struct {
	// Count of tasks for the service scheduled to this subtree
	tasks int

	// Non-leaf point to the next level of the tree. The key is the
	// value that the subtree covers.
	next map[string]*decisionTree

	// Leaf nodes contain a list of nodes
	nodeHeap nodeMaxHeap
}

// orderedNodes returns the nodes in this decision tree entry, sorted best
// (lowest) first according to the sorting function. Must be called on a leaf
// of the decision tree.
//
// The caller may modify the nodes in the returned slice. This has the effect
// of changing the nodes in the decision tree entry. The next node to
// findBestNodes on this decisionTree entry will take into account the changes
// that were made to the nodes.
func (dt *decisionTree) orderedNodes(meetsConstraints func(*NodeInfo) bool, nodeLess func(*NodeInfo, *NodeInfo) bool) []NodeInfo {
	if dt.nodeHeap.length != len(dt.nodeHeap.nodes) {
		// We already collapsed the heap into a sorted slice, so
		// re-heapify. There may have been modifications to the nodes
		// so we can't return dt.nodeHeap.nodes as-is. We also need to
		// reevaluate constraints because of the possible modifications.
		for i := 0; i < len(dt.nodeHeap.nodes); {
			if meetsConstraints(&dt.nodeHeap.nodes[i]) {
				i++
			} else {
				last := len(dt.nodeHeap.nodes) - 1
				dt.nodeHeap.nodes[i] = dt.nodeHeap.nodes[last]
				dt.nodeHeap.nodes = dt.nodeHeap.nodes[:last]
			}
		}
		dt.nodeHeap.length = len(dt.nodeHeap.nodes)
		heap.Init(&dt.nodeHeap)
	}

	// Popping every element orders the nodes from best to worst. The
	// first pop gets the worst node (since this a max-heap), and puts it
	// at position n-1. Then the next pop puts the next-worst at n-2, and
	// so on.
	for dt.nodeHeap.Len() > 0 {
		heap.Pop(&dt.nodeHeap)
	}

	return dt.nodeHeap.nodes
}
