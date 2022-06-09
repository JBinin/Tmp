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
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uritemplates

// Expand parses then expands a URI template with a set of values to produce
// the resultant URI. Two forms of the result are returned: one with all the
// elements escaped, and one with the elements unescaped.
func Expand(path string, values map[string]string) (escaped, unescaped string, err error) {
	template, err := parse(path)
	if err != nil {
		return "", "", err
	}
	escaped, unescaped = template.Expand(values)
	return escaped, unescaped, nil
}
