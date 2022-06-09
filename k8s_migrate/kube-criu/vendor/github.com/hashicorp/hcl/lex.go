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
package hcl

import (
	"unicode"
	"unicode/utf8"
)

type lexModeValue byte

const (
	lexModeUnknown lexModeValue = iota
	lexModeHcl
	lexModeJson
)

// lexMode returns whether we're going to be parsing in JSON
// mode or HCL mode.
func lexMode(v []byte) lexModeValue {
	var (
		r      rune
		w      int
		offset int
	)

	for {
		r, w = utf8.DecodeRune(v[offset:])
		offset += w
		if unicode.IsSpace(r) {
			continue
		}
		if r == '{' {
			return lexModeJson
		}
		break
	}

	return lexModeHcl
}
