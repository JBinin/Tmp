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
	"fmt"
	"runtime"
	"syscall"
)

// Via http://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/commit/?id=7b21fddd087678a70ad64afc0f632e0f1071b092
//
// We need different setns values for the different platforms and arch
// We are declaring the macro here because the SETNS syscall does not exist in th stdlib
var setNsMap = map[string]uintptr{
	"linux/386":     346,
	"linux/arm64":   268,
	"linux/amd64":   308,
	"linux/arm":     375,
	"linux/ppc":     350,
	"linux/ppc64":   350,
	"linux/ppc64le": 350,
	"linux/s390x":   339,
}

var sysSetns = setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]

func SysSetns() uint32 {
	return uint32(sysSetns)
}

func Setns(fd uintptr, flags uintptr) error {
	ns, exists := setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !exists {
		return fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	_, _, err := syscall.RawSyscall(ns, fd, flags, 0)
	if err != 0 {
		return err
	}
	return nil
}
