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
package ioutils

import "testing"

func TestFprintfIfNotEmpty(t *testing.T) {
	wc := NewWriteCounter(&NopWriter{})
	n, _ := FprintfIfNotEmpty(wc, "foo%s", "")

	if wc.Count != 0 || n != 0 {
		t.Errorf("Wrong count: %v vs. %v vs. 0", wc.Count, n)
	}

	n, _ = FprintfIfNotEmpty(wc, "foo%s", "bar")
	if wc.Count != 6 || n != 6 {
		t.Errorf("Wrong count: %v vs. %v vs. 6", wc.Count, n)
	}
}
