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
	"fmt"

	"github.com/Sirupsen/logrus"
	apierrors "github.com/docker/docker/api/errors"
	apitypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	"github.com/docker/docker/runconfig"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// GetNetworks returns all current cluster managed networks.
func (c *Cluster) GetNetworks() ([]apitypes.NetworkResource, error) {
	return c.getNetworks(nil)
}

func (c *Cluster) getNetworks(filters *swarmapi.ListNetworksRequest_Filters) ([]apitypes.NetworkResource, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return nil, c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListNetworks(ctx, &swarmapi.ListNetworksRequest{Filters: filters})
	if err != nil {
		return nil, err
	}

	var networks []apitypes.NetworkResource

	for _, network := range r.Networks {
		networks = append(networks, convert.BasicNetworkFromGRPC(*network))
	}

	return networks, nil
}

// GetNetwork returns a cluster network by an ID.
func (c *Cluster) GetNetwork(input string) (apitypes.NetworkResource, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return apitypes.NetworkResource{}, c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	network, err := getNetwork(ctx, state.controlClient, input)
	if err != nil {
		return apitypes.NetworkResource{}, err
	}
	return convert.BasicNetworkFromGRPC(*network), nil
}

// GetNetworksByName returns cluster managed networks by name.
// It is ok to have multiple networks here. #18864
func (c *Cluster) GetNetworksByName(name string) ([]apitypes.NetworkResource, error) {
	// Note that swarmapi.GetNetworkRequest.Name is not functional.
	// So we cannot just use that with c.GetNetwork.
	return c.getNetworks(&swarmapi.ListNetworksRequest_Filters{
		Names: []string{name},
	})
}

func attacherKey(target, containerID string) string {
	return containerID + ":" + target
}

// UpdateAttachment signals the attachment config to the attachment
// waiter who is trying to start or attach the container to the
// network.
func (c *Cluster) UpdateAttachment(target, containerID string, config *network.NetworkingConfig) error {
	c.mu.RLock()
	attacher, ok := c.attachers[attacherKey(target, containerID)]
	c.mu.RUnlock()
	if !ok || attacher == nil {
		return fmt.Errorf("could not find attacher for container %s to network %s", containerID, target)
	}

	attacher.attachWaitCh <- config
	close(attacher.attachWaitCh)
	return nil
}

