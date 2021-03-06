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
package container

import (
	"testing"

	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestBuildContainerListOptions(t *testing.T) {
	filters := opts.NewFilterOpt()
	assert.NilError(t, filters.Set("foo=bar"))
	assert.NilError(t, filters.Set("baz=foo"))

	contexts := []struct {
		psOpts          *psOptions
		expectedAll     bool
		expectedSize    bool
		expectedLimit   int
		expectedFilters map[string]string
	}{
		{
			psOpts: &psOptions{
				all:    true,
				size:   true,
				last:   5,
				filter: filters,
			},
			expectedAll:   true,
			expectedSize:  true,
			expectedLimit: 5,
			expectedFilters: map[string]string{
				"foo": "bar",
				"baz": "foo",
			},
		},
		{
			psOpts: &psOptions{
				all:     true,
				size:    true,
				last:    -1,
				nLatest: true,
			},
			expectedAll:     true,
			expectedSize:    true,
			expectedLimit:   1,
			expectedFilters: make(map[string]string),
		},
		{
			psOpts: &psOptions{
				all:    true,
				size:   false,
				last:   5,
				filter: filters,
				// With .Size, size should be true
				format: "{{.Size}}",
			},
			expectedAll:   true,
			expectedSize:  true,
			expectedLimit: 5,
			expectedFilters: map[string]string{
				"foo": "bar",
				"baz": "foo",
			},
		},
		{
			psOpts: &psOptions{
				all:    true,
				size:   false,
				last:   5,
				filter: filters,
				// With .Size, size should be true
				format: "{{.Size}} {{.CreatedAt}} {{.Networks}}",
			},
			expectedAll:   true,
			expectedSize:  true,
			expectedLimit: 5,
			expectedFilters: map[string]string{
				"foo": "bar",
				"baz": "foo",
			},
		},
		{
			psOpts: &psOptions{
				all:    true,
				size:   false,
				last:   5,
				filter: filters,
				// Without .Size, size should be false
				format: "{{.CreatedAt}} {{.Networks}}",
			},
			expectedAll:   true,
			expectedSize:  false,
			expectedLimit: 5,
			expectedFilters: map[string]string{
				"foo": "bar",
				"baz": "foo",
			},
		},
	}

	for _, c := range contexts {
		options, err := buildContainerListOptions(c.psOpts)
		assert.NilError(t, err)

		assert.Equal(t, c.expectedAll, options.All)
		assert.Equal(t, c.expectedSize, options.Size)
		assert.Equal(t, c.expectedLimit, options.Limit)
		assert.Equal(t, options.Filters.Len(), len(c.expectedFilters))

		for k, v := range c.expectedFilters {
			f := options.Filters
			if !f.ExactMatch(k, v) {
				t.Fatalf("Expected filter with key %s to be %s but got %s", k, v, f.Get(k))
			}
		}
	}
}
