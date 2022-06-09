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
package zfs

import (
	"fmt"
)

// Error is an error which is returned when the `zfs` or `zpool` shell
// commands return with a non-zero exit code.
type Error struct {
	Err    error
	Debug  string
	Stderr string
}

// Error returns the string representation of an Error.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %q => %s", e.Err, e.Debug, e.Stderr)
}
