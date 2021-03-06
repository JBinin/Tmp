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
	HOST_API_PROXY_TOKEN_TYPE = "hostApiProxyToken"
)

type HostApiProxyToken struct {
	Resource

	ReportedUuid string `json:"reportedUuid,omitempty" yaml:"reported_uuid,omitempty"`

	Token string `json:"token,omitempty" yaml:"token,omitempty"`

	Url string `json:"url,omitempty" yaml:"url,omitempty"`
}

type HostApiProxyTokenCollection struct {
	Collection
	Data []HostApiProxyToken `json:"data,omitempty"`
}

type HostApiProxyTokenClient struct {
	rancherClient *RancherClient
}

type HostApiProxyTokenOperations interface {
	List(opts *ListOpts) (*HostApiProxyTokenCollection, error)
	Create(opts *HostApiProxyToken) (*HostApiProxyToken, error)
	Update(existing *HostApiProxyToken, updates interface{}) (*HostApiProxyToken, error)
	ById(id string) (*HostApiProxyToken, error)
	Delete(container *HostApiProxyToken) error
}

func newHostApiProxyTokenClient(rancherClient *RancherClient) *HostApiProxyTokenClient {
	return &HostApiProxyTokenClient{
		rancherClient: rancherClient,
	}
}

func (c *HostApiProxyTokenClient) Create(container *HostApiProxyToken) (*HostApiProxyToken, error) {
	resp := &HostApiProxyToken{}
	err := c.rancherClient.doCreate(HOST_API_PROXY_TOKEN_TYPE, container, resp)
	return resp, err
}

func (c *HostApiProxyTokenClient) Update(existing *HostApiProxyToken, updates interface{}) (*HostApiProxyToken, error) {
	resp := &HostApiProxyToken{}
	err := c.rancherClient.doUpdate(HOST_API_PROXY_TOKEN_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *HostApiProxyTokenClient) List(opts *ListOpts) (*HostApiProxyTokenCollection, error) {
	resp := &HostApiProxyTokenCollection{}
	err := c.rancherClient.doList(HOST_API_PROXY_TOKEN_TYPE, opts, resp)
	return resp, err
}

func (c *HostApiProxyTokenClient) ById(id string) (*HostApiProxyToken, error) {
	resp := &HostApiProxyToken{}
	err := c.rancherClient.doById(HOST_API_PROXY_TOKEN_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *HostApiProxyTokenClient) Delete(container *HostApiProxyToken) error {
	return c.rancherClient.doResourceDelete(HOST_API_PROXY_TOKEN_TYPE, &container.Resource)
}
