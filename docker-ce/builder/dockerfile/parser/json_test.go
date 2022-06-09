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
package parser

import (
	"testing"
)

var invalidJSONArraysOfStrings = []string{
	`["a",42,"b"]`,
	`["a",123.456,"b"]`,
	`["a",{},"b"]`,
	`["a",{"c": "d"},"b"]`,
	`["a",["c"],"b"]`,
	`["a",true,"b"]`,
	`["a",false,"b"]`,
	`["a",null,"b"]`,
}

var validJSONArraysOfStrings = map[string][]string{
	`[]`:           {},
	`[""]`:         {""},
	`["a"]`:        {"a"},
	`["a","b"]`:    {"a", "b"},
	`[ "a", "b" ]`: {"a", "b"},
	`[	"a",	"b"	]`: {"a", "b"},
	`	[	"a",	"b"	]	`: {"a", "b"},
	`["abc 123", "♥", "☃", "\" \\ \/ \b \f \n \r \t \u0000"]`: {"abc 123", "♥", "☃", "\" \\ / \b \f \n \r \t \u0000"},
}

func TestJSONArraysOfStrings(t *testing.T) {
	for json, expected := range validJSONArraysOfStrings {
		d := Directive{}
		SetEscapeToken(DefaultEscapeToken, &d)

		if node, _, err := parseJSON(json, &d); err != nil {
			t.Fatalf("%q should be a valid JSON array of strings, but wasn't! (err: %q)", json, err)
		} else {
			i := 0
			for node != nil {
				if i >= len(expected) {
					t.Fatalf("expected result is shorter than parsed result (%d vs %d+) in %q", len(expected), i+1, json)
				}
				if node.Value != expected[i] {
					t.Fatalf("expected %q (not %q) in %q at pos %d", expected[i], node.Value, json, i)
				}
				node = node.Next
				i++
			}
			if i != len(expected) {
				t.Fatalf("expected result is longer than parsed result (%d vs %d) in %q", len(expected), i+1, json)
			}
		}
	}
	for _, json := range invalidJSONArraysOfStrings {
		d := Directive{}
		SetEscapeToken(DefaultEscapeToken, &d)

		if _, _, err := parseJSON(json, &d); err != errDockerfileNotStringArray {
			t.Fatalf("%q should be an invalid JSON array of strings, but wasn't!", json)
		}
	}
}
