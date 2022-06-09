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
package main

import (
	"sync/atomic"

	_ "github.com/docker/docker/autogen/winresources/docker"
)

//go:cgo_import_dynamic main.dummy CommandLineToArgvW%2 "shell32.dll"

var dummy uintptr

func init() {
	// Ensure that this import is not removed by the linker. This is used to
	// ensure that shell32.dll is loaded by the system loader, preventing
	// go#15286 from triggering on Nano Server TP5.
	atomic.LoadUintptr(&dummy)
}
