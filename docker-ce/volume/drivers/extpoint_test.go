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
package volumedrivers

import (
	"testing"

	volumetestutils "github.com/docker/docker/volume/testutils"
)

func TestGetDriver(t *testing.T) {
	_, err := GetDriver("missing")
	if err == nil {
		t.Fatal("Expected error, was nil")
	}
	Register(volumetestutils.NewFakeDriver("fake"), "fake")

	d, err := GetDriver("fake")
	if err != nil {
		t.Fatal(err)
	}
	if d.Name() != "fake" {
		t.Fatalf("Expected fake driver, got %s\n", d.Name())
	}
}
