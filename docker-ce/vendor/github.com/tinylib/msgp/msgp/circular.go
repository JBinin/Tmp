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
package msgp

import (
	"testing"
)

// EndlessReader is an io.Reader
// that loops over the same data
// endlessly. It is used for benchmarking.
type EndlessReader struct {
	tb     *testing.B
	data   []byte
	offset int
}

// NewEndlessReader returns a new endless reader
func NewEndlessReader(b []byte, tb *testing.B) *EndlessReader {
	return &EndlessReader{tb: tb, data: b, offset: 0}
}

// Read implements io.Reader. In practice, it
// always returns (len(p), nil), although it
// fills the supplied slice while the benchmark
// timer is stopped.
func (c *EndlessReader) Read(p []byte) (int, error) {
	c.tb.StopTimer()
	var n int
	l := len(p)
	m := len(c.data)
	for n < l {
		nn := copy(p[n:], c.data[c.offset:])
		n += nn
		c.offset += nn
		c.offset %= m
	}
	c.tb.StartTimer()
	return n, nil
}