// WaitForDetachment waits for the container to stop or detach from
// the network.
func (c *Cluster) WaitForDetachment(ctx context.Context, networkName, networkID, taskID, containerID string) error {
	c.mu.RLock()
	attacher, ok := c.attachers[attacherKey(networkName, containerID)]
	if !ok {
		attacher, ok = c.attachers[attacherKey(networkID, containerID)]
	}
	state := c.currentNodeState()
	if state.swarmNode == nil || state.swarmNode.Agent() == nil {
		c.mu.RUnlock()
		return errors.New("invalid cluster node while waiting for detachment")
	}

	c.mu.RUnlock()
	agent := state.swarmNode.Agent()
	if ok && attacher != nil &&
		attacher.detachWaitCh != nil &&
		attacher.attachCompleteCh != nil {
		// Attachment may be in progress still so wait for
		// attachment to complete.
		select {
		case <-attacher.attachCompleteCh:
		case <-ctx.Done():
			return ctx.Err()
		}

		if attacher.taskID == taskID {
			select {
			case <-attacher.detachWaitCh:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return agent.ResourceAllocator().DetachNetwork(ctx, taskID)
}

// AttachNetwork generates an attachment request towards the manager.
func (c *Cluster) AttachNetwork(target string, containerID string, addresses []string) (*network.NetworkingConfig, error) {
	aKey := attacherKey(target, containerID)
	c.mu.Lock()
	state := c.currentNodeState()
	if state.swarmNode == nil || state.swarmNode.Agent() == nil {
		c.mu.Unlock()
		return nil, errors.New("invalid cluster node while attaching to network")
	}
	if attacher, ok := c.attachers[aKey]; ok {
		c.mu.Unlock()
		return attacher.config, nil
	}

	agent := state.swarmNode.Agent()
	attachWaitCh := make(chan *network.NetworkingConfig)
	detachWaitCh := make(chan struct{})
	attachCompleteCh := make(chan struct{})
	c.attachers[aKey] = &attacher{
		attachWaitCh:     attachWaitCh,
		attachCompleteCh: attachCompleteCh,
		detachWaitCh:     detachWaitCh,
	}
	c.mu.Unlock()

	ctx, cancel := c.getRequestContext()
	defer cancel()

	taskID, err := agent.ResourceAllocator().AttachNetwork(ctx, containerID, target, addresses)
	if err != nil {
		c.mu.Lock()
		delete(c.attachers, aKey)
		c.mu.Unlock()
		return nil, fmt.Errorf("Could not attach to network %s: %v", target, err)
	}

	c.mu.Lock()
	c.attachers[aKey].taskID = taskID
	close(attachCompleteCh)
	c.mu.Unlock()

	logrus.Debugf("Successfully attached to network %s with task id %s", target, taskID)

	release := func() {
		ctx, cancel := c.getRequestContext()
		defer cancel()
		if err := agent.ResourceAllocator().DetachNetwork(ctx, taskID); err != nil {
			logrus.Errorf("Failed remove network attachment %s to network %s on allocation failure: %v",
				taskID, target, err)
		}
	}

	var config *network.NetworkingConfig
	select {
	case config = <-attachWaitCh:
	case <-ctx.Done():
		release()
		return nil, fmt.Errorf("attaching to network failed, make sure your network options are correct and check manager logs: %v", ctx.Err())
	}

	c.mu.Lock()
	c.attachers[aKey].config = config
	c.mu.Unlock()

	logrus.Debugf("Successfully allocated resources on network %s for task id %s", target, taskID)

	return config, nil
}

// DetachNetwork unblocks the waiters waiting on WaitForDetachment so
// that a request to detach can be generated towards the manager.
func (c *Cluster) DetachNetwork(target string, containerID string) error {
	aKey := attacherKey(target, containerID)

	c.mu.Lock()
	attacher, ok := c.attachers[aKey]
	delete(c.attachers, aKey)
	c.mu.Unlock()

	if !ok {
		return fmt.Errorf("could not find network attachment for container %s to network %s", containerID, target)
	}

	close(attacher.detachWaitCh)
	return nil
}

// CreateNetwork creates a new cluster managed network.
func (c *Cluster) CreateNetwork(s apitypes.NetworkCreateRequest) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return "", c.errNoManager(state)
	}

	if runconfig.IsPreDefinedNetwork(s.Name) {
		err := fmt.Errorf("%s is a pre-defined network and cannot be created", s.Name)
		return "", apierrors.NewRequestForbiddenError(err)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	networkSpec := convert.BasicNetworkCreateToGRPC(s)
	r, err := state.controlClient.CreateNetwork(ctx, &swarmapi.CreateNetworkRequest{Spec: &networkSpec})
	if err != nil {
		return "", err
	}

	return r.Network.ID, nil
}

// RemoveNetwork removes a cluster network.
func (c *Cluster) RemoveNetwork(input string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	network, err := getNetwork(ctx, state.controlClient, input)
	if err != nil {
		return err
	}

	_, err = state.controlClient.RemoveNetwork(ctx, &swarmapi.RemoveNetworkRequest{NetworkID: network.ID})
	return err
}

func (c *Cluster) populateNetworkID(ctx context.Context, client swarmapi.ControlClient, s *types.ServiceSpec) error {
	// Always prefer NetworkAttachmentConfigs from TaskTemplate
	// but fallback to service spec for backward compatibility
	networks := s.TaskTemplate.Networks
	if len(networks) == 0 {
		networks = s.Networks
	}

	for i, n := range networks {
		apiNetwork, err := getNetwork(ctx, client, n.Target)
		if err != nil {
			if ln, _ := c.config.Backend.FindNetwork(n.Target); ln != nil && !ln.Info().Dynamic() {
				err = fmt.Errorf("The network %s cannot be used with services. Only networks scoped to the swarm can be used, such as those created with the overlay driver.", ln.Name())
				return apierrors.NewRequestForbiddenError(err)
			}
			return err
		}
		networks[i].Target = apiNetwork.ID
	}
	return nil
}
