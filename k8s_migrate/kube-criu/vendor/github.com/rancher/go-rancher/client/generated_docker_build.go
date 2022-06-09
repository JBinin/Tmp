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
	DOCKER_BUILD_TYPE = "dockerBuild"
)

type DockerBuild struct {
	Resource

	Context string `json:"context,omitempty" yaml:"context,omitempty"`

	Dockerfile string `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`

	Forcerm bool `json:"forcerm,omitempty" yaml:"forcerm,omitempty"`

	Nocache bool `json:"nocache,omitempty" yaml:"nocache,omitempty"`

	Remote string `json:"remote,omitempty" yaml:"remote,omitempty"`

	Rm bool `json:"rm,omitempty" yaml:"rm,omitempty"`
}

type DockerBuildCollection struct {
	Collection
	Data []DockerBuild `json:"data,omitempty"`
}

type DockerBuildClient struct {
	rancherClient *RancherClient
}

type DockerBuildOperations interface {
	List(opts *ListOpts) (*DockerBuildCollection, error)
	Create(opts *DockerBuild) (*DockerBuild, error)
	Update(existing *DockerBuild, updates interface{}) (*DockerBuild, error)
	ById(id string) (*DockerBuild, error)
	Delete(container *DockerBuild) error
}

func newDockerBuildClient(rancherClient *RancherClient) *DockerBuildClient {
	return &DockerBuildClient{
		rancherClient: rancherClient,
	}
}

func (c *DockerBuildClient) Create(container *DockerBuild) (*DockerBuild, error) {
	resp := &DockerBuild{}
	err := c.rancherClient.doCreate(DOCKER_BUILD_TYPE, container, resp)
	return resp, err
}

func (c *DockerBuildClient) Update(existing *DockerBuild, updates interface{}) (*DockerBuild, error) {
	resp := &DockerBuild{}
	err := c.rancherClient.doUpdate(DOCKER_BUILD_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *DockerBuildClient) List(opts *ListOpts) (*DockerBuildCollection, error) {
	resp := &DockerBuildCollection{}
	err := c.rancherClient.doList(DOCKER_BUILD_TYPE, opts, resp)
	return resp, err
}

func (c *DockerBuildClient) ById(id string) (*DockerBuild, error) {
	resp := &DockerBuild{}
	err := c.rancherClient.doById(DOCKER_BUILD_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *DockerBuildClient) Delete(container *DockerBuild) error {
	return c.rancherClient.doResourceDelete(DOCKER_BUILD_TYPE, &container.Resource)
}
