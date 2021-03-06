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
	SERVICES_PORT_RANGE_TYPE = "servicesPortRange"
)

type ServicesPortRange struct {
	Resource

	EndPort int64 `json:"endPort,omitempty" yaml:"end_port,omitempty"`

	StartPort int64 `json:"startPort,omitempty" yaml:"start_port,omitempty"`
}

type ServicesPortRangeCollection struct {
	Collection
	Data []ServicesPortRange `json:"data,omitempty"`
}

type ServicesPortRangeClient struct {
	rancherClient *RancherClient
}

type ServicesPortRangeOperations interface {
	List(opts *ListOpts) (*ServicesPortRangeCollection, error)
	Create(opts *ServicesPortRange) (*ServicesPortRange, error)
	Update(existing *ServicesPortRange, updates interface{}) (*ServicesPortRange, error)
	ById(id string) (*ServicesPortRange, error)
	Delete(container *ServicesPortRange) error
}

func newServicesPortRangeClient(rancherClient *RancherClient) *ServicesPortRangeClient {
	return &ServicesPortRangeClient{
		rancherClient: rancherClient,
	}
}

func (c *ServicesPortRangeClient) Create(container *ServicesPortRange) (*ServicesPortRange, error) {
	resp := &ServicesPortRange{}
	err := c.rancherClient.doCreate(SERVICES_PORT_RANGE_TYPE, container, resp)
	return resp, err
}

func (c *ServicesPortRangeClient) Update(existing *ServicesPortRange, updates interface{}) (*ServicesPortRange, error) {
	resp := &ServicesPortRange{}
	err := c.rancherClient.doUpdate(SERVICES_PORT_RANGE_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ServicesPortRangeClient) List(opts *ListOpts) (*ServicesPortRangeCollection, error) {
	resp := &ServicesPortRangeCollection{}
	err := c.rancherClient.doList(SERVICES_PORT_RANGE_TYPE, opts, resp)
	return resp, err
}

func (c *ServicesPortRangeClient) ById(id string) (*ServicesPortRange, error) {
	resp := &ServicesPortRange{}
	err := c.rancherClient.doById(SERVICES_PORT_RANGE_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ServicesPortRangeClient) Delete(container *ServicesPortRange) error {
	return c.rancherClient.doResourceDelete(SERVICES_PORT_RANGE_TYPE, &container.Resource)
}
