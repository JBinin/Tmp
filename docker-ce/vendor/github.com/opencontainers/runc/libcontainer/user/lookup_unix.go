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
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package user

import (
	"io"
	"os"
)

// Unix-specific path to the passwd and group formatted files.
const (
	unixPasswdPath = "/etc/passwd"
	unixGroupPath  = "/etc/group"
)

func GetPasswdPath() (string, error) {
	return unixPasswdPath, nil
}

func GetPasswd() (io.ReadCloser, error) {
	return os.Open(unixPasswdPath)
}

func GetGroupPath() (string, error) {
	return unixGroupPath, nil
}

func GetGroup() (io.ReadCloser, error) {
	return os.Open(unixGroupPath)
}
