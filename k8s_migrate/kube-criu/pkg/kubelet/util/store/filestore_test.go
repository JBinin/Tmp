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
/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package store

import (
	"io/ioutil"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/kubernetes/pkg/util/filesystem"
)

func TestFileStore(t *testing.T) {
	path, err := ioutil.TempDir("", "FileStore")
	assert.NoError(t, err)
	store, err := NewFileStore(path, filesystem.DefaultFs{})
	assert.NoError(t, err)
	testStore(t, store)
}

func TestFakeFileStore(t *testing.T) {
	store, err := NewFileStore("/tmp/test-fake-file-store", filesystem.NewFakeFs())
	assert.NoError(t, err)
	testStore(t, store)
}

func testStore(t *testing.T, store Store) {
	testCases := []struct {
		key       string
		data      string
		expectErr bool
	}{
		{
			"id1",
			"data1",
			false,
		},
		{
			"id2",
			"data2",
			false,
		},
		{
			"/id1",
			"data1",
			true,
		},
		{
			".id1",
			"data1",
			true,
		},
		{
			"   ",
			"data2",
			true,
		},
		{
			"___",
			"data2",
			true,
		},
	}

	// Test add data.
	for _, c := range testCases {
		t.Log("test case: ", c)
		_, err := store.Read(c.key)
		assert.Error(t, err)

		err = store.Write(c.key, []byte(c.data))
		if c.expectErr {
			assert.Error(t, err)
			continue
		}

		require.NoError(t, err)
		// Test read data by key.
		data, err := store.Read(c.key)
		require.NoError(t, err)
		assert.Equal(t, string(data), c.data)
	}

	// Test list keys.
	keys, err := store.List()
	assert.NoError(t, err)
	sort.Strings(keys)
	assert.Equal(t, keys, []string{"id1", "id2"})

	// Test Delete data
	for _, c := range testCases {
		if c.expectErr {
			continue
		}

		err = store.Delete(c.key)
		require.NoError(t, err)
		_, err = store.Read(c.key)
		assert.EqualValues(t, ErrKeyNotFound, err)
	}

	// Test delete non-existent key.
	err = store.Delete("id1")
	assert.NoError(t, err)

	// Test list keys.
	keys, err = store.List()
	require.NoError(t, err)
	assert.Equal(t, len(keys), 0)
}
