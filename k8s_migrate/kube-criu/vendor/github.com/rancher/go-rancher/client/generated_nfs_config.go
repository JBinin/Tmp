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
	NFS_CONFIG_TYPE = "nfsConfig"
)

type NfsConfig struct {
	Resource

	MountOptions string `json:"mountOptions,omitempty" yaml:"mount_options,omitempty"`

	Server string `json:"server,omitempty" yaml:"server,omitempty"`

	Share string `json:"share,omitempty" yaml:"share,omitempty"`
}

type NfsConfigCollection struct {
	Collection
	Data []NfsConfig `json:"data,omitempty"`
}

type NfsConfigClient struct {
	rancherClient *RancherClient
}

type NfsConfigOperations interface {
	List(opts *ListOpts) (*NfsConfigCollection, error)
	Create(opts *NfsConfig) (*NfsConfig, error)
	Update(existing *NfsConfig, updates interface{}) (*NfsConfig, error)
	ById(id string) (*NfsConfig, error)
	Delete(container *NfsConfig) error
}

func newNfsConfigClient(rancherClient *RancherClient) *NfsConfigClient {
	return &NfsConfigClient{
		rancherClient: rancherClient,
	}
}

func (c *NfsConfigClient) Create(container *NfsConfig) (*NfsConfig, error) {
	resp := &NfsConfig{}
	err := c.rancherClient.doCreate(NFS_CONFIG_TYPE, container, resp)
	return resp, err
}

func (c *NfsConfigClient) Update(existing *NfsConfig, updates interface{}) (*NfsConfig, error) {
	resp := &NfsConfig{}
	err := c.rancherClient.doUpdate(NFS_CONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *NfsConfigClient) List(opts *ListOpts) (*NfsConfigCollection, error) {
	resp := &NfsConfigCollection{}
	err := c.rancherClient.doList(NFS_CONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *NfsConfigClient) ById(id string) (*NfsConfig, error) {
	resp := &NfsConfig{}
	err := c.rancherClient.doById(NFS_CONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *NfsConfigClient) Delete(container *NfsConfig) error {
	return c.rancherClient.doResourceDelete(NFS_CONFIG_TYPE, &container.Resource)
}
