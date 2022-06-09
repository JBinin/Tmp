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
package libcontainerd

import (
	"testing"
)

func TestEnvironmentParsing(t *testing.T) {
	env := []string{"foo=bar", "car=hat", "a=b=c"}
	result := setupEnvironmentVariables(env)
	if len(result) != 3 || result["foo"] != "bar" || result["car"] != "hat" || result["a"] != "b=c" {
		t.Fatalf("Expected map[foo:bar car:hat a:b=c], got %v", result)
	}
}
