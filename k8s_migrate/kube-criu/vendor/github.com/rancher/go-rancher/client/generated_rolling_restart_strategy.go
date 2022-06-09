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
	ROLLING_RESTART_STRATEGY_TYPE = "rollingRestartStrategy"
)

type RollingRestartStrategy struct {
	Resource

	BatchSize int64 `json:"batchSize,omitempty" yaml:"batch_size,omitempty"`

	IntervalMillis int64 `json:"intervalMillis,omitempty" yaml:"interval_millis,omitempty"`
}

type RollingRestartStrategyCollection struct {
	Collection
	Data []RollingRestartStrategy `json:"data,omitempty"`
}

type RollingRestartStrategyClient struct {
	rancherClient *RancherClient
}

type RollingRestartStrategyOperations interface {
	List(opts *ListOpts) (*RollingRestartStrategyCollection, error)
	Create(opts *RollingRestartStrategy) (*RollingRestartStrategy, error)
	Update(existing *RollingRestartStrategy, updates interface{}) (*RollingRestartStrategy, error)
	ById(id string) (*RollingRestartStrategy, error)
	Delete(container *RollingRestartStrategy) error
}

func newRollingRestartStrategyClient(rancherClient *RancherClient) *RollingRestartStrategyClient {
	return &RollingRestartStrategyClient{
		rancherClient: rancherClient,
	}
}

func (c *RollingRestartStrategyClient) Create(container *RollingRestartStrategy) (*RollingRestartStrategy, error) {
	resp := &RollingRestartStrategy{}
	err := c.rancherClient.doCreate(ROLLING_RESTART_STRATEGY_TYPE, container, resp)
	return resp, err
}

func (c *RollingRestartStrategyClient) Update(existing *RollingRestartStrategy, updates interface{}) (*RollingRestartStrategy, error) {
	resp := &RollingRestartStrategy{}
	err := c.rancherClient.doUpdate(ROLLING_RESTART_STRATEGY_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *RollingRestartStrategyClient) List(opts *ListOpts) (*RollingRestartStrategyCollection, error) {
	resp := &RollingRestartStrategyCollection{}
	err := c.rancherClient.doList(ROLLING_RESTART_STRATEGY_TYPE, opts, resp)
	return resp, err
}

func (c *RollingRestartStrategyClient) ById(id string) (*RollingRestartStrategy, error) {
	resp := &RollingRestartStrategy{}
	err := c.rancherClient.doById(ROLLING_RESTART_STRATEGY_TYPE, id, resp)
	if apiError, ok := err.(*ApiError); ok {
		if apiError.StatusCode == 404 {
			return nil, nil
		}
	}
	return resp, err
}

func (c *RollingRestartStrategyClient) Delete(container *RollingRestartStrategy) error {
	return c.rancherClient.doResourceDelete(ROLLING_RESTART_STRATEGY_TYPE, &container.Resource)
}
