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
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"errors"
)

// fixedBuffer is an io.ReadWriter backed by a fixed size buffer.
// It never allocates, but moves old data as new data is written.
type fixedBuffer struct {
	buf  []byte
	r, w int
}

var (
	errReadEmpty = errors.New("read from empty fixedBuffer")
	errWriteFull = errors.New("write on full fixedBuffer")
)

// Read copies bytes from the buffer into p.
// It is an error to read when no data is available.
func (b *fixedBuffer) Read(p []byte) (n int, err error) {
	if b.r == b.w {
		return 0, errReadEmpty
	}
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	if b.r == b.w {
		b.r = 0
		b.w = 0
	}
	return n, nil
}

// Len returns the number of bytes of the unread portion of the buffer.
func (b *fixedBuffer) Len() int {
	return b.w - b.r
}

// Write copies bytes from p into the buffer.
// It is an error to write more data than the buffer can hold.
func (b *fixedBuffer) Write(p []byte) (n int, err error) {
	// Slide existing data to beginning.
	if b.r > 0 && len(p) > len(b.buf)-b.w {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	// Write new data.
	n = copy(b.buf[b.w:], p)
	b.w += n
	if n < len(p) {
		err = errWriteFull
	}
	return n, err
}
