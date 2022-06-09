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
// +build !go1.2

package toml

// These interfaces were introduced in Go 1.2, so we add them manually when
// compiling for Go 1.1.

// TextMarshaler is a synonym for encoding.TextMarshaler. It is defined here
// so that Go 1.1 can be supported.
type TextMarshaler interface {
	MarshalText() (text []byte, err error)
}

// TextUnmarshaler is a synonym for encoding.TextUnmarshaler. It is defined
// here so that Go 1.1 can be supported.
type TextUnmarshaler interface {
	UnmarshalText(text []byte) error
}
