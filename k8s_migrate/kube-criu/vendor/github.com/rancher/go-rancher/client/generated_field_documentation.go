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
	FIELD_DOCUMENTATION_TYPE = "fieldDocumentation"
)

type FieldDocumentation struct {
	Resource

	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type FieldDocumentationCollection struct {
	Collection
	Data []FieldDocumentation `json:"data,omitempty"`
}

type FieldDocumentationClient struct {
	rancherClient *RancherClient
}

type FieldDocumentationOperations interface {
	List(opts *ListOpts) (*FieldDocumentationCollection, error)
	Create(opts *FieldDocumentation) (*FieldDocumentation, error)
	Update(existing *FieldDocumentation, updates interface{}) (*FieldDocumentation, error)
	ById(id string) (*FieldDocumentation, error)
	Delete(container *FieldDocumentation) error
}

func newFieldDocumentationClient(rancherClient *RancherClient) *FieldDocumentationClient {
	return &FieldDocumentationClient{
		rancherClient: rancherClient,
	}
}

func (c *FieldDocumentationClient) Create(container *FieldDocumentation) (*FieldDocumentation, error) {
	resp := &FieldDocumentation{}
	err := c.rancherClient.doCreate(FIELD_DOCUMENTATION_TYPE, container, resp)
	return resp, err
}

func (c *FieldDocumentationClient) Update(existing *FieldDocumentation, updates interface{}) (*FieldDocumentation, error) {
	resp := &FieldDocumentation{}
	err := c.rancherClient.doUpdate(FIELD_DOCUMENTATION_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *FieldDocumentationClient) List(opts *ListOpts) (*FieldDocumentationCollection, error) {
	resp := &FieldDocumentationCollection{}
	err := c.rancherClient.doList(FIELD_DOCUMENTATION_TYPE, opts, resp)
	return resp, err
}

func (c *FieldDocumentationClient) ById(id string) (*FieldDocumentation, error) {
	resp := &FieldDocumentation{}
	err := c.rancherClient.doById(FIELD_DOCUMENTATION_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *FieldDocumentationClient) Delete(container *FieldDocumentation) error {
	return c.rancherClient.doResourceDelete(FIELD_DOCUMENTATION_TYPE, &container.Resource)
}
