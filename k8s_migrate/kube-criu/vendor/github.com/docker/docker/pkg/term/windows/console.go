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
// +build windows

package windowsconsole

import (
	"os"

	"github.com/Azure/go-ansiterm/winterm"
)

// GetHandleInfo returns file descriptor and bool indicating whether the file is a console.
func GetHandleInfo(in interface{}) (uintptr, bool) {
	switch t := in.(type) {
	case *ansiReader:
		return t.Fd(), true
	case *ansiWriter:
		return t.Fd(), true
	}

	var inFd uintptr
	var isTerminal bool

	if file, ok := in.(*os.File); ok {
		inFd = file.Fd()
		isTerminal = IsConsole(inFd)
	}
	return inFd, isTerminal
}

// IsConsole returns true if the given file descriptor is a Windows Console.
// The code assumes that GetConsoleMode will return an error for file descriptors that are not a console.
func IsConsole(fd uintptr) bool {
	_, e := winterm.GetConsoleMode(fd)
	return e == nil
}
