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
package swagger

// Copyright 2015 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"bytes"
	"encoding/json"
)

// ApiDeclarationList maintains an ordered list of ApiDeclaration.
type ApiDeclarationList struct {
	List []ApiDeclaration
}

// At returns the ApiDeclaration by its path unless absent, then ok is false
func (l *ApiDeclarationList) At(path string) (a ApiDeclaration, ok bool) {
	for _, each := range l.List {
		if each.ResourcePath == path {
			return each, true
		}
	}
	return a, false
}

// Put adds or replaces a ApiDeclaration with this name
func (l *ApiDeclarationList) Put(path string, a ApiDeclaration) {
	// maybe replace existing
	for i, each := range l.List {
		if each.ResourcePath == path {
			// replace
			l.List[i] = a
			return
		}
	}
	// add
	l.List = append(l.List, a)
}

// Do enumerates all the properties, each with its assigned name
func (l *ApiDeclarationList) Do(block func(path string, decl ApiDeclaration)) {
	for _, each := range l.List {
		block(each.ResourcePath, each)
	}
}

// MarshalJSON writes the ModelPropertyList as if it was a map[string]ModelProperty
func (l ApiDeclarationList) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	buf.WriteString("{\n")
	for i, each := range l.List {
		buf.WriteString("\"")
		buf.WriteString(each.ResourcePath)
		buf.WriteString("\": ")
		encoder.Encode(each)
		if i < len(l.List)-1 {
			buf.WriteString(",\n")
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
