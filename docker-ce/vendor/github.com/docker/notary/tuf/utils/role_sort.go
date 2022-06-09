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
package utils

import (
	"strings"
)

// RoleList is a list of roles
type RoleList []string

// Len returns the length of the list
func (r RoleList) Len() int {
	return len(r)
}

// Less returns true if the item at i should be sorted
// before the item at j. It's an unstable partial ordering
// based on the number of segments, separated by "/", in
// the role name
func (r RoleList) Less(i, j int) bool {
	segsI := strings.Split(r[i], "/")
	segsJ := strings.Split(r[j], "/")
	if len(segsI) == len(segsJ) {
		return r[i] < r[j]
	}
	return len(segsI) < len(segsJ)
}

// Swap the items at 2 locations in the list
func (r RoleList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
