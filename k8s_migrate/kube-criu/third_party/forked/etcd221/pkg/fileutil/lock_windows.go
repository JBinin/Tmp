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
// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build windows

package fileutil

import (
	"errors"
	"os"
)

var (
	ErrLocked = errors.New("file already locked")
)

type Lock interface {
	Name() string
	TryLock() error
	Lock() error
	Unlock() error
	Destroy() error
}

type lock struct {
	fd   int
	file *os.File
}

func (l *lock) Name() string {
	return l.file.Name()
}

// TryLock acquires exclusivity on the lock without blocking
func (l *lock) TryLock() error {
	return nil
}

// Lock acquires exclusivity on the lock without blocking
func (l *lock) Lock() error {
	return nil
}

// Unlock unlocks the lock
func (l *lock) Unlock() error {
	return nil
}

func (l *lock) Destroy() error {
	return l.file.Close()
}

func NewLock(file string) (Lock, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	l := &lock{int(f.Fd()), f}
	return l, nil
}
