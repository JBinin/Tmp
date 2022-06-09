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
	LOG_CONFIG_TYPE = "logConfig"
)

type LogConfig struct {
	Resource

	Config map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`

	Driver string `json:"driver,omitempty" yaml:"driver,omitempty"`
}

type LogConfigCollection struct {
	Collection
	Data []LogConfig `json:"data,omitempty"`
}

type LogConfigClient struct {
	rancherClient *RancherClient
}

type LogConfigOperations interface {
	List(opts *ListOpts) (*LogConfigCollection, error)
	Create(opts *LogConfig) (*LogConfig, error)
	Update(existing *LogConfig, updates interface{}) (*LogConfig, error)
	ById(id string) (*LogConfig, error)
	Delete(container *LogConfig) error
}

func newLogConfigClient(rancherClient *RancherClient) *LogConfigClient {
	return &LogConfigClient{
		rancherClient: rancherClient,
	}
}

func (c *LogConfigClient) Create(container *LogConfig) (*LogConfig, error) {
	resp := &LogConfig{}
	err := c.rancherClient.doCreate(LOG_CONFIG_TYPE, container, resp)
	return resp, err
}

func (c *LogConfigClient) Update(existing *LogConfig, updates interface{}) (*LogConfig, error) {
	resp := &LogConfig{}
	err := c.rancherClient.doUpdate(LOG_CONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *LogConfigClient) List(opts *ListOpts) (*LogConfigCollection, error) {
	resp := &LogConfigCollection{}
	err := c.rancherClient.doList(LOG_CONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *LogConfigClient) ById(id string) (*LogConfig, error) {
	resp := &LogConfig{}
	err := c.rancherClient.doById(LOG_CONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *LogConfigClient) Delete(container *LogConfig) error {
	return c.rancherClient.doResourceDelete(LOG_CONFIG_TYPE, &container.Resource)
}
