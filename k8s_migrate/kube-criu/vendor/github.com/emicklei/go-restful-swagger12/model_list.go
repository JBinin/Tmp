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

// NamedModel associates a name with a Model (not using its Id)
type NamedModel struct {
	Name  string
	Model Model
}

// ModelList encapsulates a list of NamedModel (association)
type ModelList struct {
	List []NamedModel
}

// Put adds or replaces a Model by its name
func (l *ModelList) Put(name string, model Model) {
	for i, each := range l.List {
		if each.Name == name {
			// replace
			l.List[i] = NamedModel{name, model}
			return
		}
	}
	// add
	l.List = append(l.List, NamedModel{name, model})
}

// At returns a Model by its name, ok is false if absent
func (l *ModelList) At(name string) (m Model, ok bool) {
	for _, each := range l.List {
		if each.Name == name {
			return each.Model, true
		}
	}
	return m, false
}

// Do enumerates all the models, each with its assigned name
func (l *ModelList) Do(block func(name string, value Model)) {
	for _, each := range l.List {
		block(each.Name, each.Model)
	}
}

// MarshalJSON writes the ModelList as if it was a map[string]Model
func (l ModelList) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	buf.WriteString("{\n")
	for i, each := range l.List {
		buf.WriteString("\"")
		buf.WriteString(each.Name)
		buf.WriteString("\": ")
		encoder.Encode(each.Model)
		if i < len(l.List)-1 {
			buf.WriteString(",\n")
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// UnmarshalJSON reads back a ModelList. This is an expensive operation.
func (l *ModelList) UnmarshalJSON(data []byte) error {
	raw := map[string]interface{}{}
	json.NewDecoder(bytes.NewReader(data)).Decode(&raw)
	for k, v := range raw {
		// produces JSON bytes for each value
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		var m Model
		json.NewDecoder(bytes.NewReader(data)).Decode(&m)
		l.Put(k, m)
	}
	return nil
}
