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
	HOST_ACCESS_TYPE = "hostAccess"
)

type HostAccess struct {
	Resource

	Token string `json:"token,omitempty" yaml:"token,omitempty"`

	Url string `json:"url,omitempty" yaml:"url,omitempty"`
}

type HostAccessCollection struct {
	Collection
	Data []HostAccess `json:"data,omitempty"`
}

type HostAccessClient struct {
	rancherClient *RancherClient
}

type HostAccessOperations interface {
	List(opts *ListOpts) (*HostAccessCollection, error)
	Create(opts *HostAccess) (*HostAccess, error)
	Update(existing *HostAccess, updates interface{}) (*HostAccess, error)
	ById(id string) (*HostAccess, error)
	Delete(container *HostAccess) error
}

func newHostAccessClient(rancherClient *RancherClient) *HostAccessClient {
	return &HostAccessClient{
		rancherClient: rancherClient,
	}
}

func (c *HostAccessClient) Create(container *HostAccess) (*HostAccess, error) {
	resp := &HostAccess{}
	err := c.rancherClient.doCreate(HOST_ACCESS_TYPE, container, resp)
	return resp, err
}

func (c *HostAccessClient) Update(existing *HostAccess, updates interface{}) (*HostAccess, error) {
	resp := &HostAccess{}
	err := c.rancherClient.doUpdate(HOST_ACCESS_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *HostAccessClient) List(opts *ListOpts) (*HostAccessCollection, error) {
	resp := &HostAccessCollection{}
	err := c.rancherClient.doList(HOST_ACCESS_TYPE, opts, resp)
	return resp, err
}

func (c *HostAccessClient) ById(id string) (*HostAccess, error) {
	resp := &HostAccess{}
	err := c.rancherClient.doById(HOST_ACCESS_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *HostAccessClient) Delete(container *HostAccess) error {
	return c.rancherClient.doResourceDelete(HOST_ACCESS_TYPE, &container.Resource)
}
