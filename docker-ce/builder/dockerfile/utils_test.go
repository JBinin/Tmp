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
package dockerfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// createTestTempDir creates a temporary directory for testing.
// It returns the created path and a cleanup function which is meant to be used as deferred call.
// When an error occurs, it terminates the test.
func createTestTempDir(t *testing.T, dir, prefix string) (string, func()) {
	path, err := ioutil.TempDir(dir, prefix)

	if err != nil {
		t.Fatalf("Error when creating directory %s with prefix %s: %s", dir, prefix, err)
	}

	return path, func() {
		err = os.RemoveAll(path)

		if err != nil {
			t.Fatalf("Error when removing directory %s: %s", path, err)
		}
	}
}

// createTestTempFile creates a temporary file within dir with specific contents and permissions.
// When an error occurs, it terminates the test
func createTestTempFile(t *testing.T, dir, filename, contents string, perm os.FileMode) string {
	filePath := filepath.Join(dir, filename)
	err := ioutil.WriteFile(filePath, []byte(contents), perm)

	if err != nil {
		t.Fatalf("Error when creating %s file: %s", filename, err)
	}

	return filePath
}

// createTestSymlink creates a symlink file within dir which points to oldname
func createTestSymlink(t *testing.T, dir, filename, oldname string) string {
	filePath := filepath.Join(dir, filename)
	if err := os.Symlink(oldname, filePath); err != nil {
		t.Fatalf("Error when creating %s symlink to %s: %s", filename, oldname, err)
	}

	return filePath
}
