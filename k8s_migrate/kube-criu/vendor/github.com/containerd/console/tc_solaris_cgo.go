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
// +build solaris,cgo

package console

import (
	"os"

	"golang.org/x/sys/unix"
)

//#include <stdlib.h>
import "C"

const (
	cmdTcGet = unix.TCGETS
	cmdTcSet = unix.TCSETS
)

// ptsname retrieves the name of the first available pts for the given master.
func ptsname(f *os.File) (string, error) {
	ptspath, err := C.ptsname(C.int(f.Fd()))
	if err != nil {
		return "", err
	}
	return C.GoString(ptspath), nil
}

// unlockpt unlocks the slave pseudoterminal device corresponding to the master pseudoterminal referred to by f.
// unlockpt should be called before opening the slave side of a pty.
func unlockpt(f *os.File) error {
	if _, err := C.grantpt(C.int(f.Fd())); err != nil {
		return err
	}
	return nil
}
