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
package jsonlog

import (
	"testing"
	"time"
)

// Testing to ensure 'year' fields is between 0 and 9999
func TestFastTimeMarshalJSONWithInvalidDate(t *testing.T) {
	aTime := time.Date(-1, 1, 1, 0, 0, 0, 0, time.Local)
	json, err := FastTimeMarshalJSON(aTime)
	if err == nil {
		t.Fatalf("FastTimeMarshalJSON should throw an error, but was '%v'", json)
	}
	anotherTime := time.Date(10000, 1, 1, 0, 0, 0, 0, time.Local)
	json, err = FastTimeMarshalJSON(anotherTime)
	if err == nil {
		t.Fatalf("FastTimeMarshalJSON should throw an error, but was '%v'", json)
	}

}

func TestFastTimeMarshalJSON(t *testing.T) {
	aTime := time.Date(2015, 5, 29, 11, 1, 2, 3, time.UTC)
	json, err := FastTimeMarshalJSON(aTime)
	if err != nil {
		t.Fatal(err)
	}
	expected := "\"2015-05-29T11:01:02.000000003Z\""
	if json != expected {
		t.Fatalf("Expected %v, got %v", expected, json)
	}

	location, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatal(err)
	}
	aTime = time.Date(2015, 5, 29, 11, 1, 2, 3, location)
	json, err = FastTimeMarshalJSON(aTime)
	if err != nil {
		t.Fatal(err)
	}
	expected = "\"2015-05-29T11:01:02.000000003+02:00\""
	if json != expected {
		t.Fatalf("Expected %v, got %v", expected, json)
	}
}
