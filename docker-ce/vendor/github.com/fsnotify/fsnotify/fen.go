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
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package fsnotify

import (
	"errors"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events chan Event
	Errors chan error
}

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*Watcher, error) {
	return nil, errors.New("FEN based watcher not yet supported for fsnotify\n")
}

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() error {
	return nil
}

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error {
	return nil
}

// Remove stops watching the the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error {
	return nil
}
