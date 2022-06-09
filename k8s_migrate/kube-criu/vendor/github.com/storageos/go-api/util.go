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
package storageos

import (
	"fmt"
	"strings"
)

// ParseRef is a helper to split out the namespace and name from a path
// reference.
func ParseRef(ref string) (namespace string, name string, err error) {
	parts := strings.Split(ref, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Name must be prefixed with <namespace>/")
	}
	return parts[0], parts[1], nil
}
