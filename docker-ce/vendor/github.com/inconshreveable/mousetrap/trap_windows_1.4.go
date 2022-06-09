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
// +build go1.4

package mousetrap

import (
	"os"
	"syscall"
	"unsafe"
)

func getProcessEntry(pid int) (*syscall.ProcessEntry32, error) {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)
	var procEntry syscall.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = syscall.Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}
	for {
		if procEntry.ProcessID == uint32(pid) {
			return &procEntry, nil
		}
		err = syscall.Process32Next(snapshot, &procEntry)
		if err != nil {
			return nil, err
		}
	}
}

// StartedByExplorer returns true if the program was invoked by the user double-clicking
// on the executable from explorer.exe
//
// It is conservative and returns false if any of the internal calls fail.
// It does not guarantee that the program was run from a terminal. It only can tell you
// whether it was launched from explorer.exe
func StartedByExplorer() bool {
	pe, err := getProcessEntry(os.Getppid())
	if err != nil {
		return false
	}
	return "explorer.exe" == syscall.UTF16ToString(pe.ExeFile[:])
}
