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
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,solaris

package unix

func TimespecToNsec(ts Timespec) int64 { return int64(ts.Sec)*1e9 + int64(ts.Nsec) }

func NsecToTimespec(nsec int64) (ts Timespec) {
	ts.Sec = nsec / 1e9
	ts.Nsec = nsec % 1e9
	return
}

func NsecToTimeval(nsec int64) (tv Timeval) {
	nsec += 999 // round up to microsecond
	tv.Usec = nsec % 1e9 / 1e3
	tv.Sec = int64(nsec / 1e9)
	return
}

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint32(length)
}

func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	// TODO(aram): implement this, see issue 5847.
	panic("unimplemented")
}
