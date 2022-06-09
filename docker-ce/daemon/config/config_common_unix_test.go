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

package config

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestCommonUnixValidateConfigurationErrors(t *testing.T) {
	testCases := []struct {
		config *Config
	}{
		// Can't override the stock runtime
		{
			config: &Config{
				CommonUnixConfig: CommonUnixConfig{
					Runtimes: map[string]types.Runtime{
						StockRuntimeName: {},
					},
				},
			},
		},
		// Default runtime should be present in runtimes
		{
			config: &Config{
				CommonUnixConfig: CommonUnixConfig{
					Runtimes: map[string]types.Runtime{
						"foo": {},
					},
					DefaultRuntime: "bar",
				},
			},
		},
	}
	for _, tc := range testCases {
		err := Validate(tc.config)
		if err == nil {
			t.Fatalf("expected error, got nil for config %v", tc.config)
		}
	}
}
