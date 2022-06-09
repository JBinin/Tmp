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
//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package utils

import (
	"sort"
)

// Check if a sorted string list has a string
func SortedStringHas(s sort.StringSlice, x string) bool {
	index := s.Search(x)
	if index == len(s) {
		return false
	}
	return s[s.Search(x)] == x
}

// Delete a string from a sorted string list
func SortedStringsDelete(s sort.StringSlice, x string) sort.StringSlice {
	index := s.Search(x)
	if len(s) != index && s[index] == x {
		s = append(s[:index], s[index+1:]...)
	}

	return s
}
