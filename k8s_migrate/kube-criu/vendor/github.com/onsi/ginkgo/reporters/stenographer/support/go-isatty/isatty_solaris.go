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
// +build solaris
// +build !appengine

package isatty

import (
	"golang.org/x/sys/unix"
)

// IsTerminal returns true if the given file descriptor is a terminal.
// see: http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libbc/libc/gen/common/isatty.c
func IsTerminal(fd uintptr) bool {
	var termio unix.Termio
	err := unix.IoctlSetTermio(int(fd), unix.TCGETA, &termio)
	return err == nil
}
