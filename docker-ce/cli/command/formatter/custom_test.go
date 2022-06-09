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
package formatter

import (
	"reflect"
	"strings"
	"testing"
)

func compareMultipleValues(t *testing.T, value, expected string) {
	// comma-separated values means probably a map input, which won't
	// be guaranteed to have the same order as our expected value
	// We'll create maps and use reflect.DeepEquals to check instead:
	entriesMap := make(map[string]string)
	expMap := make(map[string]string)
	entries := strings.Split(value, ",")
	expectedEntries := strings.Split(expected, ",")
	for _, entry := range entries {
		keyval := strings.Split(entry, "=")
		entriesMap[keyval[0]] = keyval[1]
	}
	for _, expected := range expectedEntries {
		keyval := strings.Split(expected, "=")
		expMap[keyval[0]] = keyval[1]
	}
	if !reflect.DeepEqual(expMap, entriesMap) {
		t.Fatalf("Expected entries: %v, got: %v", expected, value)
	}
}
