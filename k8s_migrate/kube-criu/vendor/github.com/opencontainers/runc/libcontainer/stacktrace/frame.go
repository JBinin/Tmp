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
package stacktrace

import (
	"path/filepath"
	"runtime"
	"strings"
)

// NewFrame returns a new stack frame for the provided information
func NewFrame(pc uintptr, file string, line int) Frame {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return Frame{}
	}
	pack, name := parseFunctionName(fn.Name())
	return Frame{
		Line:     line,
		File:     filepath.Base(file),
		Package:  pack,
		Function: name,
	}
}

func parseFunctionName(name string) (string, string) {
	i := strings.LastIndex(name, ".")
	if i == -1 {
		return "", name
	}
	return name[:i], name[i+1:]
}

// Frame contains all the information for a stack frame within a go program
type Frame struct {
	File     string
	Function string
	Package  string
	Line     int
}
