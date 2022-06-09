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
// Copied from the golang.org/x/text repo; DO NOT EDIT

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package idna

// appendMapping appends the mapping for the respective rune. isMapped must be
// true. A mapping is a categorization of a rune as defined in UTS #46.
func (c info) appendMapping(b []byte, s string) []byte {
	index := int(c >> indexShift)
	if c&xorBit == 0 {
		s := mappings[index:]
		return append(b, s[1:s[0]+1]...)
	}
	b = append(b, s...)
	if c&inlineXOR == inlineXOR {
		// TODO: support and handle two-byte inline masks
		b[len(b)-1] ^= byte(index)
	} else {
		for p := len(b) - int(xorData[index]); p < len(b); p++ {
			index++
			b[p] ^= xorData[index]
		}
	}
	return b
}

// Sparse block handling code.

type valueRange struct {
	value  uint16 // header: value:stride
	lo, hi byte   // header: lo:n
}

type sparseBlocks struct {
	values []valueRange
	offset []uint16
}

var idnaSparse = sparseBlocks{
	values: idnaSparseValues[:],
	offset: idnaSparseOffset[:],
}

var trie = newIdnaTrie(0)

// lookup determines the type of block n and looks up the value for b.
// For n < t.cutoff, the block is a simple lookup table. Otherwise, the block
// is a list of ranges with an accompanying value. Given a matching range r,
// the value for b is by r.value + (b - r.lo) * stride.
func (t *sparseBlocks) lookup(n uint32, b byte) uint16 {
	offset := t.offset[n]
	header := t.values[offset]
	lo := offset + 1
	hi := lo + uint16(header.lo)
	for lo < hi {
		m := lo + (hi-lo)/2
		r := t.values[m]
		if r.lo <= b && b <= r.hi {
			return r.value + uint16(b-r.lo)*header.value
		}
		if b < r.lo {
			hi = m
		} else {
			lo = m + 1
		}
	}
	return 0
}
