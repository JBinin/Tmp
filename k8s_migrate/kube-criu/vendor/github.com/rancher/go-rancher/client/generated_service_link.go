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
	SERVICE_LINK_TYPE = "serviceLink"
)

type ServiceLink struct {
	Resource

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	ServiceId string `json:"serviceId,omitempty" yaml:"service_id,omitempty"`

	Uuid string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ServiceLinkCollection struct {
	Collection
	Data []ServiceLink `json:"data,omitempty"`
}

type ServiceLinkClient struct {
	rancherClient *RancherClient
}

type ServiceLinkOperations interface {
	List(opts *ListOpts) (*ServiceLinkCollection, error)
	Create(opts *ServiceLink) (*ServiceLink, error)
	Update(existing *ServiceLink, updates interface{}) (*ServiceLink, error)
	ById(id string) (*ServiceLink, error)
	Delete(container *ServiceLink) error
}

func newServiceLinkClient(rancherClient *RancherClient) *ServiceLinkClient {
	return &ServiceLinkClient{
		rancherClient: rancherClient,
	}
}

func (c *ServiceLinkClient) Create(container *ServiceLink) (*ServiceLink, error) {
	resp := &ServiceLink{}
	err := c.rancherClient.doCreate(SERVICE_LINK_TYPE, container, resp)
	return resp, err
}

func (c *ServiceLinkClient) Update(existing *ServiceLink, updates interface{}) (*ServiceLink, error) {
	resp := &ServiceLink{}
	err := c.rancherClient.doUpdate(SERVICE_LINK_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ServiceLinkClient) List(opts *ListOpts) (*ServiceLinkCollection, error) {
	resp := &ServiceLinkCollection{}
	err := c.rancherClient.doList(SERVICE_LINK_TYPE, opts, resp)
	return resp, err
}

func (c *ServiceLinkClient) ById(id string) (*ServiceLink, error) {
	resp := &ServiceLink{}
	err := c.rancherClient.doById(SERVICE_LINK_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ServiceLinkClient) Delete(container *ServiceLink) error {
	return c.rancherClient.doResourceDelete(SERVICE_LINK_TYPE, &container.Resource)
}
