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
	STATE_TRANSITION_TYPE = "stateTransition"
)

type StateTransition struct {
	Resource
}

type StateTransitionCollection struct {
	Collection
	Data []StateTransition `json:"data,omitempty"`
}

type StateTransitionClient struct {
	rancherClient *RancherClient
}

type StateTransitionOperations interface {
	List(opts *ListOpts) (*StateTransitionCollection, error)
	Create(opts *StateTransition) (*StateTransition, error)
	Update(existing *StateTransition, updates interface{}) (*StateTransition, error)
	ById(id string) (*StateTransition, error)
	Delete(container *StateTransition) error
}

func newStateTransitionClient(rancherClient *RancherClient) *StateTransitionClient {
	return &StateTransitionClient{
		rancherClient: rancherClient,
	}
}

func (c *StateTransitionClient) Create(container *StateTransition) (*StateTransition, error) {
	resp := &StateTransition{}
	err := c.rancherClient.doCreate(STATE_TRANSITION_TYPE, container, resp)
	return resp, err
}

func (c *StateTransitionClient) Update(existing *StateTransition, updates interface{}) (*StateTransition, error) {
	resp := &StateTransition{}
	err := c.rancherClient.doUpdate(STATE_TRANSITION_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *StateTransitionClient) List(opts *ListOpts) (*StateTransitionCollection, error) {
	resp := &StateTransitionCollection{}
	err := c.rancherClient.doList(STATE_TRANSITION_TYPE, opts, resp)
	return resp, err
}

func (c *StateTransitionClient) ById(id string) (*StateTransition, error) {
	resp := &StateTransition{}
	err := c.rancherClient.doById(STATE_TRANSITION_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *StateTransitionClient) Delete(container *StateTransition) error {
	return c.rancherClient.doResourceDelete(STATE_TRANSITION_TYPE, &container.Resource)
}
