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
package signal

import (
	"syscall"
)

// SignalMap is a map of FreeBSD signals.
var SignalMap = map[string]syscall.Signal{
	"ABRT":   syscall.SIGABRT,
	"ALRM":   syscall.SIGALRM,
	"BUF":    syscall.SIGBUS,
	"CHLD":   syscall.SIGCHLD,
	"CONT":   syscall.SIGCONT,
	"EMT":    syscall.SIGEMT,
	"FPE":    syscall.SIGFPE,
	"HUP":    syscall.SIGHUP,
	"ILL":    syscall.SIGILL,
	"INFO":   syscall.SIGINFO,
	"INT":    syscall.SIGINT,
	"IO":     syscall.SIGIO,
	"IOT":    syscall.SIGIOT,
	"KILL":   syscall.SIGKILL,
	"LWP":    syscall.SIGLWP,
	"PIPE":   syscall.SIGPIPE,
	"PROF":   syscall.SIGPROF,
	"QUIT":   syscall.SIGQUIT,
	"SEGV":   syscall.SIGSEGV,
	"STOP":   syscall.SIGSTOP,
	"SYS":    syscall.SIGSYS,
	"TERM":   syscall.SIGTERM,
	"THR":    syscall.SIGTHR,
	"TRAP":   syscall.SIGTRAP,
	"TSTP":   syscall.SIGTSTP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"URG":    syscall.SIGURG,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
	"VTALRM": syscall.SIGVTALRM,
	"WINCH":  syscall.SIGWINCH,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
}
