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

package vfs

import (
	"testing"

	"github.com/docker/docker/daemon/graphdriver/graphtest"

	"github.com/docker/docker/pkg/reexec"
)

func init() {
	reexec.Init()
}

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestVfsSetup and TestVfsTeardown
func TestVfsSetup(t *testing.T) {
	graphtest.GetDriver(t, "vfs")
}

func TestVfsCreateEmpty(t *testing.T) {
	graphtest.DriverTestCreateEmpty(t, "vfs")
}

func TestVfsCreateBase(t *testing.T) {
	graphtest.DriverTestCreateBase(t, "vfs")
}

func TestVfsCreateSnap(t *testing.T) {
	graphtest.DriverTestCreateSnap(t, "vfs")
}

func TestVfsTeardown(t *testing.T) {
	graphtest.PutDriver(t)
}
