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
// Copyright 2013 MongoDB, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author           tolsen
// author-github    https://github.com/tolsen
//
// repository-name  gojsonschema
// repository-desc  An implementation of JSON Schema, based on IETF's draft v4 - Go language.
//
// description      Implements a persistent (immutable w/ shared structure) singly-linked list of strings for the purpose of storing a json context
//
// created          04-09-2013

package gojsonschema

import "bytes"

// jsonContext implements a persistent linked-list of strings
type jsonContext struct {
	head string
	tail *jsonContext
}

func newJsonContext(head string, tail *jsonContext) *jsonContext {
	return &jsonContext{head, tail}
}

// String displays the context in reverse.
// This plays well with the data structure's persistent nature with
// Cons and a json document's tree structure.
func (c *jsonContext) String(del ...string) string {
	byteArr := make([]byte, 0, c.stringLen())
	buf := bytes.NewBuffer(byteArr)
	c.writeStringToBuffer(buf, del)

	return buf.String()
}

func (c *jsonContext) stringLen() int {
	length := 0
	if c.tail != nil {
		length = c.tail.stringLen() + 1 // add 1 for "."
	}

	length += len(c.head)
	return length
}

func (c *jsonContext) writeStringToBuffer(buf *bytes.Buffer, del []string) {
	if c.tail != nil {
		c.tail.writeStringToBuffer(buf, del)

		if len(del) > 0 {
			buf.WriteString(del[0])
		} else {
			buf.WriteString(".")
		}
	}

	buf.WriteString(c.head)
}
