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
// +build linux

package fsutils

import (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func testSupportsDType(t *testing.T, expected bool, mkfsCommand string, mkfsArg ...string) {
	// check whether mkfs is installed
	if _, err := exec.LookPath(mkfsCommand); err != nil {
		t.Skipf("%s not installed: %v", mkfsCommand, err)
	}

	// create a sparse image
	imageSize := int64(32 * 1024 * 1024)
	imageFile, err := ioutil.TempFile("", "fsutils-image")
	if err != nil {
		t.Fatal(err)
	}
	imageFileName := imageFile.Name()
	defer os.Remove(imageFileName)
	if _, err = imageFile.Seek(imageSize-1, 0); err != nil {
		t.Fatal(err)
	}
	if _, err = imageFile.Write([]byte{0}); err != nil {
		t.Fatal(err)
	}
	if err = imageFile.Close(); err != nil {
		t.Fatal(err)
	}

	// create a mountpoint
	mountpoint, err := ioutil.TempDir("", "fsutils-mountpoint")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(mountpoint)

	// format the image
	args := append(mkfsArg, imageFileName)
	t.Logf("Executing `%s %v`", mkfsCommand, args)
	out, err := exec.Command(mkfsCommand, args...).CombinedOutput()
	if len(out) > 0 {
		t.Log(string(out))
	}
	if err != nil {
		t.Fatal(err)
	}

	// loopback-mount the image.
	// for ease of setting up loopback device, we use os/exec rather than syscall.Mount
	out, err = exec.Command("mount", "-o", "loop", imageFileName, mountpoint).CombinedOutput()
	if len(out) > 0 {
		t.Log(string(out))
	}
	if err != nil {
		t.Skip("skipping the test because mount failed")
	}
	defer func() {
		if err := syscall.Unmount(mountpoint, 0); err != nil {
			t.Fatal(err)
		}
	}()

	// check whether it supports d_type
	result, err := SupportsDType(mountpoint)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Supports d_type: %v", result)
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestSupportsDTypeWithFType0XFS(t *testing.T) {
	testSupportsDType(t, false, "mkfs.xfs", "-m", "crc=0", "-n", "ftype=0")
}

func TestSupportsDTypeWithFType1XFS(t *testing.T) {
	testSupportsDType(t, true, "mkfs.xfs", "-m", "crc=0", "-n", "ftype=1")
}

func TestSupportsDTypeWithExt4(t *testing.T) {
	testSupportsDType(t, true, "mkfs.ext4")
}
