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
	RESTART_POLICY_TYPE = "restartPolicy"
)

type RestartPolicy struct {
	Resource

	MaximumRetryCount int64 `json:"maximumRetryCount,omitempty" yaml:"maximum_retry_count,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type RestartPolicyCollection struct {
	Collection
	Data []RestartPolicy `json:"data,omitempty"`
}

type RestartPolicyClient struct {
	rancherClient *RancherClient
}

type RestartPolicyOperations interface {
	List(opts *ListOpts) (*RestartPolicyCollection, error)
	Create(opts *RestartPolicy) (*RestartPolicy, error)
	Update(existing *RestartPolicy, updates interface{}) (*RestartPolicy, error)
	ById(id string) (*RestartPolicy, error)
	Delete(container *RestartPolicy) error
}

func newRestartPolicyClient(rancherClient *RancherClient) *RestartPolicyClient {
	return &RestartPolicyClient{
		rancherClient: rancherClient,
	}
}

func (c *RestartPolicyClient) Create(container *RestartPolicy) (*RestartPolicy, error) {
	resp := &RestartPolicy{}
	err := c.rancherClient.doCreate(RESTART_POLICY_TYPE, container, resp)
	return resp, err
}

func (c *RestartPolicyClient) Update(existing *RestartPolicy, updates interface{}) (*RestartPolicy, error) {
	resp := &RestartPolicy{}
	err := c.rancherClient.doUpdate(RESTART_POLICY_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *RestartPolicyClient) List(opts *ListOpts) (*RestartPolicyCollection, error) {
	resp := &RestartPolicyCollection{}
	err := c.rancherClient.doList(RESTART_POLICY_TYPE, opts, resp)
	return resp, err
}

func (c *RestartPolicyClient) ById(id string) (*RestartPolicy, error) {
	resp := &RestartPolicy{}
	err := c.rancherClient.doById(RESTART_POLICY_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *RestartPolicyClient) Delete(container *RestartPolicy) error {
	return c.rancherClient.doResourceDelete(RESTART_POLICY_TYPE, &container.Resource)
}
