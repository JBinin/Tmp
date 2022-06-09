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

package volume

import mounttypes "github.com/docker/docker/api/types/mount"

// DefaultPropagationMode is used only in linux. In other cases it returns
// empty string.
const DefaultPropagationMode mounttypes.Propagation = ""

// propagation modes not supported on this platform.
var propagationModes = map[mounttypes.Propagation]bool{}

// GetPropagation is not supported. Return empty string.
func GetPropagation(mode string) mounttypes.Propagation {
	return DefaultPropagationMode
}

// HasPropagation checks if there is a valid propagation mode present in
// passed string. Returns true if a valid propagation mode specifier is
// present, false otherwise.
func HasPropagation(mode string) bool {
	return false
}
