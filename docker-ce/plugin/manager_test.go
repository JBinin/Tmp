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
package plugin

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestValidatePrivileges(t *testing.T) {
	testData := map[string]struct {
		requiredPrivileges types.PluginPrivileges
		privileges         types.PluginPrivileges
		result             bool
	}{
		"diff-len": {
			requiredPrivileges: []types.PluginPrivilege{
				{"Privilege1", "Description", []string{"abc", "def", "ghi"}},
			},
			privileges: []types.PluginPrivilege{
				{"Privilege1", "Description", []string{"abc", "def", "ghi"}},
				{"Privilege2", "Description", []string{"123", "456", "789"}},
			},
			result: false,
		},
		"diff-value": {
			requiredPrivileges: []types.PluginPrivilege{
				{"Privilege1", "Description", []string{"abc", "def", "GHI"}},
				{"Privilege2", "Description", []string{"123", "456", "***"}},
			},
			privileges: []types.PluginPrivilege{
				{"Privilege1", "Description", []string{"abc", "def", "ghi"}},
				{"Privilege2", "Description", []string{"123", "456", "789"}},
			},
			result: false,
		},
		"diff-order-but-same-value": {
			requiredPrivileges: []types.PluginPrivilege{
				{"Privilege1", "Description", []string{"abc", "def", "GHI"}},
				{"Privilege2", "Description", []string{"123", "456", "789"}},
			},
			privileges: []types.PluginPrivilege{
				{"Privilege2", "Description", []string{"123", "456", "789"}},
				{"Privilege1", "Description", []string{"GHI", "abc", "def"}},
			},
			result: true,
		},
	}

	for key, data := range testData {
		err := validatePrivileges(data.requiredPrivileges, data.privileges)
		if (err == nil) != data.result {
			t.Fatalf("Test item %s expected result to be %t, got %t", key, data.result, (err == nil))
		}
	}
}
