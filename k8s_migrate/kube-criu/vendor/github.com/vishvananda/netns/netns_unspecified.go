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
// +build !linux

package netns

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

func Set(ns NsHandle) (err error) {
	return ErrNotImplemented
}

func New() (ns NsHandle, err error) {
	return -1, ErrNotImplemented
}

func Get() (NsHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromPath(path string) (NsHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromName(name string) (NsHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromPid(pid int) (NsHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromThread(pid, tid int) (NsHandle, error) {
	return -1, ErrNotImplemented
}

func GetFromDocker(id string) (NsHandle, error) {
	return -1, ErrNotImplemented
}
