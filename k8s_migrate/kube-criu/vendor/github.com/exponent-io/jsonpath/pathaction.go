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
package jsonpath

// pathNode is used to construct a trie of paths to be matched
type pathNode struct {
	matchOn    interface{} // string, or integer
	childNodes []pathNode
	action     DecodeAction
}

// match climbs the trie to find a node that matches the given JSON path.
func (n *pathNode) match(path JsonPath) *pathNode {
	var node *pathNode = n
	for _, ps := range path {
		found := false
		for i, n := range node.childNodes {
			if n.matchOn == ps {
				node = &node.childNodes[i]
				found = true
				break
			} else if _, ok := ps.(int); ok && n.matchOn == AnyIndex {
				node = &node.childNodes[i]
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return node
}

// PathActions represents a collection of DecodeAction functions that should be called at certain path positions
// when scanning the JSON stream. PathActions can be created once and used many times in one or more JSON streams.
type PathActions struct {
	node pathNode
}

// DecodeAction handlers are called by the Decoder when scanning objects. See PathActions.Add for more detail.
type DecodeAction func(d *Decoder) error

// Add specifies an action to call on the Decoder when the specified path is encountered.
func (je *PathActions) Add(action DecodeAction, path ...interface{}) {

	var node *pathNode = &je.node
	for _, ps := range path {
		found := false
		for i, n := range node.childNodes {
			if n.matchOn == ps {
				node = &node.childNodes[i]
				found = true
				break
			}
		}
		if !found {
			node.childNodes = append(node.childNodes, pathNode{matchOn: ps})
			node = &node.childNodes[len(node.childNodes)-1]
		}
	}
	node.action = action
}
