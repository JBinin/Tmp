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
	SET_PROJECT_MEMBERS_INPUT_TYPE = "setProjectMembersInput"
)

type SetProjectMembersInput struct {
	Resource

	Members []interface{} `json:"members,omitempty" yaml:"members,omitempty"`
}

type SetProjectMembersInputCollection struct {
	Collection
	Data []SetProjectMembersInput `json:"data,omitempty"`
}

type SetProjectMembersInputClient struct {
	rancherClient *RancherClient
}

type SetProjectMembersInputOperations interface {
	List(opts *ListOpts) (*SetProjectMembersInputCollection, error)
	Create(opts *SetProjectMembersInput) (*SetProjectMembersInput, error)
	Update(existing *SetProjectMembersInput, updates interface{}) (*SetProjectMembersInput, error)
	ById(id string) (*SetProjectMembersInput, error)
	Delete(container *SetProjectMembersInput) error
}

func newSetProjectMembersInputClient(rancherClient *RancherClient) *SetProjectMembersInputClient {
	return &SetProjectMembersInputClient{
		rancherClient: rancherClient,
	}
}

func (c *SetProjectMembersInputClient) Create(container *SetProjectMembersInput) (*SetProjectMembersInput, error) {
	resp := &SetProjectMembersInput{}
	err := c.rancherClient.doCreate(SET_PROJECT_MEMBERS_INPUT_TYPE, container, resp)
	return resp, err
}

func (c *SetProjectMembersInputClient) Update(existing *SetProjectMembersInput, updates interface{}) (*SetProjectMembersInput, error) {
	resp := &SetProjectMembersInput{}
	err := c.rancherClient.doUpdate(SET_PROJECT_MEMBERS_INPUT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *SetProjectMembersInputClient) List(opts *ListOpts) (*SetProjectMembersInputCollection, error) {
	resp := &SetProjectMembersInputCollection{}
	err := c.rancherClient.doList(SET_PROJECT_MEMBERS_INPUT_TYPE, opts, resp)
	return resp, err
}

func (c *SetProjectMembersInputClient) ById(id string) (*SetProjectMembersInput, error) {
	resp := &SetProjectMembersInput{}
	err := c.rancherClient.doById(SET_PROJECT_MEMBERS_INPUT_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *SetProjectMembersInputClient) Delete(container *SetProjectMembersInput) error {
	return c.rancherClient.doResourceDelete(SET_PROJECT_MEMBERS_INPUT_TYPE, &container.Resource)
}
