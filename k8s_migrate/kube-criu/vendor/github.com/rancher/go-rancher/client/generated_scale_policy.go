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
	SCALE_POLICY_TYPE = "scalePolicy"
)

type ScalePolicy struct {
	Resource

	Increment int64 `json:"increment,omitempty" yaml:"increment,omitempty"`

	Max int64 `json:"max,omitempty" yaml:"max,omitempty"`

	Min int64 `json:"min,omitempty" yaml:"min,omitempty"`
}

type ScalePolicyCollection struct {
	Collection
	Data []ScalePolicy `json:"data,omitempty"`
}

type ScalePolicyClient struct {
	rancherClient *RancherClient
}

type ScalePolicyOperations interface {
	List(opts *ListOpts) (*ScalePolicyCollection, error)
	Create(opts *ScalePolicy) (*ScalePolicy, error)
	Update(existing *ScalePolicy, updates interface{}) (*ScalePolicy, error)
	ById(id string) (*ScalePolicy, error)
	Delete(container *ScalePolicy) error
}

func newScalePolicyClient(rancherClient *RancherClient) *ScalePolicyClient {
	return &ScalePolicyClient{
		rancherClient: rancherClient,
	}
}

func (c *ScalePolicyClient) Create(container *ScalePolicy) (*ScalePolicy, error) {
	resp := &ScalePolicy{}
	err := c.rancherClient.doCreate(SCALE_POLICY_TYPE, container, resp)
	return resp, err
}

func (c *ScalePolicyClient) Update(existing *ScalePolicy, updates interface{}) (*ScalePolicy, error) {
	resp := &ScalePolicy{}
	err := c.rancherClient.doUpdate(SCALE_POLICY_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ScalePolicyClient) List(opts *ListOpts) (*ScalePolicyCollection, error) {
	resp := &ScalePolicyCollection{}
	err := c.rancherClient.doList(SCALE_POLICY_TYPE, opts, resp)
	return resp, err
}

func (c *ScalePolicyClient) ById(id string) (*ScalePolicy, error) {
	resp := &ScalePolicy{}
	err := c.rancherClient.doById(SCALE_POLICY_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ScalePolicyClient) Delete(container *ScalePolicy) error {
	return c.rancherClient.doResourceDelete(SCALE_POLICY_TYPE, &container.Resource)
}
