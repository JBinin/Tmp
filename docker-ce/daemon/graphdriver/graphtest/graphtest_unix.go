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
// +build linux freebsd solaris

package graphtest

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"reflect"
	"syscall"
	"testing"
	"unsafe"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/go-units"
)

var (
	drv *Driver
)

// Driver conforms to graphdriver.Driver interface and
// contains information such as root and reference count of the number of clients using it.
// This helps in testing drivers added into the framework.
type Driver struct {
	graphdriver.Driver
	root     string
	refCount int
}

func newDriver(t testing.TB, name string, options []string) *Driver {
	root, err := ioutil.TempDir("", "docker-graphtest-")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}

	d, err := graphdriver.GetDriver(name, nil, graphdriver.Options{DriverOptions: options, Root: root})
	if err != nil {
		t.Logf("graphdriver: %v\n", err)
		if err == graphdriver.ErrNotSupported || err == graphdriver.ErrPrerequisites || err == graphdriver.ErrIncompatibleFS {
			t.Skipf("Driver %s not supported", name)
		}
		t.Fatal(err)
	}
	return &Driver{d, root, 1}
}

func cleanup(t testing.TB, d *Driver) {
	if err := drv.Cleanup(); err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(d.root)
}

// GetDriver create a new driver with given name or return an existing driver with the name updating the reference count.
func GetDriver(t testing.TB, name string, options ...string) graphdriver.Driver {
	if drv == nil {
		drv = newDriver(t, name, options)
	} else {
		drv.refCount++
	}
	return drv
}

// PutDriver removes the driver if it is no longer used and updates the reference count.
func PutDriver(t testing.TB) {
	if drv == nil {
		t.Skip("No driver to put!")
	}
	drv.refCount--
	if drv.refCount == 0 {
		cleanup(t, drv)
		drv = nil
	}
}

// DriverTestCreateEmpty creates a new image and verifies it is empty and the right metadata
func DriverTestCreateEmpty(t testing.TB, drivername string, driverOptions ...string) {
	driver := GetDriver(t, drivername, driverOptions...)
	defer PutDriver(t)

	if err := driver.Create("empty", "", nil); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := driver.Remove("empty"); err != nil {
			t.Fatal(err)
		}
	}()

	if !driver.Exists("empty") {
		t.Fatal("Newly created image doesn't exist")
	}

	dir, err := driver.Get("empty", "")
	if err != nil {
		t.Fatal(err)
	}

	verifyFile(t, dir, 0755|os.ModeDir, 0, 0)

	// Verify that the directory is empty
	fis, err := readDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(fis) != 0 {
		t.Fatal("New directory not empty")
	}

	driver.Put("empty")
}

// DriverTestCreateBase create a base driver and verify.
func DriverTestCreateBase(t testing.TB, drivername string, driverOptions ...string) {
	driver := GetDriver(t, drivername, driverOptions...)
	defer PutDriver(t)

	createBase(t, driver, "Base")
	defer func() {
		if err := driver.Remove("Base"); err != nil {
			t.Fatal(err)
		}
	}()
	verifyBase(t, driver, "Base")
}

// DriverTestCreateSnap Create a driver and snap and verify.
func DriverTestCreateSnap(t testing.TB, drivername string, driverOptions ...string) {
	driver := GetDriver(t, drivername, driverOptions...)
	defer PutDriver(t)

	createBase(t, driver, "Base")

	defer func() {
		if err := driver.Remove("Base"); err != nil {
			t.Fatal(err)
		}
	}()

	if err := driver.Create("Snap", "Base", nil); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := driver.Remove("Snap"); err != nil {
			t.Fatal(err)
		}
	}()

	verifyBase(t, driver, "Snap")
}

// DriverTestDeepLayerRead reads a file from a lower layer under a given number of layers
func DriverTestDeepLayerRead(t testing.TB, layerCount int, drivername string, driverOptions ...string) {
	driver := GetDriver(t, drivername, driverOptions...)
	defer PutDriver(t)

	base := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil {
		t.Fatal(err)
	}

	content := []byte("test content")
	if err := addFile(driver, base, "testfile.txt", content); err != nil {
		t.Fatal(err)
	}

	topLayer, err := addManyLayers(driver, base, layerCount)
	if err != nil {
		t.Fatal(err)
	}

	err = checkManyLayers(driver, topLayer, layerCount)
	if err != nil {
		t.Fatal(err)
	}

	if err := checkFile(driver, topLayer, "testfile.txt", content); err != nil {
		t.Fatal(err)
	}
}

