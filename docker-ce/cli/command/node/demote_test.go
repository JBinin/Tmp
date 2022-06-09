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

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/internal/test"
	// Import builders to get the builder function as package function
	. "github.com/docker/docker/cli/internal/test/builders"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestNodeDemoteErrors(t *testing.T) {
	testCases := []struct {
		args            []string
		nodeInspectFunc func() (swarm.Node, []byte, error)
		nodeUpdateFunc  func(nodeID string, version swarm.Version, node swarm.NodeSpec) error
		expectedError   string
	}{
		{
			expectedError: "requires at least 1 argument",
		},
		{
			args: []string{"nodeID"},
			nodeInspectFunc: func() (swarm.Node, []byte, error) {
				return swarm.Node{}, []byte{}, fmt.Errorf("error inspecting the node")
			},
			expectedError: "error inspecting the node",
		},
		{
			args: []string{"nodeID"},
			nodeUpdateFunc: func(nodeID string, version swarm.Version, node swarm.NodeSpec) error {
				return fmt.Errorf("error updating the node")
			},
			expectedError: "error updating the node",
		},
	}
	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		cmd := newDemoteCommand(
			test.NewFakeCli(&fakeClient{
				nodeInspectFunc: tc.nodeInspectFunc,
				nodeUpdateFunc:  tc.nodeUpdateFunc,
			}, buf))
		cmd.SetArgs(tc.args)
		cmd.SetOutput(ioutil.Discard)
		assert.Error(t, cmd.Execute(), tc.expectedError)
	}
}

func TestNodeDemoteNoChange(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newDemoteCommand(
		test.NewFakeCli(&fakeClient{
			nodeInspectFunc: func() (swarm.Node, []byte, error) {
				return *Node(), []byte{}, nil
			},
			nodeUpdateFunc: func(nodeID string, version swarm.Version, node swarm.NodeSpec) error {
				if node.Role != swarm.NodeRoleWorker {
					return fmt.Errorf("expected role worker, got %s", node.Role)
				}
				return nil
			},
		}, buf))
	cmd.SetArgs([]string{"nodeID"})
	assert.NilError(t, cmd.Execute())
}

func TestNodeDemoteMultipleNode(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newDemoteCommand(
		test.NewFakeCli(&fakeClient{
			nodeInspectFunc: func() (swarm.Node, []byte, error) {
				return *Node(Manager()), []byte{}, nil
			},
			nodeUpdateFunc: func(nodeID string, version swarm.Version, node swarm.NodeSpec) error {
				if node.Role != swarm.NodeRoleWorker {
					return fmt.Errorf("expected role worker, got %s", node.Role)
				}
				return nil
			},
		}, buf))
	cmd.SetArgs([]string{"nodeID1", "nodeID2"})
	assert.NilError(t, cmd.Execute())
}
