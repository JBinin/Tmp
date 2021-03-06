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

package gensupport

import (
	"net/url"

	"google.golang.org/api/googleapi"
)

// URLParams is a simplified replacement for url.Values
// that safely builds up URL parameters for encoding.
type URLParams map[string][]string

// Get returns the first value for the given key, or "".
func (u URLParams) Get(key string) string {
	vs := u[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value.
// It replaces any existing values.
func (u URLParams) Set(key, value string) {
	u[key] = []string{value}
}

// SetMulti sets the key to an array of values.
// It replaces any existing values.
// Note that values must not be modified after calling SetMulti
// so the caller is responsible for making a copy if necessary.
func (u URLParams) SetMulti(key string, values []string) {
	u[key] = values
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (u URLParams) Encode() string {
	return url.Values(u).Encode()
}

func SetOptions(u URLParams, opts ...googleapi.CallOption) {
	for _, o := range opts {
		u.Set(o.Get())
	}
}
