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
	CHANGE_SECRET_INPUT_TYPE = "changeSecretInput"
)

type ChangeSecretInput struct {
	Resource

	NewSecret string `json:"newSecret,omitempty" yaml:"new_secret,omitempty"`

	OldSecret string `json:"oldSecret,omitempty" yaml:"old_secret,omitempty"`
}

type ChangeSecretInputCollection struct {
	Collection
	Data []ChangeSecretInput `json:"data,omitempty"`
}

type ChangeSecretInputClient struct {
	rancherClient *RancherClient
}

type ChangeSecretInputOperations interface {
	List(opts *ListOpts) (*ChangeSecretInputCollection, error)
	Create(opts *ChangeSecretInput) (*ChangeSecretInput, error)
	Update(existing *ChangeSecretInput, updates interface{}) (*ChangeSecretInput, error)
	ById(id string) (*ChangeSecretInput, error)
	Delete(container *ChangeSecretInput) error
}

func newChangeSecretInputClient(rancherClient *RancherClient) *ChangeSecretInputClient {
	return &ChangeSecretInputClient{
		rancherClient: rancherClient,
	}
}

func (c *ChangeSecretInputClient) Create(container *ChangeSecretInput) (*ChangeSecretInput, error) {
	resp := &ChangeSecretInput{}
	err := c.rancherClient.doCreate(CHANGE_SECRET_INPUT_TYPE, container, resp)
	return resp, err
}

func (c *ChangeSecretInputClient) Update(existing *ChangeSecretInput, updates interface{}) (*ChangeSecretInput, error) {
	resp := &ChangeSecretInput{}
	err := c.rancherClient.doUpdate(CHANGE_SECRET_INPUT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ChangeSecretInputClient) List(opts *ListOpts) (*ChangeSecretInputCollection, error) {
	resp := &ChangeSecretInputCollection{}
	err := c.rancherClient.doList(CHANGE_SECRET_INPUT_TYPE, opts, resp)
	return resp, err
}

func (c *ChangeSecretInputClient) ById(id string) (*ChangeSecretInput, error) {
	resp := &ChangeSecretInput{}
	err := c.rancherClient.doById(CHANGE_SECRET_INPUT_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ChangeSecretInputClient) Delete(container *ChangeSecretInput) error {
	return c.rancherClient.doResourceDelete(CHANGE_SECRET_INPUT_TYPE, &container.Resource)
}
