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
package useragent

import "testing"

func TestVersionInfo(t *testing.T) {
	vi := VersionInfo{"foo", "bar"}
	if !vi.isValid() {
		t.Fatalf("VersionInfo should be valid")
	}
	vi = VersionInfo{"", "bar"}
	if vi.isValid() {
		t.Fatalf("Expected VersionInfo to be invalid")
	}
	vi = VersionInfo{"foo", ""}
	if vi.isValid() {
		t.Fatalf("Expected VersionInfo to be invalid")
	}
}

func TestAppendVersions(t *testing.T) {
	vis := []VersionInfo{
		{"foo", "1.0"},
		{"bar", "0.1"},
		{"pi", "3.1.4"},
	}
	v := AppendVersions("base", vis...)
	expect := "base foo/1.0 bar/0.1 pi/3.1.4"
	if v != expect {
		t.Fatalf("expected %q, got %q", expect, v)
	}
}
