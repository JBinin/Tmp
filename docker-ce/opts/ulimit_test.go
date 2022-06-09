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
package opts

import (
	"testing"

	"github.com/docker/go-units"
)

func TestUlimitOpt(t *testing.T) {
	ulimitMap := map[string]*units.Ulimit{
		"nofile": {"nofile", 1024, 512},
	}

	ulimitOpt := NewUlimitOpt(&ulimitMap)

	expected := "[nofile=512:1024]"
	if ulimitOpt.String() != expected {
		t.Fatalf("Expected %v, got %v", expected, ulimitOpt)
	}

	// Valid ulimit append to opts
	if err := ulimitOpt.Set("core=1024:1024"); err != nil {
		t.Fatal(err)
	}

	// Invalid ulimit type returns an error and do not append to opts
	if err := ulimitOpt.Set("notavalidtype=1024:1024"); err == nil {
		t.Fatalf("Expected error on invalid ulimit type")
	}
	expected = "[nofile=512:1024 core=1024:1024]"
	expected2 := "[core=1024:1024 nofile=512:1024]"
	result := ulimitOpt.String()
	if result != expected && result != expected2 {
		t.Fatalf("Expected %v or %v, got %v", expected, expected2, ulimitOpt)
	}

	// And test GetList
	ulimits := ulimitOpt.GetList()
	if len(ulimits) != 2 {
		t.Fatalf("Expected a ulimit list of 2, got %v", ulimits)
	}
}
