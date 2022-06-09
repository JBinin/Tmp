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
package ast

import "fmt"

// WalkFunc describes a function to be called for each node during a Walk. The
// returned node can be used to rewrite the AST. Walking stops the returned
// bool is false.
type WalkFunc func(Node) (Node, bool)

// Walk traverses an AST in depth-first order: It starts by calling fn(node);
// node must not be nil. If fn returns true, Walk invokes fn recursively for
// each of the non-nil children of node, followed by a call of fn(nil). The
// returned node of fn can be used to rewrite the passed node to fn.
func Walk(node Node, fn WalkFunc) Node {
	rewritten, ok := fn(node)
	if !ok {
		return rewritten
	}

	switch n := node.(type) {
	case *File:
		n.Node = Walk(n.Node, fn)
	case *ObjectList:
		for i, item := range n.Items {
			n.Items[i] = Walk(item, fn).(*ObjectItem)
		}
	case *ObjectKey:
		// nothing to do
	case *ObjectItem:
		for i, k := range n.Keys {
			n.Keys[i] = Walk(k, fn).(*ObjectKey)
		}

		if n.Val != nil {
			n.Val = Walk(n.Val, fn)
		}
	case *LiteralType:
		// nothing to do
	case *ListType:
		for i, l := range n.List {
			n.List[i] = Walk(l, fn)
		}
	case *ObjectType:
		n.List = Walk(n.List, fn).(*ObjectList)
	default:
		// should we panic here?
		fmt.Printf("unknown type: %T\n", n)
	}

	fn(nil)
	return rewritten
}
