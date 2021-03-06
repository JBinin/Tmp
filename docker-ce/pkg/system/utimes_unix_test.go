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
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// prepareFiles creates files for testing in the temp directory
func prepareFiles(t *testing.T) (string, string, string, string) {
	dir, err := ioutil.TempDir("", "docker-system-test")
	if err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(dir, "exist")
	if err := ioutil.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	invalid := filepath.Join(dir, "doesnt-exist")

	symlink := filepath.Join(dir, "symlink")
	if err := os.Symlink(file, symlink); err != nil {
		t.Fatal(err)
	}

	return file, invalid, symlink, dir
}

func TestLUtimesNano(t *testing.T) {
	file, invalid, symlink, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	before, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}

	ts := []syscall.Timespec{{Sec: 0, Nsec: 0}, {Sec: 0, Nsec: 0}}
	if err := LUtimesNano(symlink, ts); err != nil {
		t.Fatal(err)
	}

	symlinkInfo, err := os.Lstat(symlink)
	if err != nil {
		t.Fatal(err)
	}
	if before.ModTime().Unix() == symlinkInfo.ModTime().Unix() {
		t.Fatal("The modification time of the symlink should be different")
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if before.ModTime().Unix() != fileInfo.ModTime().Unix() {
		t.Fatal("The modification time of the file should be same")
	}

	if err := LUtimesNano(invalid, ts); err == nil {
		t.Fatal("Doesn't return an error on a non-existing file")
	}
}
