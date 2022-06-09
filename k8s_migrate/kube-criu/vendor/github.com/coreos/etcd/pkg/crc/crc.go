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
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crc provides utility function for cyclic redundancy check
// algorithms.
package crc

import (
	"hash"
	"hash/crc32"
)

// The size of a CRC-32 checksum in bytes.
const Size = 4

type digest struct {
	crc uint32
	tab *crc32.Table
}

// New creates a new hash.Hash32 computing the CRC-32 checksum
// using the polynomial represented by the Table.
// Modified by xiangli to take a prevcrc.
func New(prev uint32, tab *crc32.Table) hash.Hash32 { return &digest{prev, tab} }

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Reset() { d.crc = 0 }

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = crc32.Update(d.crc, d.tab, p)
	return len(p), nil
}

func (d *digest) Sum32() uint32 { return d.crc }

func (d *digest) Sum(in []byte) []byte {
	s := d.Sum32()
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}
