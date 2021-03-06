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
	IP_ADDRESS_ASSOCIATE_INPUT_TYPE = "ipAddressAssociateInput"
)

type IpAddressAssociateInput struct {
	Resource

	IpAddressId string `json:"ipAddressId,omitempty" yaml:"ip_address_id,omitempty"`
}

type IpAddressAssociateInputCollection struct {
	Collection
	Data []IpAddressAssociateInput `json:"data,omitempty"`
}

type IpAddressAssociateInputClient struct {
	rancherClient *RancherClient
}

type IpAddressAssociateInputOperations interface {
	List(opts *ListOpts) (*IpAddressAssociateInputCollection, error)
	Create(opts *IpAddressAssociateInput) (*IpAddressAssociateInput, error)
	Update(existing *IpAddressAssociateInput, updates interface{}) (*IpAddressAssociateInput, error)
	ById(id string) (*IpAddressAssociateInput, error)
	Delete(container *IpAddressAssociateInput) error
}

func newIpAddressAssociateInputClient(rancherClient *RancherClient) *IpAddressAssociateInputClient {
	return &IpAddressAssociateInputClient{
		rancherClient: rancherClient,
	}
}

func (c *IpAddressAssociateInputClient) Create(container *IpAddressAssociateInput) (*IpAddressAssociateInput, error) {
	resp := &IpAddressAssociateInput{}
	err := c.rancherClient.doCreate(IP_ADDRESS_ASSOCIATE_INPUT_TYPE, container, resp)
	return resp, err
}

func (c *IpAddressAssociateInputClient) Update(existing *IpAddressAssociateInput, updates interface{}) (*IpAddressAssociateInput, error) {
	resp := &IpAddressAssociateInput{}
	err := c.rancherClient.doUpdate(IP_ADDRESS_ASSOCIATE_INPUT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *IpAddressAssociateInputClient) List(opts *ListOpts) (*IpAddressAssociateInputCollection, error) {
	resp := &IpAddressAssociateInputCollection{}
	err := c.rancherClient.doList(IP_ADDRESS_ASSOCIATE_INPUT_TYPE, opts, resp)
	return resp, err
}

func (c *IpAddressAssociateInputClient) ById(id string) (*IpAddressAssociateInput, error) {
	resp := &IpAddressAssociateInput{}
	err := c.rancherClient.doById(IP_ADDRESS_ASSOCIATE_INPUT_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *IpAddressAssociateInputClient) Delete(container *IpAddressAssociateInput) error {
	return c.rancherClient.doResourceDelete(IP_ADDRESS_ASSOCIATE_INPUT_TYPE, &container.Resource)
}
