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
	TYPE_DOCUMENTATION_TYPE = "typeDocumentation"
)

type TypeDocumentation struct {
	Resource

	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	ResourceFields map[string]interface{} `json:"resourceFields,omitempty" yaml:"resource_fields,omitempty"`
}

type TypeDocumentationCollection struct {
	Collection
	Data []TypeDocumentation `json:"data,omitempty"`
}

type TypeDocumentationClient struct {
	rancherClient *RancherClient
}

type TypeDocumentationOperations interface {
	List(opts *ListOpts) (*TypeDocumentationCollection, error)
	Create(opts *TypeDocumentation) (*TypeDocumentation, error)
	Update(existing *TypeDocumentation, updates interface{}) (*TypeDocumentation, error)
	ById(id string) (*TypeDocumentation, error)
	Delete(container *TypeDocumentation) error
}

func newTypeDocumentationClient(rancherClient *RancherClient) *TypeDocumentationClient {
	return &TypeDocumentationClient{
		rancherClient: rancherClient,
	}
}

func (c *TypeDocumentationClient) Create(container *TypeDocumentation) (*TypeDocumentation, error) {
	resp := &TypeDocumentation{}
	err := c.rancherClient.doCreate(TYPE_DOCUMENTATION_TYPE, container, resp)
	return resp, err
}

func (c *TypeDocumentationClient) Update(existing *TypeDocumentation, updates interface{}) (*TypeDocumentation, error) {
	resp := &TypeDocumentation{}
	err := c.rancherClient.doUpdate(TYPE_DOCUMENTATION_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *TypeDocumentationClient) List(opts *ListOpts) (*TypeDocumentationCollection, error) {
	resp := &TypeDocumentationCollection{}
	err := c.rancherClient.doList(TYPE_DOCUMENTATION_TYPE, opts, resp)
	return resp, err
}

func (c *TypeDocumentationClient) ById(id string) (*TypeDocumentation, error) {
	resp := &TypeDocumentation{}
	err := c.rancherClient.doById(TYPE_DOCUMENTATION_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *TypeDocumentationClient) Delete(container *TypeDocumentation) error {
	return c.rancherClient.doResourceDelete(TYPE_DOCUMENTATION_TYPE, &container.Resource)
}
