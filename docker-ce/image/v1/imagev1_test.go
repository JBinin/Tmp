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
package v1

import (
	"encoding/json"
	"testing"

	"github.com/docker/docker/image"
)

func TestMakeV1ConfigFromConfig(t *testing.T) {
	img := &image.Image{
		V1Image: image.V1Image{
			ID:     "v2id",
			Parent: "v2parent",
			OS:     "os",
		},
		OSVersion: "osversion",
		RootFS: &image.RootFS{
			Type: "layers",
		},
	}
	v2js, err := json.Marshal(img)
	if err != nil {
		t.Fatal(err)
	}

	// Convert the image back in order to get RawJSON() support.
	img, err = image.NewFromJSON(v2js)
	if err != nil {
		t.Fatal(err)
	}

	js, err := MakeV1ConfigFromConfig(img, "v1id", "v1parent", false)
	if err != nil {
		t.Fatal(err)
	}

	newimg := &image.Image{}
	err = json.Unmarshal(js, newimg)
	if err != nil {
		t.Fatal(err)
	}

	if newimg.V1Image.ID != "v1id" || newimg.Parent != "v1parent" {
		t.Error("ids should have changed", newimg.V1Image.ID, newimg.V1Image.Parent)
	}

	if newimg.RootFS != nil {
		t.Error("rootfs should have been removed")
	}

	if newimg.V1Image.OS != "os" {
		t.Error("os should have been preserved")
	}
}
