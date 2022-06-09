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
// Package pq is a pure Go Postgres driver for the database/sql package.
package pq

import (
	"path/filepath"
	"syscall"
)

// Perform Windows user name lookup identically to libpq.
//
// The PostgreSQL code makes use of the legacy Win32 function
// GetUserName, and that function has not been imported into stock Go.
// GetUserNameEx is available though, the difference being that a
// wider range of names are available.  To get the output to be the
// same as GetUserName, only the base (or last) component of the
// result is returned.
func userCurrent() (string, error) {
	pw_name := make([]uint16, 128)
	pwname_size := uint32(len(pw_name)) - 1
	err := syscall.GetUserNameEx(syscall.NameSamCompatible, &pw_name[0], &pwname_size)
	if err != nil {
		return "", ErrCouldNotDetectUsername
	}
	s := syscall.UTF16ToString(pw_name)
	u := filepath.Base(s)
	return u, nil
}
