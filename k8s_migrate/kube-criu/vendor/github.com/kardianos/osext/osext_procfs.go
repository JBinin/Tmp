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
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux netbsd openbsd solaris dragonfly

package osext

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func executable() (string, error) {
	switch runtime.GOOS {
	case "linux":
		const deletedTag = " (deleted)"
		execpath, err := os.Readlink("/proc/self/exe")
		if err != nil {
			return execpath, err
		}
		execpath = strings.TrimSuffix(execpath, deletedTag)
		execpath = strings.TrimPrefix(execpath, deletedTag)
		return execpath, nil
	case "netbsd":
		return os.Readlink("/proc/curproc/exe")
	case "openbsd", "dragonfly":
		return os.Readlink("/proc/curproc/file")
	case "solaris":
		return os.Readlink(fmt.Sprintf("/proc/%d/path/a.out", os.Getpid()))
	}
	return "", errors.New("ExecPath not implemented for " + runtime.GOOS)
}
