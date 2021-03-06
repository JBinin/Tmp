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
	apitypes "github.com/docker/docker/api/types"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	swarmapi "github.com/docker/swarmkit/api"
)

// GetSecret returns a secret from a managed swarm cluster
func (c *Cluster) GetSecret(input string) (types.Secret, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return types.Secret{}, c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	secret, err := getSecret(ctx, state.controlClient, input)
	if err != nil {
		return types.Secret{}, err
	}
	return convert.SecretFromGRPC(secret), nil
}

// GetSecrets returns all secrets of a managed swarm cluster.
func (c *Cluster) GetSecrets(options apitypes.SecretListOptions) ([]types.Secret, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return nil, c.errNoManager(state)
	}

	filters, err := newListSecretsFilters(options.Filters)
	if err != nil {
		return nil, err
	}
	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListSecrets(ctx,
		&swarmapi.ListSecretsRequest{Filters: filters})
	if err != nil {
		return nil, err
	}

	secrets := []types.Secret{}

	for _, secret := range r.Secrets {
		secrets = append(secrets, convert.SecretFromGRPC(secret))
	}

	return secrets, nil
}

// CreateSecret creates a new secret in a managed swarm cluster.
func (c *Cluster) CreateSecret(s types.SecretSpec) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return "", c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	secretSpec := convert.SecretSpecToGRPC(s)

	r, err := state.controlClient.CreateSecret(ctx,
		&swarmapi.CreateSecretRequest{Spec: &secretSpec})
	if err != nil {
		return "", err
	}

	return r.Secret.ID, nil
}

// RemoveSecret removes a secret from a managed swarm cluster.
func (c *Cluster) RemoveSecret(input string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	secret, err := getSecret(ctx, state.controlClient, input)
	if err != nil {
		return err
	}

	req := &swarmapi.RemoveSecretRequest{
		SecretID: secret.ID,
	}

	_, err = state.controlClient.RemoveSecret(ctx, req)
	return err
}

// UpdateSecret updates a secret in a managed swarm cluster.
// Note: this is not exposed to the CLI but is available from the API only
func (c *Cluster) UpdateSecret(id string, version uint64, spec types.SecretSpec) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() {
		return c.errNoManager(state)
	}

	ctx, cancel := c.getRequestContext()
	defer cancel()

	secretSpec := convert.SecretSpecToGRPC(spec)

	_, err := state.controlClient.UpdateSecret(ctx,
		&swarmapi.UpdateSecretRequest{
			SecretID: id,
			SecretVersion: &swarmapi.Version{
				Index: version,
			},
			Spec: &secretSpec,
		})
	return err
}
