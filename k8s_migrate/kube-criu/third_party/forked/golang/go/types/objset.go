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
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements objsets.
//
// An objset is similar to a Scope but objset elements
// are identified by their unique id, instead of their
// object name.

package types

// An objset is a set of objects identified by their unique id.
// The zero value for objset is a ready-to-use empty objset.
type objset map[string]Object // initialized lazily

// insert attempts to insert an object obj into objset s.
// If s already contains an alternative object alt with
// the same name, insert leaves s unchanged and returns alt.
// Otherwise it inserts obj and returns nil.
func (s *objset) insert(obj Object) Object {
	id := obj.Id()
	if alt := (*s)[id]; alt != nil {
		return alt
	}
	if *s == nil {
		*s = make(map[string]Object)
	}
	(*s)[id] = obj
	return nil
}
