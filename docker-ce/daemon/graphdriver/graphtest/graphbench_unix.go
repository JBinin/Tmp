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

package graphtest

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/stringid"
)

// DriverBenchExists benchmarks calls to exist
func DriverBenchExists(b *testing.B, drivername string, driveroptions ...string) {
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()

	if err := driver.Create(base, "", nil); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !driver.Exists(base) {
			b.Fatal("Newly created image doesn't exist")
		}
	}
}

// DriverBenchGetEmpty benchmarks calls to get on an empty layer
func DriverBenchGetEmpty(b *testing.B, drivername string, driveroptions ...string) {
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()

	if err := driver.Create(base, "", nil); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := driver.Get(base, "")
		b.StopTimer()
		if err != nil {
			b.Fatalf("Error getting mount: %s", err)
		}
		if err := driver.Put(base); err != nil {
			b.Fatalf("Error putting mount: %s", err)
		}
		b.StartTimer()
	}
}

// DriverBenchDiffBase benchmarks calls to diff on a root layer
func DriverBenchDiffBase(b *testing.B, drivername string, driveroptions ...string) {
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil {
		b.Fatal(err)
	}

	if err := addFiles(driver, base, 3); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arch, err := driver.Diff(base, "")
		if err != nil {
			b.Fatal(err)
		}
		_, err = io.Copy(ioutil.Discard, arch)
		if err != nil {
			b.Fatalf("Error copying archive: %s", err)
		}
		arch.Close()
	}
}

// DriverBenchDiffN benchmarks calls to diff on two layers with
// a provided number of files on the lower and upper layers.
func DriverBenchDiffN(b *testing.B, bottom, top int, drivername string, driveroptions ...string) {
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)
	base := stringid.GenerateRandomID()
	upper := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil {
		b.Fatal(err)
	}

	if err := addManyFiles(driver, base, bottom, 3); err != nil {
		b.Fatal(err)
	}

	if err := driver.Create(upper, base, nil); err != nil {
		b.Fatal(err)
	}

	if err := addManyFiles(driver, upper, top, 6); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arch, err := driver.Diff(upper, "")
		if err != nil {
			b.Fatal(err)
		}
		_, err = io.Copy(ioutil.Discard, arch)
		if err != nil {
			b.Fatalf("Error copying archive: %s", err)
		}
		arch.Close()
	}
}

// DriverBenchDiffApplyN benchmarks calls to diff and apply together
func DriverBenchDiffApplyN(b *testing.B, fileCount int, drivername string, driveroptions ...string) {
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)
	base := stringid.GenerateRandomID()
	upper := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil {
		b.Fatal(err)
	}

	if err := addManyFiles(driver, base, fileCount, 3); err != nil {
		b.Fatal(err)
	}

	if err := driver.Create(upper, base, nil); err != nil {
		b.Fatal(err)
	}

	if err := addManyFiles(driver, upper, fileCount, 6); err != nil {
		b.Fatal(err)
	}
	diffSize, err := driver.DiffSize(upper, "")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		diff := stringid.GenerateRandomID()
		if err := driver.Create(diff, base, nil); err != nil {
			b.Fatal(err)
		}

		if err := checkManyFiles(driver, diff, fileCount, 3); err != nil {
			b.Fatal(err)
		}

		b.StartTimer()

		arch, err := driver.Diff(upper, "")
		if err != nil {
			b.Fatal(err)
		}

		applyDiffSize, err := driver.ApplyDiff(diff, "", arch)
		if err != nil {
			b.Fatal(err)
		}

		b.StopTimer()
		arch.Close()

		if applyDiffSize != diffSize {
			// TODO: enforce this
			//b.Fatalf("Apply diff size different, got %d, expected %s", applyDiffSize, diffSize)
		}
		if err := checkManyFiles(driver, diff, fileCount, 6); err != nil {
			b.Fatal(err)
		}
	}
}

// DriverBenchDeepLayerDiff benchmarks calls to diff on top of a given number of layers.
func DriverBenchDeepLayerDiff(b *testing.B, layerCount int, drivername string, driveroptions ...string) {
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil {
		b.Fatal(err)
	}

	if err := addFiles(driver, base, 50); err != nil {
		b.Fatal(err)
	}

	topLayer, err := addManyLayers(driver, base, layerCount)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arch, err := driver.Diff(topLayer, "")
		if err != nil {
			b.Fatal(err)
		}
		_, err = io.Copy(ioutil.Discard, arch)
		if err != nil {
			b.Fatalf("Error copying archive: %s", err)
		}
		arch.Close()
	}
}

// DriverBenchDeepLayerRead benchmarks calls to read a file under a given number of layers.
func DriverBenchDeepLayerRead(b *testing.B, layerCount int, drivername string, driveroptions ...string) {
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil {
		b.Fatal(err)
	}

	content := []byte("test content")
	if err := addFile(driver, base, "testfile.txt", content); err != nil {
		b.Fatal(err)
	}

	topLayer, err := addManyLayers(driver, base, layerCount)
	if err != nil {
		b.Fatal(err)
	}

	root, err := driver.Get(topLayer, "")
	if err != nil {
		b.Fatal(err)
	}
	defer driver.Put(topLayer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		// Read content
		c, err := ioutil.ReadFile(filepath.Join(root, "testfile.txt"))
		if err != nil {
			b.Fatal(err)
		}

		b.StopTimer()
		if bytes.Compare(c, content) != 0 {
			b.Fatalf("Wrong content in file %v, expected %v", c, content)
		}
		b.StartTimer()
	}
}
