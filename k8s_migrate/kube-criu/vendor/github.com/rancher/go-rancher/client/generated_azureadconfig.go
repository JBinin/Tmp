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
	AZUREADCONFIG_TYPE = "azureadconfig"
)

type Azureadconfig struct {
	Resource

	AccessMode string `json:"accessMode,omitempty" yaml:"access_mode,omitempty"`

	AdminAccountPassword string `json:"adminAccountPassword,omitempty" yaml:"admin_account_password,omitempty"`

	AdminAccountUsername string `json:"adminAccountUsername,omitempty" yaml:"admin_account_username,omitempty"`

	ClientId string `json:"clientId,omitempty" yaml:"client_id,omitempty"`

	Domain string `json:"domain,omitempty" yaml:"domain,omitempty"`

	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	TenantId string `json:"tenantId,omitempty" yaml:"tenant_id,omitempty"`
}

type AzureadconfigCollection struct {
	Collection
	Data []Azureadconfig `json:"data,omitempty"`
}

type AzureadconfigClient struct {
	rancherClient *RancherClient
}

type AzureadconfigOperations interface {
	List(opts *ListOpts) (*AzureadconfigCollection, error)
	Create(opts *Azureadconfig) (*Azureadconfig, error)
	Update(existing *Azureadconfig, updates interface{}) (*Azureadconfig, error)
	ById(id string) (*Azureadconfig, error)
	Delete(container *Azureadconfig) error
}

func newAzureadconfigClient(rancherClient *RancherClient) *AzureadconfigClient {
	return &AzureadconfigClient{
		rancherClient: rancherClient,
	}
}

func (c *AzureadconfigClient) Create(container *Azureadconfig) (*Azureadconfig, error) {
	resp := &Azureadconfig{}
	err := c.rancherClient.doCreate(AZUREADCONFIG_TYPE, container, resp)
	return resp, err
}

func (c *AzureadconfigClient) Update(existing *Azureadconfig, updates interface{}) (*Azureadconfig, error) {
	resp := &Azureadconfig{}
	err := c.rancherClient.doUpdate(AZUREADCONFIG_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *AzureadconfigClient) List(opts *ListOpts) (*AzureadconfigCollection, error) {
	resp := &AzureadconfigCollection{}
	err := c.rancherClient.doList(AZUREADCONFIG_TYPE, opts, resp)
	return resp, err
}

func (c *AzureadconfigClient) ById(id string) (*Azureadconfig, error) {
	resp := &Azureadconfig{}
	err := c.rancherClient.doById(AZUREADCONFIG_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *AzureadconfigClient) Delete(container *Azureadconfig) error {
	return c.rancherClient.doResourceDelete(AZUREADCONFIG_TYPE, &container.Resource)
}
