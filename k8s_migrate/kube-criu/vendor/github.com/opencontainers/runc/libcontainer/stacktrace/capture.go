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

import "runtime"

// Capture captures a stacktrace for the current calling go program
//
// skip is the number of frames to skip
func Capture(userSkip int) Stacktrace {
	var (
		skip   = userSkip + 1 // add one for our own function
		frames []Frame
		prevPc uintptr
	)
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		//detect if caller is repeated to avoid loop, gccgo
		//currently runs  into a loop without this check
		if !ok || pc == prevPc {
			break
		}
		frames = append(frames, NewFrame(pc, file, line))
		prevPc = pc
	}
	return Stacktrace{
		Frames: frames,
	}
}
