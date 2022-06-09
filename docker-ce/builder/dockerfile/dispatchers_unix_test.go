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
// +build !windows

package dockerfile

import (
	"testing"
)

func TestNormaliseWorkdir(t *testing.T) {
	testCases := []struct{ current, requested, expected, expectedError string }{
		{``, ``, ``, `cannot normalise nothing`},
		{``, `foo`, `/foo`, ``},
		{``, `/foo`, `/foo`, ``},
		{`/foo`, `bar`, `/foo/bar`, ``},
		{`/foo`, `/bar`, `/bar`, ``},
	}

	for _, test := range testCases {
		normalised, err := normaliseWorkdir(test.current, test.requested)

		if test.expectedError != "" && err == nil {
			t.Fatalf("NormaliseWorkdir should return an error %s, got nil", test.expectedError)
		}

		if test.expectedError != "" && err.Error() != test.expectedError {
			t.Fatalf("NormaliseWorkdir returned wrong error. Expected %s, got %s", test.expectedError, err.Error())
		}

		if normalised != test.expected {
			t.Fatalf("NormaliseWorkdir error. Expected %s for current %s and requested %s, got %s", test.expected, test.current, test.requested, normalised)
		}
	}
}
