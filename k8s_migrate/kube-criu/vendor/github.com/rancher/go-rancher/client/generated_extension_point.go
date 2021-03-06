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
	EXTENSION_POINT_TYPE = "extensionPoint"
)

type ExtensionPoint struct {
	Resource

	ExcludeSetting string `json:"excludeSetting,omitempty" yaml:"exclude_setting,omitempty"`

	Implementations []interface{} `json:"implementations,omitempty" yaml:"implementations,omitempty"`

	IncludeSetting string `json:"includeSetting,omitempty" yaml:"include_setting,omitempty"`

	ListSetting string `json:"listSetting,omitempty" yaml:"list_setting,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ExtensionPointCollection struct {
	Collection
	Data []ExtensionPoint `json:"data,omitempty"`
}

type ExtensionPointClient struct {
	rancherClient *RancherClient
}

type ExtensionPointOperations interface {
	List(opts *ListOpts) (*ExtensionPointCollection, error)
	Create(opts *ExtensionPoint) (*ExtensionPoint, error)
	Update(existing *ExtensionPoint, updates interface{}) (*ExtensionPoint, error)
	ById(id string) (*ExtensionPoint, error)
	Delete(container *ExtensionPoint) error
}

func newExtensionPointClient(rancherClient *RancherClient) *ExtensionPointClient {
	return &ExtensionPointClient{
		rancherClient: rancherClient,
	}
}

func (c *ExtensionPointClient) Create(container *ExtensionPoint) (*ExtensionPoint, error) {
	resp := &ExtensionPoint{}
	err := c.rancherClient.doCreate(EXTENSION_POINT_TYPE, container, resp)
	return resp, err
}

func (c *ExtensionPointClient) Update(existing *ExtensionPoint, updates interface{}) (*ExtensionPoint, error) {
	resp := &ExtensionPoint{}
	err := c.rancherClient.doUpdate(EXTENSION_POINT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ExtensionPointClient) List(opts *ListOpts) (*ExtensionPointCollection, error) {
	resp := &ExtensionPointCollection{}
	err := c.rancherClient.doList(EXTENSION_POINT_TYPE, opts, resp)
	return resp, err
}

func (c *ExtensionPointClient) ById(id string) (*ExtensionPoint, error) {
	resp := &ExtensionPoint{}
	err := c.rancherClient.doById(EXTENSION_POINT_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *ExtensionPointClient) Delete(container *ExtensionPoint) error {
	return c.rancherClient.doResourceDelete(EXTENSION_POINT_TYPE, &container.Resource)
}
