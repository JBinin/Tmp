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
// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po

import (
	"io"
	"strings"
)

type lineReader struct {
	lines []string
	pos   int
}

func newLineReader(data string) *lineReader {
	data = strings.Replace(data, "\r", "", -1)
	lines := strings.Split(data, "\n")
	return &lineReader{lines: lines}
}

func (r *lineReader) skipBlankLine() error {
	for ; r.pos < len(r.lines); r.pos++ {
		if strings.TrimSpace(r.lines[r.pos]) != "" {
			break
		}
	}
	if r.pos >= len(r.lines) {
		return io.EOF
	}
	return nil
}

func (r *lineReader) currentPos() int {
	return r.pos
}

func (r *lineReader) currentLine() (s string, pos int, err error) {
	if r.pos >= len(r.lines) {
		err = io.EOF
		return
	}
	s, pos = r.lines[r.pos], r.pos
	return
}

func (r *lineReader) readLine() (s string, pos int, err error) {
	if r.pos >= len(r.lines) {
		err = io.EOF
		return
	}
	s, pos = r.lines[r.pos], r.pos
	r.pos++
	return
}

func (r *lineReader) unreadLine() {
	if r.pos >= 0 {
		r.pos--
	}
}
