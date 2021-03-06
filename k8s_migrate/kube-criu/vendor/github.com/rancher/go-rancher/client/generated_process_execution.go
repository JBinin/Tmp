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
	PROCESS_EXECUTION_TYPE = "processExecution"
)

type ProcessExecution struct {
	Resource

	Created string `json:"created,omitempty" yaml:"created,omitempty"`

	Log map[string]interface{} `json:"log,omitempty" yaml:"log,omitempty"`

	ProcessInstanceId string `json:"processInstanceId,omitempty" yaml:"process_instance_id,omitempty"`

	Uuid string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ProcessExecutionCollection struct {
	Collection
	Data []ProcessExecution `json:"data,omitempty"`
}

type ProcessExecutionClient struct {
	rancherClient *RancherClient
}

type ProcessExecutionOperations interface {
	List(opts *ListOpts) (*ProcessExecutionCollection, error)
	Create(opts *ProcessExecution) (*ProcessExecution, error)
	Update(existing *ProcessExecution, updates interface{}) (*ProcessExecution, error)
	ById(id string) (*ProcessExecution, error)
	Delete(container *ProcessExecution) error
}

func newProcessExecutionClient(rancherClient *RancherClient) *ProcessExecutionClient {
	return &ProcessExecutionClient{
		rancherClient: rancherClient,
	}
}

func (c *ProcessExecutionClient) Create(container *ProcessExecution) (*ProcessExecution, error) {
	resp := &ProcessExecution{}
	err := c.rancherClient.doCreate(PROCESS_EXECUTION_TYPE, container, resp)
	return resp, err
}

func (c *ProcessExecutionClient) Update(existing *ProcessExecution, updates interface{}) (*ProcessExecution, error) {
	resp := &ProcessExecution{}
	err := c.rancherClient.doUpdate(PROCESS_EXECUTION_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ProcessExecutionClient) List(opts *ListOpts) (*ProcessExecutionCollection, error) {
	resp := &ProcessExecutionCollection{}
	err := c.rancherClient.doList(PROCESS_EXECUTION_TYPE, opts, resp)
	return resp, err
}

func (c *ProcessExecutionClient) ById(id string) (*ProcessExecution, error) {
	resp := &ProcessExecution{}
	err := c.rancherClient.doById(PROCESS_EXECUTION_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ProcessExecutionClient) Delete(container *ProcessExecution) error {
	return c.rancherClient.doResourceDelete(PROCESS_EXECUTION_TYPE, &container.Resource)
}
