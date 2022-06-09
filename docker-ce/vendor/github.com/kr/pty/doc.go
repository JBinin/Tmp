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
// Package pty provides functions for working with Unix terminals.
package pty

import (
	"errors"
	"os"
)

// ErrUnsupported is returned if a function is not
// available on the current platform.
var ErrUnsupported = errors.New("unsupported")

// Opens a pty and its corresponding tty.
func Open() (pty, tty *os.File, err error) {
	return open()
}
