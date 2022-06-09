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
package volume

import "strings"

// {<copy mode>=isEnabled}
var copyModes = map[string]bool{
	"nocopy": false,
}

func copyModeExists(mode string) bool {
	_, exists := copyModes[mode]
	return exists
}

// GetCopyMode gets the copy mode from the mode string for mounts
func getCopyMode(mode string) (bool, bool) {
	for _, o := range strings.Split(mode, ",") {
		if isEnabled, exists := copyModes[o]; exists {
			return isEnabled, true
		}
	}
	return DefaultCopyMode, false
}
