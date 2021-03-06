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
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Functions to access/create device major and minor numbers matching the
// encoding used by the Linux kernel and glibc.
//
// The information below is extracted and adapted from bits/sysmacros.h in the
// glibc sources:
//
// dev_t in glibc is 64-bit, with 32-bit major and minor numbers. glibc's
// default encoding is MMMM Mmmm mmmM MMmm, where M is a hex digit of the major
// number and m is a hex digit of the minor number. This is backward compatible
// with legacy systems where dev_t is 16 bits wide, encoded as MMmm. It is also
// backward compatible with the Linux kernel, which for some architectures uses
// 32-bit dev_t, encoded as mmmM MMmm.

package unix

// Major returns the major component of a Linux device number.
func Major(dev uint64) uint32 {
	major := uint32((dev & 0x00000000000fff00) >> 8)
	major |= uint32((dev & 0xfffff00000000000) >> 32)
	return major
}

// Minor returns the minor component of a Linux device number.
func Minor(dev uint64) uint32 {
	minor := uint32((dev & 0x00000000000000ff) >> 0)
	minor |= uint32((dev & 0x00000ffffff00000) >> 12)
	return minor
}

// Mkdev returns a Linux device number generated from the given major and minor
// components.
func Mkdev(major, minor uint32) uint64 {
	dev := (uint64(major) & 0x00000fff) << 8
	dev |= (uint64(major) & 0xfffff000) << 32
	dev |= (uint64(minor) & 0x000000ff) << 0
	dev |= (uint64(minor) & 0xffffff00) << 12
	return dev
}
