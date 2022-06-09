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
// +build windows

package dockerfile

import "testing"

func TestNormaliseWorkdir(t *testing.T) {
	tests := []struct{ current, requested, expected, etext string }{
		{``, ``, ``, `cannot normalise nothing`},
		{``, `C:`, ``, `C:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`},
		{``, `C:.`, ``, `C:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`},
		{`c:`, `\a`, ``, `c:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`},
		{`c:.`, `\a`, ``, `c:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`},
		{``, `a`, `C:\a`, ``},
		{``, `c:\foo`, `C:\foo`, ``},
		{``, `c:\\foo`, `C:\foo`, ``},
		{``, `\foo`, `C:\foo`, ``},
		{``, `\\foo`, `C:\foo`, ``},
		{``, `/foo`, `C:\foo`, ``},
		{``, `C:/foo`, `C:\foo`, ``},
		{`C:\foo`, `bar`, `C:\foo\bar`, ``},
		{`C:\foo`, `/bar`, `C:\bar`, ``},
		{`C:\foo`, `\bar`, `C:\bar`, ``},
	}
	for _, i := range tests {
		r, e := normaliseWorkdir(i.current, i.requested)

		if i.etext != "" && e == nil {
			t.Fatalf("TestNormaliseWorkingDir Expected error %s for '%s' '%s', got no error", i.etext, i.current, i.requested)
		}

		if i.etext != "" && e.Error() != i.etext {
			t.Fatalf("TestNormaliseWorkingDir Expected error %s for '%s' '%s', got %s", i.etext, i.current, i.requested, e.Error())
		}

		if r != i.expected {
			t.Fatalf("TestNormaliseWorkingDir Expected '%s' for '%s' '%s', got '%s'", i.expected, i.current, i.requested, r)
		}
	}
}
