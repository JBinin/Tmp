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
// +build linux

package console

import (
	"io"
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

const (
	maxEvents = 128
)

// Epoller manages multiple epoll consoles using edge-triggered epoll api so we
// dont have to deal with repeated wake-up of EPOLLER or EPOLLHUP.
// For more details, see:
// - https://github.com/systemd/systemd/pull/4262
// - https://github.com/moby/moby/issues/27202
//
// Example usage of Epoller and EpollConsole can be as follow:
//
//	epoller, _ := NewEpoller()
//	epollConsole, _ := epoller.Add(console)
//	go epoller.Wait()
//	var (
//		b  bytes.Buffer
//		wg sync.WaitGroup
//	)
//	wg.Add(1)
//	go func() {
//		io.Copy(&b, epollConsole)
//		wg.Done()
//	}()
//	// perform I/O on the console
//	epollConsole.Shutdown(epoller.CloseConsole)
//	wg.Wait()
//	epollConsole.Close()
type Epoller struct {
	efd       int
	mu        sync.Mutex
	fdMapping map[int]*EpollConsole
}

// NewEpoller returns an instance of epoller with a valid epoll fd.
func NewEpoller() (*Epoller, error) {
	efd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	return &Epoller{
		efd:       efd,
		fdMapping: make(map[int]*EpollConsole),
	}, nil
}

// Add creates a epoll console based on the provided console. The console will
// be registered with EPOLLET (i.e. using edge-triggered notification) and its
// file descriptor will be set to non-blocking mode. After this, user should use
// the return console to perform I/O.
func (e *Epoller) Add(console Console) (*EpollConsole, error) {
	sysfd := int(console.Fd())
	// Set sysfd to non-blocking mode
	if err := unix.SetNonblock(sysfd, true); err != nil {
		return nil, err
	}

	ev := unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLOUT | unix.EPOLLRDHUP | unix.EPOLLET,
		Fd:     int32(sysfd),
	}
	if err := unix.EpollCtl(e.efd, unix.EPOLL_CTL_ADD, sysfd, &ev); err != nil {
		return nil, err
	}
	ef := &EpollConsole{
		Console: console,
		sysfd:   sysfd,
		readc:   sync.NewCond(&sync.Mutex{}),
		writec:  sync.NewCond(&sync.Mutex{}),
	}
	e.mu.Lock()
	e.fdMapping[sysfd] = ef
	e.mu.Unlock()
	return ef, nil
}

// Wait starts the loop to wait for its consoles' notifications and signal
// appropriate console that it can perform I/O.
func (e *Epoller) Wait() error {
	events := make([]unix.EpollEvent, maxEvents)
	for {
		n, err := unix.EpollWait(e.efd, events, -1)
		if err != nil {
			// EINTR: The call was interrupted by a signal handler before either
			// any of the requested events occurred or the timeout expired
			if err == unix.EINTR {
				continue
			}
			return err
		}
		for i := 0; i < n; i++ {
			ev := &events[i]
			// the console is ready to be read from
			if ev.Events&(unix.EPOLLIN|unix.EPOLLHUP|unix.EPOLLERR) != 0 {
				if epfile := e.getConsole(int(ev.Fd)); epfile != nil {
					epfile.signalRead()
				}
			}
			// the console is ready to be written to
			if ev.Events&(unix.EPOLLOUT|unix.EPOLLHUP|unix.EPOLLERR) != 0 {
				if epfile := e.getConsole(int(ev.Fd)); epfile != nil {
					epfile.signalWrite()
				}
			}
		}
	}
}

// Close unregister the console's file descriptor from epoll interface
func (e *Epoller) CloseConsole(fd int) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.fdMapping, fd)
	return unix.EpollCtl(e.efd, unix.EPOLL_CTL_DEL, fd, &unix.EpollEvent{})
}

func (e *Epoller) getConsole(sysfd int) *EpollConsole {
	e.mu.Lock()
	f := e.fdMapping[sysfd]
	e.mu.Unlock()
	return f
}

// Close the epoll fd
func (e *Epoller) Close() error {
	return unix.Close(e.efd)
}

// EpollConsole acts like a console but register its file descriptor with a
// epoll fd and uses epoll API to perform I/O.
type EpollConsole struct {
	Console
	readc  *sync.Cond
	writec *sync.Cond
	sysfd  int
	closed bool
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read
// (0 <= n <= len(p)) and any error encountered.
//
// If the console's read returns EAGAIN or EIO, we assumes that its a
// temporary error because the other side went away and wait for the signal
// generated by epoll event to continue.
func (ec *EpollConsole) Read(p []byte) (n int, err error) {
	var read int
	ec.readc.L.Lock()
	defer ec.readc.L.Unlock()
	for {
		read, err = ec.Console.Read(p[n:])
		n += read
		if err != nil {
			var hangup bool
			if perr, ok := err.(*os.PathError); ok {
				hangup = (perr.Err == unix.EAGAIN || perr.Err == unix.EIO)
			} else {
				hangup = (err == unix.EAGAIN || err == unix.EIO)
			}
			// if the other end disappear, assume this is temporary and wait for the
			// signal to continue again. Unless we didnt read anything and the
			// console is already marked as closed then we should exit
			if hangup && !(n == 0 && len(p) > 0 && ec.closed) {
				ec.readc.Wait()
				continue
			}
		}
		break
	}
	// if we didnt read anything then return io.EOF to end gracefully
	if n == 0 && len(p) > 0 && err == nil {
		err = io.EOF
	}
	// signal for others that we finished the read
	ec.readc.Signal()
	return n, err
}

// Writes len(p) bytes from p to the console. It returns the number of bytes
// written from p (0 <= n <= len(p)) and any error encountered that caused
// the write to stop early.
//
// If writes to the console returns EAGAIN or EIO, we assumes that its a
// temporary error because the other side went away and wait for the signal
// generated by epoll event to continue.
func (ec *EpollConsole) Write(p []byte) (n int, err error) {
	var written int
	ec.writec.L.Lock()
	defer ec.writec.L.Unlock()
	for {
		written, err = ec.Console.Write(p[n:])
		n += written
		if err != nil {
			var hangup bool
			if perr, ok := err.(*os.PathError); ok {
				hangup = (perr.Err == unix.EAGAIN || perr.Err == unix.EIO)
			} else {
				hangup = (err == unix.EAGAIN || err == unix.EIO)
			}
			// if the other end disappear, assume this is temporary and wait for the
			// signal to continue again.
			if hangup {
				ec.writec.Wait()
				continue
			}
		}
		// unrecoverable error, break the loop and return the error
		break
	}
	if n < len(p) && err == nil {
		err = io.ErrShortWrite
	}
	// signal for others that we finished the write
	ec.writec.Signal()
	return n, err
}

// Close closed the file descriptor and signal call waiters for this fd.
// It accepts a callback which will be called with the console's fd. The
// callback typically will be used to do further cleanup such as unregister the
// console's fd from the epoll interface.
// User should call Shutdown and wait for all I/O operation to be finished
// before closing the console.
func (ec *EpollConsole) Shutdown(close func(int) error) error {
	ec.readc.L.Lock()
	defer ec.readc.L.Unlock()
	ec.writec.L.Lock()
	defer ec.writec.L.Unlock()

	ec.readc.Broadcast()
	ec.writec.Broadcast()
	ec.closed = true
	return close(ec.sysfd)
}

// signalRead signals that the console is readable.
func (ec *EpollConsole) signalRead() {
	ec.readc.Signal()
}

// signalWrite signals that the console is writable.
func (ec *EpollConsole) signalWrite() {
	ec.writec.Signal()
}
