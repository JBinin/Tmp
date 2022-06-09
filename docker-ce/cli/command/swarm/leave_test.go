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
package swarm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/cli/internal/test"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestSwarmLeaveErrors(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		swarmLeaveFunc func() error
		expectedError  string
	}{
		{
			name:          "too-many-args",
			args:          []string{"foo"},
			expectedError: "accepts no argument(s)",
		},
		{
			name: "leave-failed",
			swarmLeaveFunc: func() error {
				return fmt.Errorf("error leaving the swarm")
			},
			expectedError: "error leaving the swarm",
		},
	}
	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		cmd := newLeaveCommand(
			test.NewFakeCli(&fakeClient{
				swarmLeaveFunc: tc.swarmLeaveFunc,
			}, buf))
		cmd.SetArgs(tc.args)
		cmd.SetOutput(ioutil.Discard)
		assert.Error(t, cmd.Execute(), tc.expectedError)
	}
}

func TestSwarmLeave(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newLeaveCommand(
		test.NewFakeCli(&fakeClient{}, buf))
	assert.NilError(t, cmd.Execute())
	assert.Equal(t, strings.TrimSpace(buf.String()), "Node left the swarm.")
}
