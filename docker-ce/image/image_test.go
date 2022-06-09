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
package image

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
)

const sampleImageJSON = `{
	"architecture": "amd64",
	"os": "linux",
	"config": {},
	"rootfs": {
		"type": "layers",
		"diff_ids": []
	}
}`

func TestJSON(t *testing.T) {
	img, err := NewFromJSON([]byte(sampleImageJSON))
	if err != nil {
		t.Fatal(err)
	}
	rawJSON := img.RawJSON()
	if string(rawJSON) != sampleImageJSON {
		t.Fatalf("raw JSON of config didn't match: expected %+v, got %v", sampleImageJSON, rawJSON)
	}
}

func TestInvalidJSON(t *testing.T) {
	_, err := NewFromJSON([]byte("{}"))
	if err == nil {
		t.Fatal("expected JSON parse error")
	}
}

func TestMarshalKeyOrder(t *testing.T) {
	b, err := json.Marshal(&Image{
		V1Image: V1Image{
			Comment:      "a",
			Author:       "b",
			Architecture: "c",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedOrder := []string{"architecture", "author", "comment"}
	var indexes []int
	for _, k := range expectedOrder {
		indexes = append(indexes, strings.Index(string(b), k))
	}

	if !sort.IntsAreSorted(indexes) {
		t.Fatal("invalid key order in JSON: ", string(b))
	}
}
