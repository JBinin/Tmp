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
package node

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/cli/internal/test"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestNodeRemoveErrors(t *testing.T) {
	testCases := []struct {
		args           []string
		nodeRemoveFunc func() error
		expectedError  string
	}{
		{
			expectedError: "requires at least 1 argument",
		},
		{
			args: []string{"nodeID"},
			nodeRemoveFunc: func() error {
				return fmt.Errorf("error removing the node")
			},
			expectedError: "error removing the node",
		},
	}
	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		cmd := newRemoveCommand(
			test.NewFakeCli(&fakeClient{
				nodeRemoveFunc: tc.nodeRemoveFunc,
			}, buf))
		cmd.SetArgs(tc.args)
		cmd.SetOutput(ioutil.Discard)
		assert.Error(t, cmd.Execute(), tc.expectedError)
	}
}

func TestNodeRemoveMultiple(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newRemoveCommand(test.NewFakeCli(&fakeClient{}, buf))
	cmd.SetArgs([]string{"nodeID1", "nodeID2"})
	assert.NilError(t, cmd.Execute())
}
