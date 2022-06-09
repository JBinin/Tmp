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
package system

import (
	"os"
	"time"
)

// Chtimes changes the access time and modified time of a file at the given path
func Chtimes(name string, atime time.Time, mtime time.Time) error {
	unixMinTime := time.Unix(0, 0)
	unixMaxTime := maxTime

	// If the modified time is prior to the Unix Epoch, or after the
	// end of Unix Time, os.Chtimes has undefined behavior
	// default to Unix Epoch in this case, just in case

	if atime.Before(unixMinTime) || atime.After(unixMaxTime) {
		atime = unixMinTime
	}

	if mtime.Before(unixMinTime) || mtime.After(unixMaxTime) {
		mtime = unixMinTime
	}

	if err := os.Chtimes(name, atime, mtime); err != nil {
		return err
	}

	// Take platform specific action for setting create time.
	if err := setCTime(name, mtime); err != nil {
		return err
	}

	return nil
}
