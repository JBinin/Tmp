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
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import "golang.org/x/text/transform"

// Transformer implements the transform.Transformer interface.
type Transformer struct {
	t transform.Transformer
}

// Reset implements the transform.Transformer interface.
func (t Transformer) Reset() { t.t.Reset() }

// Transform implements the transform.Transformer interface.
func (t Transformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	return t.t.Transform(dst, src, atEOF)
}

// Bytes returns a new byte slice with the result of applying t to b.
func (t Transformer) Bytes(b []byte) []byte {
	b, _, _ = transform.Bytes(t, b)
	return b
}

// String returns a string with the result of applying t to s.
func (t Transformer) String(s string) string {
	s, _, _ = transform.String(t, s)
	return s
}
