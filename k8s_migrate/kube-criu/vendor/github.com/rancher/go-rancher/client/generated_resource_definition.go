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
	RESOURCE_DEFINITION_TYPE = "resourceDefinition"
)

type ResourceDefinition struct {
	Resource

	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ResourceDefinitionCollection struct {
	Collection
	Data []ResourceDefinition `json:"data,omitempty"`
}

type ResourceDefinitionClient struct {
	rancherClient *RancherClient
}

type ResourceDefinitionOperations interface {
	List(opts *ListOpts) (*ResourceDefinitionCollection, error)
	Create(opts *ResourceDefinition) (*ResourceDefinition, error)
	Update(existing *ResourceDefinition, updates interface{}) (*ResourceDefinition, error)
	ById(id string) (*ResourceDefinition, error)
	Delete(container *ResourceDefinition) error
}

func newResourceDefinitionClient(rancherClient *RancherClient) *ResourceDefinitionClient {
	return &ResourceDefinitionClient{
		rancherClient: rancherClient,
	}
}

func (c *ResourceDefinitionClient) Create(container *ResourceDefinition) (*ResourceDefinition, error) {
	resp := &ResourceDefinition{}
	err := c.rancherClient.doCreate(RESOURCE_DEFINITION_TYPE, container, resp)
	return resp, err
}

func (c *ResourceDefinitionClient) Update(existing *ResourceDefinition, updates interface{}) (*ResourceDefinition, error) {
	resp := &ResourceDefinition{}
	err := c.rancherClient.doUpdate(RESOURCE_DEFINITION_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ResourceDefinitionClient) List(opts *ListOpts) (*ResourceDefinitionCollection, error) {
	resp := &ResourceDefinitionCollection{}
	err := c.rancherClient.doList(RESOURCE_DEFINITION_TYPE, opts, resp)
	return resp, err
}

func (c *ResourceDefinitionClient) ById(id string) (*ResourceDefinition, error) {
	resp := &ResourceDefinition{}
	err := c.rancherClient.doById(RESOURCE_DEFINITION_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ResourceDefinitionClient) Delete(container *ResourceDefinition) error {
	return c.rancherClient.doResourceDelete(RESOURCE_DEFINITION_TYPE, &container.Resource)
}
