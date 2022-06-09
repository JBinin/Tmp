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
// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package monotime provides a fast monotonic clock source.
package monotime

import (
	_ "unsafe" // required to use //go:linkname
)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// Now returns the current time in nanoseconds from a monotonic clock.
// The time returned is based on some arbitrary platform-specific point in the
// past.  The time returned is guaranteed to increase monotonically at a
// constant rate, unlike time.Now() from the Go standard library, which may
// slow down, speed up, jump forward or backward, due to NTP activity or leap
// seconds.
func Now() Time {
	return Time(nanotime())
}
