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
package cluster

import (
	apierrors "github.com/docker/docker/api/errors"
	apitypes "github.com/docker/docker/api/types"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	swarmapi "github.com/docker/swarmkit/api"
)

// GetNodes returns a list of all nodes known to a cluster.
func (c *Cluster) GetNodes(options apitypes.NodeListOptions) ([]types.Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return nil, c.errNoManager(state)
	}

	filters, err := newListNodesFilters(options.Filters)
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListNodes(
		ctx,
		&swarmapi.ListNodesRequest{Filters: filters})
	if err != nil {
		return nil, err
	}

	nodes := []types.Node{}

	for _, node := range r.Nodes {
		nodes = append(nodes, convert.NodeFromGRPC(*node))
	}
	return nodes, nil
}

// GetNode returns a node based on an ID.
func (c *Cluster) GetNode(input string) (types.Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return types.Node{}, c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	node, err := getNode(ctx, state.controlClient, input)
	if err != nil {
		return types.Node{}, err
	}
	return convert.NodeFromGRPC(*node), nil
}

// UpdateNode updates existing nodes properties.
func (c *Cluster) UpdateNode(input string, version uint64, spec types.NodeSpec) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return c.errNoManager(state)
	}

	nodeSpec, err := convert.NodeSpecToGRPC(spec)
	if err != nil {
		return apierrors.NewBadRequestError(err)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	currentNode, err := getNode(ctx, state.controlClient, input)
	if err != nil {
		return err
	}

	_, err = state.controlClient.UpdateNode(
		ctx,
		&swarmapi.UpdateNodeRequest{
			NodeID: currentNode.ID,
			Spec:   &nodeSpec,
			NodeVersion: &swarmapi.Version{
				Index: version,
			},
		},
	)
	return err
}

// RemoveNode removes a node from a cluster
func (c *Cluster) RemoveNode(input string, force bool) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	node, err := getNode(ctx, state.controlClient, input)
	if err != nil {
		return err
	}

	_, err = state.controlClient.RemoveNode(ctx, &swarmapi.RemoveNodeRequest{NodeID: node.ID, Force: force})
	return err
}
