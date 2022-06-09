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

package signal

import (
	"syscall"
)

// Signals used in cli/command (no windows equivalent, use
// invalid signals so they don't get handled)
const (
	SIGCHLD  = syscall.Signal(0xff)
	SIGWINCH = syscall.Signal(0xff)
	SIGPIPE  = syscall.Signal(0xff)
	// DefaultStopSignal is the syscall signal used to stop a container in windows systems.
	DefaultStopSignal = "15"
)

// SignalMap is a map of "supported" signals. As per the comment in GOLang's
// ztypes_windows.go: "More invented values for signals". Windows doesn't
// really support signals in any way, shape or form that Unix does.
//
// We have these so that docker kill can be used to gracefully (TERM) and
// forcibly (KILL) terminate a container on Windows.
var SignalMap = map[string]syscall.Signal{
	"KILL": syscall.SIGKILL,
	"TERM": syscall.SIGTERM,
}
