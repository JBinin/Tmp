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
// This file is included to the build if any of the buildtags below
// are defined. Refer to README notes for more details.

//+build easyjson_nounsafe appengine

package jlexer

// bytesToStr creates a string normally from []byte
//
// Note that this method is roughly 1.5x slower than using the 'unsafe' method.
func bytesToStr(data []byte) string {
	return string(data)
}
