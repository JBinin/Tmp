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
package term

import (
	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TCGETS
	setTermios = unix.TCSETS
)

// Termios is the Unix API for terminal I/O.
type Termios unix.Termios

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd uintptr) (*State, error) {
	termios, err := unix.IoctlGetTermios(int(fd), getTermios)
	if err != nil {
		return nil, err
	}

	var oldState State
	oldState.termios = Termios(*termios)

	termios.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termios.Cflag &^= (unix.CSIZE | unix.PARENB)
	termios.Cflag |= unix.CS8

	if err := unix.IoctlSetTermios(int(fd), setTermios, termios); err != nil {
		return nil, err
	}
	return &oldState, nil
}
