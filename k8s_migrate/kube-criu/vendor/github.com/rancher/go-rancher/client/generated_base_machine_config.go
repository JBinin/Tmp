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
	BASE_MACHINE_CONFIG_TYPE = "baseMachineConfig"
)

type BaseMachineConfig struct {
	Resource
}

type BaseMachineConfigCollection struct {
	Collection
	Data []BaseMachineConfig `json:"data,omitempty"`
}

type BaseMachineConfigClient struct {
	rancherClient *RancherClient
}

type BaseMachineConfigOperations interface {
	List(opts *ListOpts) (*BaseMachineConfigCollection, error)
	Create(opts *BaseMachineConfig) (*BaseMachineConfig, error)
	Update(existing *BaseMachineConfig, updates interface{}) (*BaseMachineConfig, error)
	ById(id string) (*BaseMachineConfig, error)
	Delete(container *BaseMachineConfig) error
}

func newBaseMachineConfigClient(rancherClient *RancherClient) *BaseMachineConfigClient {
	return &BaseMachineConfigClient{
		rancherClient: rancherClient,
	}
}

func (c *BaseMachineConfigClient) Create(container *BaseMachineConfig) (*BaseMachineConfig, error) {
	resp := &BaseMachineConfig{}
	err := c.rancherClient.doCreate(BASE_MACHINE_CONFIG_TYPE, container, resp)
	return resp, err
}

func (c *BaseMachineConfigClient) Update(existing *BaseMachineConfig, updates interface{}) (*BaseMachineConfig, error) {
	resp := &BaseMachineConfig{}
	err := c.rancherClient.doUpdate(BASE_MACHINE_CONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BaseMachineConfigClient) List(opts *ListOpts) (*BaseMachineConfigCollection, error) {
	resp := &BaseMachineConfigCollection{}
	err := c.rancherClient.doList(BASE_MACHINE_CONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *BaseMachineConfigClient) ById(id string) (*BaseMachineConfig, error) {
	resp := &BaseMachineConfig{}
	err := c.rancherClient.doById(BASE_MACHINE_CONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *BaseMachineConfigClient) Delete(container *BaseMachineConfig) error {
	return c.rancherClient.doResourceDelete(BASE_MACHINE_CONFIG_TYPE, &container.Resource)
}
