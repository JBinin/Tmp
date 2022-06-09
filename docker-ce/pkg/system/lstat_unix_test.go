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
// +build linux freebsd

package system

import (
	"os"
	"testing"
)

// TestLstat tests Lstat for existing and non existing files
func TestLstat(t *testing.T) {
	file, invalid, _, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	statFile, err := Lstat(file)
	if err != nil {
		t.Fatal(err)
	}
	if statFile == nil {
		t.Fatal("returned empty stat for existing file")
	}

	statInvalid, err := Lstat(invalid)
	if err == nil {
		t.Fatal("did not return error for non-existing file")
	}
	if statInvalid != nil {
		t.Fatal("returned non-nil stat for non-existing file")
	}
}
