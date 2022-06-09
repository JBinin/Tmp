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
package console

import (
	"errors"
	"io"
	"os"
)

var ErrNotAConsole = errors.New("provided file is not a console")

type Console interface {
	io.Reader
	io.Writer
	io.Closer

	// Resize resizes the console to the provided window size
	Resize(WinSize) error
	// ResizeFrom resizes the calling console to the size of the
	// provided console
	ResizeFrom(Console) error
	// SetRaw sets the console in raw mode
	SetRaw() error
	// DisableEcho disables echo on the console
	DisableEcho() error
	// Reset restores the console to its orignal state
	Reset() error
	// Size returns the window size of the console
	Size() (WinSize, error)
	// Fd returns the console's file descriptor
	Fd() uintptr
	// Name returns the console's file name
	Name() string
}

// WinSize specifies the window size of the console
type WinSize struct {
	// Height of the console
	Height uint16
	// Width of the console
	Width uint16
	x     uint16
	y     uint16
}

// Current returns the current processes console
func Current() Console {
	c, err := ConsoleFromFile(os.Stdin)
	if err != nil {
		// stdin should always be a console for the design
		// of this function
		panic(err)
	}
	return c
}

// ConsoleFromFile returns a console using the provided file
func ConsoleFromFile(f *os.File) (Console, error) {
	if err := checkConsole(f); err != nil {
		return nil, err
	}
	return newMaster(f)
}
