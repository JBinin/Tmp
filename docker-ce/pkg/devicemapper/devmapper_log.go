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
// +build linux

package devicemapper

import "C"

import (
	"strings"
)

// Due to the way cgo works this has to be in a separate file, as devmapper.go has
// definitions in the cgo block, which is incompatible with using "//export"

// DevmapperLogCallback exports the devmapper log callback for cgo.
//export DevmapperLogCallback
func DevmapperLogCallback(level C.int, file *C.char, line C.int, dmErrnoOrClass C.int, message *C.char) {
	msg := C.GoString(message)
	if level < 7 {
		if strings.Contains(msg, "busy") {
			dmSawBusy = true
		}

		if strings.Contains(msg, "File exists") {
			dmSawExist = true
		}

		if strings.Contains(msg, "No such device or address") {
			dmSawEnxio = true
		}
	}

	if dmLogger != nil {
		dmLogger.DMLog(int(level), C.GoString(file), int(line), int(dmErrnoOrClass), msg)
	}
}
