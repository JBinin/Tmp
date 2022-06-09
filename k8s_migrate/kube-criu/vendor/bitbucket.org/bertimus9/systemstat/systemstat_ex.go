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
// Copyright (c) 2013 Phillip Bond
// Licensed under the MIT License
// see file LICENSE

// +build !linux

package systemstat

import (
	"time"
)

func getUptime(procfile string) (uptime UptimeSample) {
	notImplemented("getUptime")
	uptime.Time = time.Now()
	return
}

func getLoadAvgSample(procfile string) (samp LoadAvgSample) {
	notImplemented("getLoadAvgSample")
	samp.Time = time.Now()
	return
}

func getMemSample(procfile string) (samp MemSample) {
	notImplemented("getMemSample")
	samp.Time = time.Now()
	return
}

func getProcCPUSample() (s ProcCPUSample) {
	notImplemented("getProcCPUSample")
	s.Time = time.Now()
	return
}

func getCPUSample(procfile string) (samp CPUSample) {
	notImplemented("getCPUSample")
	samp.Time = time.Now()
	return
}
