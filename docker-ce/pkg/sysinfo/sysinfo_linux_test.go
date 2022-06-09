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
package sysinfo

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestReadProcBool(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-sysinfo-proc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	procFile := filepath.Join(tmpDir, "read-proc-bool")
	if err := ioutil.WriteFile(procFile, []byte("1"), 644); err != nil {
		t.Fatal(err)
	}

	if !readProcBool(procFile) {
		t.Fatal("expected proc bool to be true, got false")
	}

	if err := ioutil.WriteFile(procFile, []byte("0"), 644); err != nil {
		t.Fatal(err)
	}
	if readProcBool(procFile) {
		t.Fatal("expected proc bool to be false, got false")
	}

	if readProcBool(path.Join(tmpDir, "no-exist")) {
		t.Fatal("should be false for non-existent entry")
	}

}

func TestCgroupEnabled(t *testing.T) {
	cgroupDir, err := ioutil.TempDir("", "cgroup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cgroupDir)

	if cgroupEnabled(cgroupDir, "test") {
		t.Fatal("cgroupEnabled should be false")
	}

	if err := ioutil.WriteFile(path.Join(cgroupDir, "test"), []byte{}, 644); err != nil {
		t.Fatal(err)
	}

	if !cgroupEnabled(cgroupDir, "test") {
		t.Fatal("cgroupEnabled should be true")
	}
}
