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
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

// Message authentication support

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

type macMode struct {
	keySize int
	etm     bool
	new     func(key []byte) hash.Hash
}

// truncatingMAC wraps around a hash.Hash and truncates the output digest to
// a given size.
type truncatingMAC struct {
	length int
	hmac   hash.Hash
}

func (t truncatingMAC) Write(data []byte) (int, error) {
	return t.hmac.Write(data)
}

func (t truncatingMAC) Sum(in []byte) []byte {
	out := t.hmac.Sum(in)
	return out[:len(in)+t.length]
}

func (t truncatingMAC) Reset() {
	t.hmac.Reset()
}

func (t truncatingMAC) Size() int {
	return t.length
}

func (t truncatingMAC) BlockSize() int { return t.hmac.BlockSize() }

var macModes = map[string]*macMode{
	"hmac-sha2-256-etm@openssh.com": {32, true, func(key []byte) hash.Hash {
		return hmac.New(sha256.New, key)
	}},
	"hmac-sha2-256": {32, false, func(key []byte) hash.Hash {
		return hmac.New(sha256.New, key)
	}},
	"hmac-sha1": {20, false, func(key []byte) hash.Hash {
		return hmac.New(sha1.New, key)
	}},
	"hmac-sha1-96": {20, false, func(key []byte) hash.Hash {
		return truncatingMAC{12, hmac.New(sha1.New, key)}
	}},
}
