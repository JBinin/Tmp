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
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package debug provides facilities to execute svc.Handler on console.
//
package debug

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/windows/svc"
)

// Run executes service name by calling appropriate handler function.
// The process is running on console, unlike real service. Use Ctrl+C to
// send "Stop" command to your service.
func Run(name string, handler svc.Handler) error {
	cmds := make(chan svc.ChangeRequest)
	changes := make(chan svc.Status)

	sig := make(chan os.Signal)
	signal.Notify(sig)

	go func() {
		status := svc.Status{State: svc.Stopped}
		for {
			select {
			case <-sig:
				cmds <- svc.ChangeRequest{svc.Stop, status}
			case status = <-changes:
			}
		}
	}()

	_, errno := handler.Execute([]string{name}, cmds, changes)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}
