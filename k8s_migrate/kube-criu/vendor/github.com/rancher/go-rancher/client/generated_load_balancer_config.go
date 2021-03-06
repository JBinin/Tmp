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
	LOAD_BALANCER_CONFIG_TYPE = "loadBalancerConfig"
)

type LoadBalancerConfig struct {
	Resource

	HaproxyConfig *HaproxyConfig `json:"haproxyConfig,omitempty" yaml:"haproxy_config,omitempty"`

	LbCookieStickinessPolicy *LoadBalancerCookieStickinessPolicy `json:"lbCookieStickinessPolicy,omitempty" yaml:"lb_cookie_stickiness_policy,omitempty"`
}

type LoadBalancerConfigCollection struct {
	Collection
	Data []LoadBalancerConfig `json:"data,omitempty"`
}

type LoadBalancerConfigClient struct {
	rancherClient *RancherClient
}

type LoadBalancerConfigOperations interface {
	List(opts *ListOpts) (*LoadBalancerConfigCollection, error)
	Create(opts *LoadBalancerConfig) (*LoadBalancerConfig, error)
	Update(existing *LoadBalancerConfig, updates interface{}) (*LoadBalancerConfig, error)
	ById(id string) (*LoadBalancerConfig, error)
	Delete(container *LoadBalancerConfig) error
}

func newLoadBalancerConfigClient(rancherClient *RancherClient) *LoadBalancerConfigClient {
	return &LoadBalancerConfigClient{
		rancherClient: rancherClient,
	}
}

func (c *LoadBalancerConfigClient) Create(container *LoadBalancerConfig) (*LoadBalancerConfig, error) {
	resp := &LoadBalancerConfig{}
	err := c.rancherClient.doCreate(LOAD_BALANCER_CONFIG_TYPE, container, resp)
	return resp, err
}

func (c *LoadBalancerConfigClient) Update(existing *LoadBalancerConfig, updates interface{}) (*LoadBalancerConfig, error) {
	resp := &LoadBalancerConfig{}
	err := c.rancherClient.doUpdate(LOAD_BALANCER_CONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *LoadBalancerConfigClient) List(opts *ListOpts) (*LoadBalancerConfigCollection, error) {
	resp := &LoadBalancerConfigCollection{}
	err := c.rancherClient.doList(LOAD_BALANCER_CONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *LoadBalancerConfigClient) ById(id string) (*LoadBalancerConfig, error) {
	resp := &LoadBalancerConfig{}
	err := c.rancherClient.doById(LOAD_BALANCER_CONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *LoadBalancerConfigClient) Delete(container *LoadBalancerConfig) error {
	return c.rancherClient.doResourceDelete(LOAD_BALANCER_CONFIG_TYPE, &container.Resource)
}
