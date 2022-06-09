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
// Code generated by "stringer -type=Kind"; DO NOT EDIT.

package width

import "fmt"

const _Kind_name = "NeutralEastAsianAmbiguousEastAsianWideEastAsianNarrowEastAsianFullwidthEastAsianHalfwidth"

var _Kind_index = [...]uint8{0, 7, 25, 38, 53, 71, 89}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)-1) {
		return fmt.Sprintf("Kind(%d)", i)
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
