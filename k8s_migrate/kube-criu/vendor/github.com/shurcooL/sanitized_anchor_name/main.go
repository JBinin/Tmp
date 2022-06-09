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
// Package sanitized_anchor_name provides a func to create sanitized anchor names.
//
// Its logic can be reused by multiple packages to create interoperable anchor names
// and links to those anchors.
//
// At this time, it does not try to ensure that generated anchor names
// are unique, that responsibility falls on the caller.
package sanitized_anchor_name

import "unicode"

// Create returns a sanitized anchor name for the given text.
func Create(text string) string {
	var anchorName []rune
	var futureDash = false
	for _, r := range []rune(text) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}
