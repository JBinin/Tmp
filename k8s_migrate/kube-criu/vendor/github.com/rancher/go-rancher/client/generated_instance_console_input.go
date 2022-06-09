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
	INSTANCE_CONSOLE_INPUT_TYPE = "instanceConsoleInput"
)

type InstanceConsoleInput struct {
	Resource
}

type InstanceConsoleInputCollection struct {
	Collection
	Data []InstanceConsoleInput `json:"data,omitempty"`
}

type InstanceConsoleInputClient struct {
	rancherClient *RancherClient
}

type InstanceConsoleInputOperations interface {
	List(opts *ListOpts) (*InstanceConsoleInputCollection, error)
	Create(opts *InstanceConsoleInput) (*InstanceConsoleInput, error)
	Update(existing *InstanceConsoleInput, updates interface{}) (*InstanceConsoleInput, error)
	ById(id string) (*InstanceConsoleInput, error)
	Delete(container *InstanceConsoleInput) error
}

func newInstanceConsoleInputClient(rancherClient *RancherClient) *InstanceConsoleInputClient {
	return &InstanceConsoleInputClient{
		rancherClient: rancherClient,
	}
}

func (c *InstanceConsoleInputClient) Create(container *InstanceConsoleInput) (*InstanceConsoleInput, error) {
	resp := &InstanceConsoleInput{}
	err := c.rancherClient.doCreate(INSTANCE_CONSOLE_INPUT_TYPE, container, resp)
	return resp, err
}

func (c *InstanceConsoleInputClient) Update(existing *InstanceConsoleInput, updates interface{}) (*InstanceConsoleInput, error) {
	resp := &InstanceConsoleInput{}
	err := c.rancherClient.doUpdate(INSTANCE_CONSOLE_INPUT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *InstanceConsoleInputClient) List(opts *ListOpts) (*InstanceConsoleInputCollection, error) {
	resp := &InstanceConsoleInputCollection{}
	err := c.rancherClient.doList(INSTANCE_CONSOLE_INPUT_TYPE, opts, resp)
	return resp, err
}

func (c *InstanceConsoleInputClient) ById(id string) (*InstanceConsoleInput, error) {
	resp := &InstanceConsoleInput{}
	err := c.rancherClient.doById(INSTANCE_CONSOLE_INPUT_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *InstanceConsoleInputClient) Delete(container *InstanceConsoleInput) error {
	return c.rancherClient.doResourceDelete(INSTANCE_CONSOLE_INPUT_TYPE, &container.Resource)
}
