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
package pidfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewAndRemove(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test-pidfile")
	if err != nil {
		t.Fatal("Could not create test directory")
	}

	path := filepath.Join(dir, "testfile")
	file, err := New(path)
	if err != nil {
		t.Fatal("Could not create test file", err)
	}

	_, err = New(path)
	if err == nil {
		t.Fatal("Test file creation not blocked")
	}

	if err := file.Remove(); err != nil {
		t.Fatal("Could not delete created test file")
	}
}

func TestRemoveInvalidPath(t *testing.T) {
	file := PIDFile{path: filepath.Join("foo", "bar")}

	if err := file.Remove(); err == nil {
		t.Fatal("Non-existing file doesn't give an error on delete")
	}
}
