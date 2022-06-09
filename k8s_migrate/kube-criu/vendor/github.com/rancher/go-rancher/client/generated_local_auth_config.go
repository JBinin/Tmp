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
	LOCAL_AUTH_CONFIG_TYPE = "localAuthConfig"
)

type LocalAuthConfig struct {
	Resource

	AccessMode string `json:"accessMode,omitempty" yaml:"access_mode,omitempty"`

	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	Username string `json:"username,omitempty" yaml:"username,omitempty"`
}

type LocalAuthConfigCollection struct {
	Collection
	Data []LocalAuthConfig `json:"data,omitempty"`
}

type LocalAuthConfigClient struct {
	rancherClient *RancherClient
}

type LocalAuthConfigOperations interface {
	List(opts *ListOpts) (*LocalAuthConfigCollection, error)
	Create(opts *LocalAuthConfig) (*LocalAuthConfig, error)
	Update(existing *LocalAuthConfig, updates interface{}) (*LocalAuthConfig, error)
	ById(id string) (*LocalAuthConfig, error)
	Delete(container *LocalAuthConfig) error
}

func newLocalAuthConfigClient(rancherClient *RancherClient) *LocalAuthConfigClient {
	return &LocalAuthConfigClient{
		rancherClient: rancherClient,
	}
}

func (c *LocalAuthConfigClient) Create(container *LocalAuthConfig) (*LocalAuthConfig, error) {
	resp := &LocalAuthConfig{}
	err := c.rancherClient.doCreate(LOCAL_AUTH_CONFIG_TYPE, container, resp)
	return resp, err
}

func (c *LocalAuthConfigClient) Update(existing *LocalAuthConfig, updates interface{}) (*LocalAuthConfig, error) {
	resp := &LocalAuthConfig{}
	err := c.rancherClient.doUpdate(LOCAL_AUTH_CONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *LocalAuthConfigClient) List(opts *ListOpts) (*LocalAuthConfigCollection, error) {
	resp := &LocalAuthConfigCollection{}
	err := c.rancherClient.doList(LOCAL_AUTH_CONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *LocalAuthConfigClient) ById(id string) (*LocalAuthConfig, error) {
	resp := &LocalAuthConfig{}
	err := c.rancherClient.doById(LOCAL_AUTH_CONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *LocalAuthConfigClient) Delete(container *LocalAuthConfig) error {
	return c.rancherClient.doResourceDelete(LOCAL_AUTH_CONFIG_TYPE, &container.Resource)
}
