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
package longpath

import (
	"strings"
	"testing"
)

func TestStandardLongPath(t *testing.T) {
	c := `C:\simple\path`
	longC := AddPrefix(c)
	if !strings.EqualFold(longC, `\\?\C:\simple\path`) {
		t.Errorf("Wrong long path returned. Original = %s ; Long = %s", c, longC)
	}
}

func TestUNCLongPath(t *testing.T) {
	c := `\\server\share\path`
	longC := AddPrefix(c)
	if !strings.EqualFold(longC, `\\?\UNC\server\share\path`) {
		t.Errorf("Wrong UNC long path returned. Original = %s ; Long = %s", c, longC)
	}
}
