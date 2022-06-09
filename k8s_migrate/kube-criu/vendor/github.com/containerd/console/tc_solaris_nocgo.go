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
// +build solaris,!cgo

//
// Implementing the functions below requires cgo support.  Non-cgo stubs
// versions are defined below to enable cross-compilation of source code
// that depends on these functions, but the resultant cross-compiled
// binaries cannot actually be used.  If the stub function(s) below are
// actually invoked they will display an error message and cause the
// calling process to exit.
//

package console

import (
	"os"

	"golang.org/x/sys/unix"
)

const (
	cmdTcGet = unix.TCGETS
	cmdTcSet = unix.TCSETS
)

func ptsname(f *os.File) (string, error) {
	panic("ptsname() support requires cgo.")
}

func unlockpt(f *os.File) error {
	panic("unlockpt() support requires cgo.")
}
