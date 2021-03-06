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
package client

const (
	HA_CONFIG_TYPE = "haConfig"
)

type HaConfig struct {
	Resource

	ClusterSize int64 `json:"clusterSize,omitempty" yaml:"cluster_size,omitempty"`

	DbHost string `json:"dbHost,omitempty" yaml:"db_host,omitempty"`

	DbSize int64 `json:"dbSize,omitempty" yaml:"db_size,omitempty"`

	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type HaConfigCollection struct {
	Collection
	Data []HaConfig `json:"data,omitempty"`
}

type HaConfigClient struct {
	rancherClient *RancherClient
}

type HaConfigOperations interface {
	List(opts *ListOpts) (*HaConfigCollection, error)
	Create(opts *HaConfig) (*HaConfig, error)
	Update(existing *HaConfig, updates interface{}) (*HaConfig, error)
	ById(id string) (*HaConfig, error)
	Delete(container *HaConfig) error
}

func newHaConfigClient(rancherClient *RancherClient) *HaConfigClient {
	return &HaConfigClient{
		rancherClient: rancherClient,
	}
}

func (c *HaConfigClient) Create(container *HaConfig) (*HaConfig, error) {
	resp := &HaConfig{}
	err := c.rancherClient.doCreate(HA_CONFIG_TYPE, container, resp)
	return resp, err
}

func (c *HaConfigClient) Update(existing *HaConfig, updates interface{}) (*HaConfig, error) {
	resp := &HaConfig{}
	err := c.rancherClient.doUpdate(HA_CONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *HaConfigClient) List(opts *ListOpts) (*HaConfigCollection, error) {
	resp := &HaConfigCollection{}
	err := c.rancherClient.doList(HA_CONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *HaConfigClient) ById(id string) (*HaConfig, error) {
	resp := &HaConfig{}
	err := c.rancherClient.doById(HA_CONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *HaConfigClient) Delete(container *HaConfig) error {
	return c.rancherClient.doResourceDelete(HA_CONFIG_TYPE, &container.Resource)
}
