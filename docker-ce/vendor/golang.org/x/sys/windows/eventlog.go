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

package windows

const (
	EVENTLOG_SUCCESS          = 0
	EVENTLOG_ERROR_TYPE       = 1
	EVENTLOG_WARNING_TYPE     = 2
	EVENTLOG_INFORMATION_TYPE = 4
	EVENTLOG_AUDIT_SUCCESS    = 8
	EVENTLOG_AUDIT_FAILURE    = 16
)

//sys	RegisterEventSource(uncServerName *uint16, sourceName *uint16) (handle Handle, err error) [failretval==0] = advapi32.RegisterEventSourceW
//sys	DeregisterEventSource(handle Handle) (err error) = advapi32.DeregisterEventSource
//sys	ReportEvent(log Handle, etype uint16, category uint16, eventId uint32, usrSId uintptr, numStrings uint16, dataSize uint32, strings **uint16, rawData *byte) (err error) = advapi32.ReportEventW
