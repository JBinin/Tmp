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

package svc

import (
	"errors"

	"golang.org/x/sys/windows"
)

// event represents auto-reset, initially non-signaled Windows event.
// It is used to communicate between go and asm parts of this package.
type event struct {
	h windows.Handle
}

func newEvent() (*event, error) {
	h, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return nil, err
	}
	return &event{h: h}, nil
}

func (e *event) Close() error {
	return windows.CloseHandle(e.h)
}

func (e *event) Set() error {
	return windows.SetEvent(e.h)
}

func (e *event) Wait() error {
	s, err := windows.WaitForSingleObject(e.h, windows.INFINITE)
	switch s {
	case windows.WAIT_OBJECT_0:
		break
	case windows.WAIT_FAILED:
		return err
	default:
		return errors.New("unexpected result from WaitForSingleObject")
	}
	return nil
}
