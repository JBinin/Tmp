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
// Package golden provides function and helpers to use golden file for
// testing purpose.
package golden

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("test.update", false, "update golden file")

// Get returns the golden file content. If the `test.update` is specified, it updates the
// file with the current output and returns it.
func Get(t *testing.T, actual []byte, filename string) []byte {
	golden := filepath.Join("testdata", filename)
	if *update {
		if err := ioutil.WriteFile(golden, actual, 0644); err != nil {
			t.Fatal(err)
		}
	}
	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	return expected
}
