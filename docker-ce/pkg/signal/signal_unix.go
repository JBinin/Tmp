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
// +build !windows

package signal

import (
	"syscall"
)

// Signals used in cli/command (no windows equivalent, use
// invalid signals so they don't get handled)

const (
	// SIGCHLD is a signal sent to a process when a child process terminates, is interrupted, or resumes after being interrupted.
	SIGCHLD = syscall.SIGCHLD
	// SIGWINCH is a signal sent to a process when its controlling terminal changes its size
	SIGWINCH = syscall.SIGWINCH
	// SIGPIPE is a signal sent to a process when a pipe is written to before the other end is open for reading
	SIGPIPE = syscall.SIGPIPE
	// DefaultStopSignal is the syscall signal used to stop a container in unix systems.
	DefaultStopSignal = "SIGTERM"
)