// DriverTestDiffApply tests diffing and applying produces the same layer
func DriverTestDiffApply(t testing.TB, fileCount int, drivername string, driverOptions ...string) {
	driver := GetDriver(t, drivername, driverOptions...)
	defer PutDriver(t)
	base := stringid.GenerateRandomID()
	upper := stringid.GenerateRandomID()
	deleteFile := "file-remove.txt"
	deleteFileContent := []byte("This file should get removed in upper!")
	deleteDir := "var/lib"

	if err := driver.Create(base, "", nil); err != nil {
		t.Fatal(err)
	}

	if err := addManyFiles(driver, base, fileCount, 3); err != nil {
		t.Fatal(err)
	}

	if err := addFile(driver, base, deleteFile, deleteFileContent); err != nil {
		t.Fatal(err)
	}

	if err := addDirectory(driver, base, deleteDir); err != nil {
		t.Fatal(err)
	}

	if err := driver.Create(upper, base, nil); err != nil {
		t.Fatal(err)
	}

	if err := addManyFiles(driver, upper, fileCount, 6); err != nil {
		t.Fatal(err)
	}

	if err := removeAll(driver, upper, deleteFile, deleteDir); err != nil {
		t.Fatal(err)
	}

	diffSize, err := driver.DiffSize(upper, "")
	if err != nil {
		t.Fatal(err)
	}

	diff := stringid.GenerateRandomID()
	if err := driver.Create(diff, base, nil); err != nil {
		t.Fatal(err)
	}

	if err := checkManyFiles(driver, diff, fileCount, 3); err != nil {
		t.Fatal(err)
	}

	if err := checkFile(driver, diff, deleteFile, deleteFileContent); err != nil {
		t.Fatal(err)
	}

	arch, err := driver.Diff(upper, base)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(arch); err != nil {
		t.Fatal(err)
	}
	if err := arch.Close(); err != nil {
		t.Fatal(err)
	}

	applyDiffSize, err := driver.ApplyDiff(diff, base, bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	if applyDiffSize != diffSize {
		t.Fatalf("Apply diff size different, got %d, expected %d", applyDiffSize, diffSize)
	}

	if err := checkManyFiles(driver, diff, fileCount, 6); err != nil {
		t.Fatal(err)
	}

	if err := checkFileRemoved(driver, diff, deleteFile); err != nil {
		t.Fatal(err)
	}

	if err := checkFileRemoved(driver, diff, deleteDir); err != nil {
		t.Fatal(err)
	}
}

// DriverTestChanges tests computed changes on a layer matches changes made
func DriverTestChanges(t testing.TB, drivername string, driverOptions ...string) {
	driver := GetDriver(t, drivername, driverOptions...)
	defer PutDriver(t)
	base := stringid.GenerateRandomID()
	upper := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil {
		t.Fatal(err)
	}

	if err := addManyFiles(driver, base, 20, 3); err != nil {
		t.Fatal(err)
	}

	if err := driver.Create(upper, base, nil); err != nil {
		t.Fatal(err)
	}

	expectedChanges, err := changeManyFiles(driver, upper, 20, 6)
	if err != nil {
		t.Fatal(err)
	}

	changes, err := driver.Changes(upper, base)
	if err != nil {
		t.Fatal(err)
	}

	if err = checkChanges(expectedChanges, changes); err != nil {
		t.Fatal(err)
	}
}

func writeRandomFile(path string, size uint64) error {
	buf := make([]int64, size/8)

	r := rand.NewSource(0)
	for i := range buf {
		buf[i] = r.Int63()
	}

	// Cast to []byte
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&buf))
	header.Len *= 8
	header.Cap *= 8
	data := *(*[]byte)(unsafe.Pointer(&header))

	return ioutil.WriteFile(path, data, 0700)
}

// DriverTestSetQuota Create a driver and test setting quota.
func DriverTestSetQuota(t *testing.T, drivername string) {
	driver := GetDriver(t, drivername)
	defer PutDriver(t)

	createBase(t, driver, "Base")
	createOpts := &graphdriver.CreateOpts{}
	createOpts.StorageOpt = make(map[string]string, 1)
	createOpts.StorageOpt["size"] = "50M"
	if err := driver.Create("zfsTest", "Base", createOpts); err != nil {
		t.Fatal(err)
	}

	mountPath, err := driver.Get("zfsTest", "")
	if err != nil {
		t.Fatal(err)
	}

	quota := uint64(50 * units.MiB)
	err = writeRandomFile(path.Join(mountPath, "file"), quota*2)
	if pathError, ok := err.(*os.PathError); ok && pathError.Err != syscall.EDQUOT {
		t.Fatalf("expect write() to fail with %v, got %v", syscall.EDQUOT, err)
	}

}
