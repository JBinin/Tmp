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
package homedir

import (
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) {
	home := Get()
	if home == "" {
		t.Fatal("returned home directory is empty")
	}

	if !filepath.IsAbs(home) {
		t.Fatalf("returned path is not absolute: %s", home)
	}
}

func TestGetShortcutString(t *testing.T) {
	shortcut := GetShortcutString()
	if shortcut == "" {
		t.Fatal("returned shortcut string is empty")
	}
}
