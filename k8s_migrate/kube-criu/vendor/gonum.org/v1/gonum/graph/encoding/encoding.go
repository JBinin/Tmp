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
// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding

import "gonum.org/v1/gonum/graph"

// Builder is a graph that can have user-defined nodes and edges added.
type Builder interface {
	graph.Graph
	graph.Builder
}

// AttributeSetter is implemented by types that can set an encoded graph
// attribute.
type AttributeSetter interface {
	SetAttribute(Attribute) error
}

// Attributer defines graph.Node or graph.Edge values that can
// specify graph attributes.
type Attributer interface {
	Attributes() []Attribute
}

// Attribute is an encoded key value attribute pair use in graph encoding.
type Attribute struct {
	Key, Value string
}
