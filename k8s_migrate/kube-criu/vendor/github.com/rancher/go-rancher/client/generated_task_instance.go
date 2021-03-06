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
	TASK_INSTANCE_TYPE = "taskInstance"
)

type TaskInstance struct {
	Resource

	EndTime string `json:"endTime,omitempty" yaml:"end_time,omitempty"`

	Exception string `json:"exception,omitempty" yaml:"exception,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	ServerId string `json:"serverId,omitempty" yaml:"server_id,omitempty"`

	StartTime string `json:"startTime,omitempty" yaml:"start_time,omitempty"`

	TaskId string `json:"taskId,omitempty" yaml:"task_id,omitempty"`
}

type TaskInstanceCollection struct {
	Collection
	Data []TaskInstance `json:"data,omitempty"`
}

type TaskInstanceClient struct {
	rancherClient *RancherClient
}

type TaskInstanceOperations interface {
	List(opts *ListOpts) (*TaskInstanceCollection, error)
	Create(opts *TaskInstance) (*TaskInstance, error)
	Update(existing *TaskInstance, updates interface{}) (*TaskInstance, error)
	ById(id string) (*TaskInstance, error)
	Delete(container *TaskInstance) error
}

func newTaskInstanceClient(rancherClient *RancherClient) *TaskInstanceClient {
	return &TaskInstanceClient{
		rancherClient: rancherClient,
	}
}

func (c *TaskInstanceClient) Create(container *TaskInstance) (*TaskInstance, error) {
	resp := &TaskInstance{}
	err := c.rancherClient.doCreate(TASK_INSTANCE_TYPE, container, resp)
	return resp, err
}

func (c *TaskInstanceClient) Update(existing *TaskInstance, updates interface{}) (*TaskInstance, error) {
	resp := &TaskInstance{}
	err := c.rancherClient.doUpdate(TASK_INSTANCE_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *TaskInstanceClient) List(opts *ListOpts) (*TaskInstanceCollection, error) {
	resp := &TaskInstanceCollection{}
	err := c.rancherClient.doList(TASK_INSTANCE_TYPE, opts, resp)
	return resp, err
}

func (c *TaskInstanceClient) ById(id string) (*TaskInstance, error) {
	resp := &TaskInstance{}
	err := c.rancherClient.doById(TASK_INSTANCE_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *TaskInstanceClient) Delete(container *TaskInstance) error {
	return c.rancherClient.doResourceDelete(TASK_INSTANCE_TYPE, &container.Resource)
}
