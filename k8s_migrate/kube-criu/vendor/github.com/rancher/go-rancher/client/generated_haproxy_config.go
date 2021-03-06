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
	HAPROXY_CONFIG_TYPE = "haproxyConfig"
)

type HaproxyConfig struct {
	Resource

	Defaults string `json:"defaults,omitempty" yaml:"defaults,omitempty"`

	Global string `json:"global,omitempty" yaml:"global,omitempty"`
}

type HaproxyConfigCollection struct {
	Collection
	Data []HaproxyConfig `json:"data,omitempty"`
}

type HaproxyConfigClient struct {
	rancherClient *RancherClient
}

type HaproxyConfigOperations interface {
	List(opts *ListOpts) (*HaproxyConfigCollection, error)
	Create(opts *HaproxyConfig) (*HaproxyConfig, error)
	Update(existing *HaproxyConfig, updates interface{}) (*HaproxyConfig, error)
	ById(id string) (*HaproxyConfig, error)
	Delete(container *HaproxyConfig) error
}

func newHaproxyConfigClient(rancherClient *RancherClient) *HaproxyConfigClient {
	return &HaproxyConfigClient{
		rancherClient: rancherClient,
	}
}

func (c *HaproxyConfigClient) Create(container *HaproxyConfig) (*HaproxyConfig, error) {
	resp := &HaproxyConfig{}
	err := c.rancherClient.doCreate(HAPROXY_CONFIG_TYPE, container, resp)
	return resp, err
}

func (c *HaproxyConfigClient) Update(existing *HaproxyConfig, updates interface{}) (*HaproxyConfig, error) {
	resp := &HaproxyConfig{}
	err := c.rancherClient.doUpdate(HAPROXY_CONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *HaproxyConfigClient) List(opts *ListOpts) (*HaproxyConfigCollection, error) {
	resp := &HaproxyConfigCollection{}
	err := c.rancherClient.doList(HAPROXY_CONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *HaproxyConfigClient) ById(id string) (*HaproxyConfig, error) {
	resp := &HaproxyConfig{}
	err := c.rancherClient.doById(HAPROXY_CONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *HaproxyConfigClient) Delete(container *HaproxyConfig) error {
	return c.rancherClient.doResourceDelete(HAPROXY_CONFIG_TYPE, &container.Resource)
}
