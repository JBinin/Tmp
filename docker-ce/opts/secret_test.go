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
package opts

import (
	"os"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

func TestSecretOptionsSimple(t *testing.T) {
	var opt SecretOpt

	testCase := "app-secret"
	assert.NilError(t, opt.Set(testCase))

	reqs := opt.Value()
	assert.Equal(t, len(reqs), 1)
	req := reqs[0]
	assert.Equal(t, req.Source, "app-secret")
	assert.Equal(t, req.Target, "app-secret")
	assert.Equal(t, req.UID, "0")
	assert.Equal(t, req.GID, "0")
}

func TestSecretOptionsSourceTarget(t *testing.T) {
	var opt SecretOpt

	testCase := "source=foo,target=testing"
	assert.NilError(t, opt.Set(testCase))

	reqs := opt.Value()
	assert.Equal(t, len(reqs), 1)
	req := reqs[0]
	assert.Equal(t, req.Source, "foo")
	assert.Equal(t, req.Target, "testing")
}

func TestSecretOptionsShorthand(t *testing.T) {
	var opt SecretOpt

	testCase := "src=foo,target=testing"
	assert.NilError(t, opt.Set(testCase))

	reqs := opt.Value()
	assert.Equal(t, len(reqs), 1)
	req := reqs[0]
	assert.Equal(t, req.Source, "foo")
}

func TestSecretOptionsCustomUidGid(t *testing.T) {
	var opt SecretOpt

	testCase := "source=foo,target=testing,uid=1000,gid=1001"
	assert.NilError(t, opt.Set(testCase))

	reqs := opt.Value()
	assert.Equal(t, len(reqs), 1)
	req := reqs[0]
	assert.Equal(t, req.Source, "foo")
	assert.Equal(t, req.Target, "testing")
	assert.Equal(t, req.UID, "1000")
	assert.Equal(t, req.GID, "1001")
}

func TestSecretOptionsCustomMode(t *testing.T) {
	var opt SecretOpt

	testCase := "source=foo,target=testing,uid=1000,gid=1001,mode=0444"
	assert.NilError(t, opt.Set(testCase))

	reqs := opt.Value()
	assert.Equal(t, len(reqs), 1)
	req := reqs[0]
	assert.Equal(t, req.Source, "foo")
	assert.Equal(t, req.Target, "testing")
	assert.Equal(t, req.UID, "1000")
	assert.Equal(t, req.GID, "1001")
	assert.Equal(t, req.Mode, os.FileMode(0444))
}
