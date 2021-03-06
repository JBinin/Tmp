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
	COMPOSE_CONFIG_INPUT_TYPE = "composeConfigInput"
)

type ComposeConfigInput struct {
	Resource

	ServiceIds []string `json:"serviceIds,omitempty" yaml:"service_ids,omitempty"`
}

type ComposeConfigInputCollection struct {
	Collection
	Data []ComposeConfigInput `json:"data,omitempty"`
}

type ComposeConfigInputClient struct {
	rancherClient *RancherClient
}

type ComposeConfigInputOperations interface {
	List(opts *ListOpts) (*ComposeConfigInputCollection, error)
	Create(opts *ComposeConfigInput) (*ComposeConfigInput, error)
	Update(existing *ComposeConfigInput, updates interface{}) (*ComposeConfigInput, error)
	ById(id string) (*ComposeConfigInput, error)
	Delete(container *ComposeConfigInput) error
}

func newComposeConfigInputClient(rancherClient *RancherClient) *ComposeConfigInputClient {
	return &ComposeConfigInputClient{
		rancherClient: rancherClient,
	}
}

func (c *ComposeConfigInputClient) Create(container *ComposeConfigInput) (*ComposeConfigInput, error) {
	resp := &ComposeConfigInput{}
	err := c.rancherClient.doCreate(COMPOSE_CONFIG_INPUT_TYPE, container, resp)
	return resp, err
}

func (c *ComposeConfigInputClient) Update(existing *ComposeConfigInput, updates interface{}) (*ComposeConfigInput, error) {
	resp := &ComposeConfigInput{}
	err := c.rancherClient.doUpdate(COMPOSE_CONFIG_INPUT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ComposeConfigInputClient) List(opts *ListOpts) (*ComposeConfigInputCollection, error) {
	resp := &ComposeConfigInputCollection{}
	err := c.rancherClient.doList(COMPOSE_CONFIG_INPUT_TYPE, opts, resp)
	return resp, err
}

func (c *ComposeConfigInputClient) ById(id string) (*ComposeConfigInput, error) {
	resp := &ComposeConfigInput{}
	err := c.rancherClient.doById(COMPOSE_CONFIG_INPUT_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ComposeConfigInputClient) Delete(container *ComposeConfigInput) error {
	return c.rancherClient.doResourceDelete(COMPOSE_CONFIG_INPUT_TYPE, &container.Resource)
}
