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
// +build !windows

package ct

import (
	"fmt"
	"os"
	"strconv"
)

func isDumbTerm() bool {
	return os.Getenv("TERM") == "dumb"
}

func resetColor() {
	if isDumbTerm() {
		return
	}
	fmt.Print("\x1b[0m")
}

func ansiText(fg Color, fgBright bool, bg Color, bgBright bool) string {
	if fg == None && bg == None {
		return ""
	}
	s := []byte("\x1b[0")
	if fg != None {
		s = strconv.AppendUint(append(s, ";"...), 30+(uint64)(fg-Black), 10)
		if fgBright {
			s = append(s, ";1"...)
		}
	}
	if bg != None {
		s = strconv.AppendUint(append(s, ";"...), 40+(uint64)(bg-Black), 10)
	}
	s = append(s, "m"...)
	return string(s)
}

func changeColor(fg Color, fgBright bool, bg Color, bgBright bool) {
	if isDumbTerm() {
		return
	}
	if fg == None && bg == None {
		return
	}
	fmt.Print(ansiText(fg, fgBright, bg, bgBright))
}
