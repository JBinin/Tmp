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
// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
)

var (
	newline    = []byte{'\n'}
	dataindent = []byte{'\t', '\t'}
	space      = []byte{' '}
)

type ByteWriter struct {
	io.Writer
	c int
}

func (w *ByteWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	for n = range p {
		if w.c%12 == 0 {
			w.Writer.Write(newline)
			w.Writer.Write(dataindent)
			w.c = 0
		} else {
			w.Writer.Write(space)
		}

		fmt.Fprintf(w.Writer, "0x%02x,", p[n])
		w.c++
	}

	n++

	return
}
