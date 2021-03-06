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
package json

import (
	"reflect"
	"testing"
)

func TestLookupPtrToStruct(t *testing.T) {
	type Elem struct {
		Key   string
		Value string
	}
	type Outer struct {
		Inner []Elem `json:"inner" patchStrategy:"merge" patchMergeKey:"key"`
	}
	outer := &Outer{}
	elemType, patchStrategies, patchMergeKey, err := LookupPatchMetadataForStruct(reflect.TypeOf(outer), "inner")
	if err != nil {
		t.Fatal(err)
	}
	if elemType != reflect.TypeOf([]Elem{}) {
		t.Errorf("elemType = %v, want: %v", elemType, reflect.TypeOf([]Elem{}))
	}
	if !reflect.DeepEqual(patchStrategies, []string{"merge"}) {
		t.Errorf("patchStrategies = %v, want: %v", patchStrategies, []string{"merge"})
	}
	if patchMergeKey != "key" {
		t.Errorf("patchMergeKey = %v, want: %v", patchMergeKey, "key")
	}
}
