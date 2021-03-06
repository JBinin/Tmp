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

package devmapper

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/graphtest"
)

func init() {
	// Reduce the size of the base fs and loopback for the tests
	defaultDataLoopbackSize = 300 * 1024 * 1024
	defaultMetaDataLoopbackSize = 200 * 1024 * 1024
	defaultBaseFsSize = 300 * 1024 * 1024
	defaultUdevSyncOverride = true
	if err := graphtest.InitLoopbacks(); err != nil {
		panic(err)
	}
}

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestDevmapperSetup and TestDevmapperTeardown
func TestDevmapperSetup(t *testing.T) {
	graphtest.GetDriver(t, "devicemapper")
}

func TestDevmapperCreateEmpty(t *testing.T) {
	graphtest.DriverTestCreateEmpty(t, "devicemapper")
}

func TestDevmapperCreateBase(t *testing.T) {
	graphtest.DriverTestCreateBase(t, "devicemapper")
}

func TestDevmapperCreateSnap(t *testing.T) {
	graphtest.DriverTestCreateSnap(t, "devicemapper")
}

func TestDevmapperTeardown(t *testing.T) {
	graphtest.PutDriver(t)
}

func TestDevmapperReduceLoopBackSize(t *testing.T) {
	tenMB := int64(10 * 1024 * 1024)
	testChangeLoopBackSize(t, -tenMB, defaultDataLoopbackSize, defaultMetaDataLoopbackSize)
}

func TestDevmapperIncreaseLoopBackSize(t *testing.T) {
	tenMB := int64(10 * 1024 * 1024)
	testChangeLoopBackSize(t, tenMB, defaultDataLoopbackSize+tenMB, defaultMetaDataLoopbackSize+tenMB)
}

func testChangeLoopBackSize(t *testing.T, delta, expectDataSize, expectMetaDataSize int64) {
	driver := graphtest.GetDriver(t, "devicemapper").(*graphtest.Driver).Driver.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	defer graphtest.PutDriver(t)
	// make sure data or metadata loopback size are the default size
	if s := driver.DeviceSet.Status(); s.Data.Total != uint64(defaultDataLoopbackSize) || s.Metadata.Total != uint64(defaultMetaDataLoopbackSize) {
		t.Fatalf("data or metadata loop back size is incorrect")
	}
	if err := driver.Cleanup(); err != nil {
		t.Fatal(err)
	}
	//Reload
	d, err := Init(driver.home, []string{
		fmt.Sprintf("dm.loopdatasize=%d", defaultDataLoopbackSize+delta),
		fmt.Sprintf("dm.loopmetadatasize=%d", defaultMetaDataLoopbackSize+delta),
	}, nil, nil)
	if err != nil {
		t.Fatalf("error creating devicemapper driver: %v", err)
	}
	driver = d.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	if s := driver.DeviceSet.Status(); s.Data.Total != uint64(expectDataSize) || s.Metadata.Total != uint64(expectMetaDataSize) {
		t.Fatalf("data or metadata loop back size is incorrect")
	}
	if err := driver.Cleanup(); err != nil {
		t.Fatal(err)
	}
}

// Make sure devices.Lock() has been release upon return from cleanupDeletedDevices() function
func TestDevmapperLockReleasedDeviceDeletion(t *testing.T) {
	driver := graphtest.GetDriver(t, "devicemapper").(*graphtest.Driver).Driver.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	defer graphtest.PutDriver(t)

	// Call cleanupDeletedDevices() and after the call take and release
	// DeviceSet Lock. If lock has not been released, this will hang.
	driver.DeviceSet.cleanupDeletedDevices()

	doneChan := make(chan bool)

	go func() {
		driver.DeviceSet.Lock()
		defer driver.DeviceSet.Unlock()
		doneChan <- true
	}()

	select {
	case <-time.After(time.Second * 5):
		// Timer expired. That means lock was not released upon
		// function return and we are deadlocked. Release lock
		// here so that cleanup could succeed and fail the test.
		driver.DeviceSet.Unlock()
		t.Fatalf("Could not acquire devices lock after call to cleanupDeletedDevices()")
	case <-doneChan:
	}
}
