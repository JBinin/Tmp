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
// +build linux

package homedir

import (
	"os"

	"github.com/docker/docker/pkg/idtools"
)

// GetStatic returns the home directory for the current user without calling
// os/user.Current(). This is useful for static-linked binary on glibc-based
// system, because a call to os/user.Current() in a static binary leads to
// segfault due to a glibc issue that won't be fixed in a short term.
// (#29344, golang/go#13470, https://sourceware.org/bugzilla/show_bug.cgi?id=19341)
func GetStatic() (string, error) {
	uid := os.Getuid()
	usr, err := idtools.LookupUID(uid)
	if err != nil {
		return "", err
	}
	return usr.Home, nil
}
