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
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package registry

import "syscall"

const (
	_REG_OPTION_NON_VOLATILE = 0

	_REG_CREATED_NEW_KEY     = 1
	_REG_OPENED_EXISTING_KEY = 2

	_ERROR_NO_MORE_ITEMS syscall.Errno = 259
)

func LoadRegLoadMUIString() error {
	return procRegLoadMUIStringW.Find()
}

//sys	regCreateKeyEx(key syscall.Handle, subkey *uint16, reserved uint32, class *uint16, options uint32, desired uint32, sa *syscall.SecurityAttributes, result *syscall.Handle, disposition *uint32) (regerrno error) = advapi32.RegCreateKeyExW
//sys	regDeleteKey(key syscall.Handle, subkey *uint16) (regerrno error) = advapi32.RegDeleteKeyW
//sys	regSetValueEx(key syscall.Handle, valueName *uint16, reserved uint32, vtype uint32, buf *byte, bufsize uint32) (regerrno error) = advapi32.RegSetValueExW
//sys	regEnumValue(key syscall.Handle, index uint32, name *uint16, nameLen *uint32, reserved *uint32, valtype *uint32, buf *byte, buflen *uint32) (regerrno error) = advapi32.RegEnumValueW
//sys	regDeleteValue(key syscall.Handle, name *uint16) (regerrno error) = advapi32.RegDeleteValueW
//sys   regLoadMUIString(key syscall.Handle, name *uint16, buf *uint16, buflen uint32, buflenCopied *uint32, flags uint32, dir *uint16) (regerrno error) = advapi32.RegLoadMUIStringW

//sys	expandEnvironmentStrings(src *uint16, dst *uint16, size uint32) (n uint32, err error) = kernel32.ExpandEnvironmentStringsW
