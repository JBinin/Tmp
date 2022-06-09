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

package debug

import (
	"os"
	"strconv"
)

// Log interface allows different log implementations to be used.
type Log interface {
	Close() error
	Info(eid uint32, msg string) error
	Warning(eid uint32, msg string) error
	Error(eid uint32, msg string) error
}

// ConsoleLog provides access to the console.
type ConsoleLog struct {
	Name string
}

// New creates new ConsoleLog.
func New(source string) *ConsoleLog {
	return &ConsoleLog{Name: source}
}

// Close closes console log l.
func (l *ConsoleLog) Close() error {
	return nil
}

func (l *ConsoleLog) report(kind string, eid uint32, msg string) error {
	s := l.Name + "." + kind + "(" + strconv.Itoa(int(eid)) + "): " + msg + "\n"
	_, err := os.Stdout.Write([]byte(s))
	return err
}

// Info writes an information event msg with event id eid to the console l.
func (l *ConsoleLog) Info(eid uint32, msg string) error {
	return l.report("info", eid, msg)
}

// Warning writes an warning event msg with event id eid to the console l.
func (l *ConsoleLog) Warning(eid uint32, msg string) error {
	return l.report("warn", eid, msg)
}

// Error writes an error event msg with event id eid to the console l.
func (l *ConsoleLog) Error(eid uint32, msg string) error {
	return l.report("error", eid, msg)
}
