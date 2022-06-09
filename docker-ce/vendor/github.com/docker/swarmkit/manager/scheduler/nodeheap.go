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

type nodeMaxHeap struct {
	nodes    []NodeInfo
	lessFunc func(*NodeInfo, *NodeInfo) bool
	length   int
}

func (h nodeMaxHeap) Len() int {
	return h.length
}

func (h nodeMaxHeap) Swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
}

func (h nodeMaxHeap) Less(i, j int) bool {
	// reversed to make a max-heap
	return h.lessFunc(&h.nodes[j], &h.nodes[i])
}

func (h *nodeMaxHeap) Push(x interface{}) {
	h.nodes = append(h.nodes, x.(NodeInfo))
	h.length++
}

func (h *nodeMaxHeap) Pop() interface{} {
	h.length--
	// return value is never used
	return nil
}
