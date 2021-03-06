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

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package unix

// TimespecToNsec converts a Timespec value into a number of
// nanoseconds since the Unix epoch.
func TimespecToNsec(ts Timespec) int64 { return int64(ts.Sec)*1e9 + int64(ts.Nsec) }

// NsecToTimespec takes a number of nanoseconds since the Unix epoch
// and returns the corresponding Timespec value.
func NsecToTimespec(nsec int64) Timespec {
	sec := nsec / 1e9
	nsec = nsec % 1e9
	if nsec < 0 {
		nsec += 1e9
		sec--
	}
	return setTimespec(sec, nsec)
}

// TimevalToNsec converts a Timeval value into a number of nanoseconds
// since the Unix epoch.
func TimevalToNsec(tv Timeval) int64 { return int64(tv.Sec)*1e9 + int64(tv.Usec)*1e3 }

// NsecToTimeval takes a number of nanoseconds since the Unix epoch
// and returns the corresponding Timeval value.
func NsecToTimeval(nsec int64) Timeval {
	nsec += 999 // round up to microsecond
	usec := nsec % 1e9 / 1e3
	sec := nsec / 1e9
	if usec < 0 {
		usec += 1e6
		sec--
	}
	return setTimeval(sec, usec)
}

// Unix returns ts as the number of seconds and nanoseconds elapsed since the
// Unix epoch.
func (ts *Timespec) Unix() (sec int64, nsec int64) {
	return int64(ts.Sec), int64(ts.Nsec)
}

// Unix returns tv as the number of seconds and nanoseconds elapsed since the
// Unix epoch.
func (tv *Timeval) Unix() (sec int64, nsec int64) {
	return int64(tv.Sec), int64(tv.Usec) * 1000
}

// Nano returns ts as the number of nanoseconds elapsed since the Unix epoch.
func (ts *Timespec) Nano() int64 {
	return int64(ts.Sec)*1e9 + int64(ts.Nsec)
}

// Nano returns tv as the number of nanoseconds elapsed since the Unix epoch.
func (tv *Timeval) Nano() int64 {
	return int64(tv.Sec)*1e9 + int64(tv.Usec)*1000
}
