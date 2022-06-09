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
package semver

import (
	"sort"
)

// Versions represents multiple versions.
type Versions []Version

// Len returns length of version collection
func (s Versions) Len() int {
	return len(s)
}

// Swap swaps two versions inside the collection by its indices
func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less checks if version at index i is less than version at index j
func (s Versions) Less(i, j int) bool {
	return s[i].LT(s[j])
}

// Sort sorts a slice of versions
func Sort(versions []Version) {
	sort.Sort(Versions(versions))
}
