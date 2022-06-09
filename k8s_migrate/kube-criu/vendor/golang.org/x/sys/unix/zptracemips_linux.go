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
// Code generated by linux/mkall.go generatePtracePair(mips, mips64). DO NOT EDIT.

// +build linux
// +build mips mips64

package unix

import "unsafe"

// PtraceRegsMips is the registers used by mips binaries.
type PtraceRegsMips struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}

// PtraceGetRegsMips fetches the registers used by mips binaries.
func PtraceGetRegsMips(pid int, regsout *PtraceRegsMips) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}

// PtraceSetRegsMips sets the registers used by mips binaries.
func PtraceSetRegsMips(pid int, regs *PtraceRegsMips) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}

// PtraceRegsMips64 is the registers used by mips64 binaries.
type PtraceRegsMips64 struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}

// PtraceGetRegsMips64 fetches the registers used by mips64 binaries.
func PtraceGetRegsMips64(pid int, regsout *PtraceRegsMips64) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}

// PtraceSetRegsMips64 sets the registers used by mips64 binaries.
func PtraceSetRegsMips64(pid int, regs *PtraceRegsMips64) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}
