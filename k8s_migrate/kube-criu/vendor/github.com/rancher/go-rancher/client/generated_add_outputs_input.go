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
	ADD_OUTPUTS_INPUT_TYPE = "addOutputsInput"
)

type AddOutputsInput struct {
	Resource

	Outputs map[string]interface{} `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

type AddOutputsInputCollection struct {
	Collection
	Data []AddOutputsInput `json:"data,omitempty"`
}

type AddOutputsInputClient struct {
	rancherClient *RancherClient
}

type AddOutputsInputOperations interface {
	List(opts *ListOpts) (*AddOutputsInputCollection, error)
	Create(opts *AddOutputsInput) (*AddOutputsInput, error)
	Update(existing *AddOutputsInput, updates interface{}) (*AddOutputsInput, error)
	ById(id string) (*AddOutputsInput, error)
	Delete(container *AddOutputsInput) error
}

func newAddOutputsInputClient(rancherClient *RancherClient) *AddOutputsInputClient {
	return &AddOutputsInputClient{
		rancherClient: rancherClient,
	}
}

func (c *AddOutputsInputClient) Create(container *AddOutputsInput) (*AddOutputsInput, error) {
	resp := &AddOutputsInput{}
	err := c.rancherClient.doCreate(ADD_OUTPUTS_INPUT_TYPE, container, resp)
	return resp, err
}

func (c *AddOutputsInputClient) Update(existing *AddOutputsInput, updates interface{}) (*AddOutputsInput, error) {
	resp := &AddOutputsInput{}
	err := c.rancherClient.doUpdate(ADD_OUTPUTS_INPUT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *AddOutputsInputClient) List(opts *ListOpts) (*AddOutputsInputCollection, error) {
	resp := &AddOutputsInputCollection{}
	err := c.rancherClient.doList(ADD_OUTPUTS_INPUT_TYPE, opts, resp)
	return resp, err
}

func (c *AddOutputsInputClient) ById(id string) (*AddOutputsInput, error) {
	resp := &AddOutputsInput{}
	err := c.rancherClient.doById(ADD_OUTPUTS_INPUT_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *AddOutputsInputClient) Delete(container *AddOutputsInput) error {
	return c.rancherClient.doResourceDelete(ADD_OUTPUTS_INPUT_TYPE, &container.Resource)
}
