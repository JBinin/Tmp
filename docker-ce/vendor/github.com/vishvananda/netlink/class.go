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
package netlink

import (
	"fmt"
)

type Class interface {
	Attrs() *ClassAttrs
	Type() string
}

// ClassAttrs represents a netlink class. A filter is associated with a link,
// has a handle and a parent. The root filter of a device should have a
// parent == HANDLE_ROOT.
type ClassAttrs struct {
	LinkIndex int
	Handle    uint32
	Parent    uint32
	Leaf      uint32
}

func (q ClassAttrs) String() string {
	return fmt.Sprintf("{LinkIndex: %d, Handle: %s, Parent: %s, Leaf: %d}", q.LinkIndex, HandleStr(q.Handle), HandleStr(q.Parent), q.Leaf)
}

type HtbClassAttrs struct {
	// TODO handle all attributes
	Rate    uint64
	Ceil    uint64
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
}

func (q HtbClassAttrs) String() string {
	return fmt.Sprintf("{Rate: %d, Ceil: %d, Buffer: %d, Cbuffer: %d}", q.Rate, q.Ceil, q.Buffer, q.Cbuffer)
}

// HtbClass represents an Htb class
type HtbClass struct {
	ClassAttrs
	Rate    uint64
	Ceil    uint64
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
}

func (q HtbClass) String() string {
	return fmt.Sprintf("{Rate: %d, Ceil: %d, Buffer: %d, Cbuffer: %d}", q.Rate, q.Ceil, q.Buffer, q.Cbuffer)
}

func (q *HtbClass) Attrs() *ClassAttrs {
	return &q.ClassAttrs
}

func (q *HtbClass) Type() string {
	return "htb"
}

// GenericClass classes represent types that are not currently understood
// by this netlink library.
type GenericClass struct {
	ClassAttrs
	ClassType string
}

func (class *GenericClass) Attrs() *ClassAttrs {
	return &class.ClassAttrs
}

func (class *GenericClass) Type() string {
	return class.ClassType
}
