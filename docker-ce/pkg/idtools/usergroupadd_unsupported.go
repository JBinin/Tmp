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
// +build !linux

package idtools

import "fmt"

// AddNamespaceRangesUser takes a name and finds an unused uid, gid pair
// and calls the appropriate helper function to add the group and then
// the user to the group in /etc/group and /etc/passwd respectively.
func AddNamespaceRangesUser(name string) (int, int, error) {
	return -1, -1, fmt.Errorf("No support for adding users or groups on this OS")
}
