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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/internal/test"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestSwarmJoinErrors(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		swarmJoinFunc func() error
		infoFunc      func() (types.Info, error)
		expectedError string
	}{
		{
			name:          "not-enough-args",
			expectedError: "requires exactly 1 argument",
		},
		{
			name:          "too-many-args",
			args:          []string{"remote1", "remote2"},
			expectedError: "requires exactly 1 argument",
		},
		{
			name: "join-failed",
			args: []string{"remote"},
			swarmJoinFunc: func() error {
				return fmt.Errorf("error joining the swarm")
			},
			expectedError: "error joining the swarm",
		},
		{
			name: "join-failed-on-init",
			args: []string{"remote"},
			infoFunc: func() (types.Info, error) {
				return types.Info{}, fmt.Errorf("error asking for node info")
			},
			expectedError: "error asking for node info",
		},
	}
	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		cmd := newJoinCommand(
			test.NewFakeCli(&fakeClient{
				swarmJoinFunc: tc.swarmJoinFunc,
				infoFunc:      tc.infoFunc,
			}, buf))
		cmd.SetArgs(tc.args)
		cmd.SetOutput(ioutil.Discard)
		assert.Error(t, cmd.Execute(), tc.expectedError)
	}
}

func TestSwarmJoin(t *testing.T) {
	testCases := []struct {
		name     string
		infoFunc func() (types.Info, error)
		expected string
	}{
		{
			name: "join-as-manager",
			infoFunc: func() (types.Info, error) {
				return types.Info{
					Swarm: swarm.Info{
						ControlAvailable: true,
					},
				}, nil
			},
			expected: "This node joined a swarm as a manager.",
		},
		{
			name: "join-as-worker",
			infoFunc: func() (types.Info, error) {
				return types.Info{
					Swarm: swarm.Info{
						ControlAvailable: false,
					},
				}, nil
			},
			expected: "This node joined a swarm as a worker.",
		},
	}
	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		cmd := newJoinCommand(
			test.NewFakeCli(&fakeClient{
				infoFunc: tc.infoFunc,
			}, buf))
		cmd.SetArgs([]string{"remote"})
		assert.NilError(t, cmd.Execute())
		assert.Equal(t, strings.TrimSpace(buf.String()), tc.expected)
	}
}
