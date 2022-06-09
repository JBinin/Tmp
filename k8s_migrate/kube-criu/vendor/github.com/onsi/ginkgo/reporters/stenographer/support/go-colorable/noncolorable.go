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
package colorable

import (
	"bytes"
	"fmt"
	"io"
)

type NonColorable struct {
	out     io.Writer
	lastbuf bytes.Buffer
}

func NewNonColorable(w io.Writer) io.Writer {
	return &NonColorable{out: w}
}

func (w *NonColorable) Write(data []byte) (n int, err error) {
	er := bytes.NewBuffer(data)
loop:
	for {
		c1, _, err := er.ReadRune()
		if err != nil {
			break loop
		}
		if c1 != 0x1b {
			fmt.Fprint(w.out, string(c1))
			continue
		}
		c2, _, err := er.ReadRune()
		if err != nil {
			w.lastbuf.WriteRune(c1)
			break loop
		}
		if c2 != 0x5b {
			w.lastbuf.WriteRune(c1)
			w.lastbuf.WriteRune(c2)
			continue
		}

		var buf bytes.Buffer
		for {
			c, _, err := er.ReadRune()
			if err != nil {
				w.lastbuf.WriteRune(c1)
				w.lastbuf.WriteRune(c2)
				w.lastbuf.Write(buf.Bytes())
				break loop
			}
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' {
				break
			}
			buf.Write([]byte(string(c)))
		}
	}
	return len(data) - w.lastbuf.Len(), nil
}
