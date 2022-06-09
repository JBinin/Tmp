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
// +build !linux,!solaris

package kernel

import (
	"errors"
)

// Utsname represents the system name structure.
// It is defined here to make it portable as it is available on linux but not
// on windows.
type Utsname struct {
	Release [65]byte
}

func uname() (*Utsname, error) {
	return nil, errors.New("Kernel version detection is available only on linux")
}
