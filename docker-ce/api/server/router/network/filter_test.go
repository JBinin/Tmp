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

package network

import (
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func TestFilterNetworks(t *testing.T) {
	networks := []types.NetworkResource{
		{
			Name:   "host",
			Driver: "host",
		},
		{
			Name:   "bridge",
			Driver: "bridge",
		},
		{
			Name:   "none",
			Driver: "null",
		},
		{
			Name:   "myoverlay",
			Driver: "overlay",
		},
		{
			Name:   "mydrivernet",
			Driver: "mydriver",
		},
	}

	bridgeDriverFilters := filters.NewArgs()
	bridgeDriverFilters.Add("driver", "bridge")

	overlayDriverFilters := filters.NewArgs()
	overlayDriverFilters.Add("driver", "overlay")

	nonameDriverFilters := filters.NewArgs()
	nonameDriverFilters.Add("driver", "noname")

	customDriverFilters := filters.NewArgs()
	customDriverFilters.Add("type", "custom")

	builtinDriverFilters := filters.NewArgs()
	builtinDriverFilters.Add("type", "builtin")

	invalidDriverFilters := filters.NewArgs()
	invalidDriverFilters.Add("type", "invalid")

	testCases := []struct {
		filter      filters.Args
		resultCount int
		err         string
	}{
		{
			filter:      bridgeDriverFilters,
			resultCount: 1,
			err:         "",
		},
		{
			filter:      overlayDriverFilters,
			resultCount: 1,
			err:         "",
		},
		{
			filter:      nonameDriverFilters,
			resultCount: 0,
			err:         "",
		},
		{
			filter:      customDriverFilters,
			resultCount: 2,
			err:         "",
		},
		{
			filter:      builtinDriverFilters,
			resultCount: 3,
			err:         "",
		},
		{
			filter:      invalidDriverFilters,
			resultCount: 0,
			err:         "Invalid filter: 'type'='invalid'",
		},
	}

	for _, testCase := range testCases {
		result, err := filterNetworks(networks, testCase.filter)
		if testCase.err != "" {
			if err == nil {
				t.Fatalf("expect error '%s', got no error", testCase.err)

			} else if !strings.Contains(err.Error(), testCase.err) {
				t.Fatalf("expect error '%s', got '%s'", testCase.err, err)
			}
		} else {
			if err != nil {
				t.Fatalf("expect no error, got error '%s'", err)
			}
			// Make sure result is not nil
			if result == nil {
				t.Fatal("filterNetworks should not return nil")
			}

			if len(result) != testCase.resultCount {
				t.Fatalf("expect '%d' networks, got '%d' networks", testCase.resultCount, len(result))
			}
		}
	}
}
