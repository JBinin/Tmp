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
// Copyright 2014 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import "errors"

func (u UUID) MarshalJSON() ([]byte, error) {
	if len(u) != 16 {
		return []byte(`""`), nil
	}
	var js [38]byte
	js[0] = '"'
	encodeHex(js[1:], u)
	js[37] = '"'
	return js[:], nil
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	if string(data) == `""` {
		return nil
	}
	if data[0] != '"' {
		return errors.New("invalid UUID format")
	}
	data = data[1 : len(data)-1]
	uu := Parse(string(data))
	if uu == nil {
		return errors.New("invalid UUID format")
	}
	*u = uu
	return nil
}
