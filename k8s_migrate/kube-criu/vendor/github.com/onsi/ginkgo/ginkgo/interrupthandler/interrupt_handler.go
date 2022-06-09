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
package interrupthandler

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type InterruptHandler struct {
	interruptCount int
	lock           *sync.Mutex
	C              chan bool
}

func NewInterruptHandler() *InterruptHandler {
	h := &InterruptHandler{
		lock: &sync.Mutex{},
		C:    make(chan bool, 0),
	}

	go h.handleInterrupt()
	SwallowSigQuit()

	return h
}

func (h *InterruptHandler) WasInterrupted() bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.interruptCount > 0
}

func (h *InterruptHandler) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	signal.Stop(c)

	h.lock.Lock()
	h.interruptCount++
	if h.interruptCount == 1 {
		close(h.C)
	} else if h.interruptCount > 5 {
		os.Exit(1)
	}
	h.lock.Unlock()

	go h.handleInterrupt()
}
