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
package command

import (
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/term"
)

// OutStream is an output stream used by the DockerCli to write normal program
// output.
type OutStream struct {
	out        io.Writer
	fd         uintptr
	isTerminal bool
	state      *term.State
}

func (o *OutStream) Write(p []byte) (int, error) {
	return o.out.Write(p)
}

// FD returns the file descriptor number for this stream
func (o *OutStream) FD() uintptr {
	return o.fd
}

// IsTerminal returns true if this stream is connected to a terminal
func (o *OutStream) IsTerminal() bool {
	return o.isTerminal
}

// SetRawTerminal sets raw mode on the output terminal
func (o *OutStream) SetRawTerminal() (err error) {
	if os.Getenv("NORAW") != "" || !o.isTerminal {
		return nil
	}
	o.state, err = term.SetRawTerminalOutput(o.fd)
	return err
}

// RestoreTerminal restores normal mode to the terminal
func (o *OutStream) RestoreTerminal() {
	if o.state != nil {
		term.RestoreTerminal(o.fd, o.state)
	}
}

// GetTtySize returns the height and width in characters of the tty
func (o *OutStream) GetTtySize() (uint, uint) {
	if !o.isTerminal {
		return 0, 0
	}
	ws, err := term.GetWinsize(o.fd)
	if err != nil {
		logrus.Debugf("Error getting size: %s", err)
		if ws == nil {
			return 0, 0
		}
	}
	return uint(ws.Height), uint(ws.Width)
}

// NewOutStream returns a new OutStream object from a Writer
func NewOutStream(out io.Writer) *OutStream {
	fd, isTerminal := term.GetFdInfo(out)
	return &OutStream{out: out, fd: fd, isTerminal: isTerminal}
}
