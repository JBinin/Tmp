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
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build s390x,go1.11,!gccgo,!appengine

package poly1305

// hasVectorFacility reports whether the machine supports
// the vector facility (vx).
func hasVectorFacility() bool

// hasVMSLFacility reports whether the machine supports
// Vector Multiply Sum Logical (VMSL).
func hasVMSLFacility() bool

var hasVX = hasVectorFacility()
var hasVMSL = hasVMSLFacility()

// poly1305vx is an assembly implementation of Poly1305 that uses vector
// instructions. It must only be called if the vector facility (vx) is
// available.
//go:noescape
func poly1305vx(out *[16]byte, m *byte, mlen uint64, key *[32]byte)

// poly1305vmsl is an assembly implementation of Poly1305 that uses vector
// instructions, including VMSL. It must only be called if the vector facility (vx) is
// available and if VMSL is supported.
//go:noescape
func poly1305vmsl(out *[16]byte, m *byte, mlen uint64, key *[32]byte)

// Sum generates an authenticator for m using a one-time key and puts the
// 16-byte result into out. Authenticating two different messages with the same
// key allows an attacker to forge messages at will.
func Sum(out *[16]byte, m []byte, key *[32]byte) {
	if hasVX {
		var mPtr *byte
		if len(m) > 0 {
			mPtr = &m[0]
		}
		if hasVMSL && len(m) > 256 {
			poly1305vmsl(out, mPtr, uint64(len(m)), key)
		} else {
			poly1305vx(out, mPtr, uint64(len(m)), key)
		}
	} else {
		sumGeneric(out, m, key)
	}
}
