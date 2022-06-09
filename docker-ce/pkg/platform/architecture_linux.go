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
// Package platform provides helper function to get the runtime architecture
// for different platforms.
package platform

import (
	"syscall"
)

// runtimeArchitecture gets the name of the current architecture (x86, x86_64, …)
func runtimeArchitecture() (string, error) {
	utsname := &syscall.Utsname{}
	if err := syscall.Uname(utsname); err != nil {
		return "", err
	}
	return charsToString(utsname.Machine), nil
}
