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
package storageos

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/storageos/go-api/types"
)

var (
	// LoggerAPIPrefix is a partial path to the HTTP endpoint.
	LoggerAPIPrefix = "logs"
)

// LoggerConfig returns every cluster node's logging configuration.
func (c *Client) LoggerConfig(opts types.ListOptions) ([]*types.Logger, error) {

	listOpts := doOptions{
		fieldSelector: opts.FieldSelector,
		labelSelector: opts.LabelSelector,
		context:       opts.Context,
	}

	if opts.LabelSelector != "" {
		query := url.Values{}
		query.Add("labelSelector", opts.LabelSelector)
		listOpts.values = query
	}

	resp, err := c.do("GET", LoggerAPIPrefix+"/cluster/config", listOpts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var loggers []*types.Logger
	if err := json.NewDecoder(resp.Body).Decode(&loggers); err != nil {
		return nil, err
	}
	return loggers, nil

}

// LoggerUpdate patches updates to logging configuration.  Fields to update must
// be listed in the Fields value, and if a list of Nodes is given it will only
// apply to the nodes listed.  Returns the updated configuration.
func (c *Client) LoggerUpdate(opts types.LoggerUpdateOptions) ([]*types.Logger, error) {

	resp, err := c.do("PATCH", LoggerAPIPrefix+"/cluster/config", doOptions{
		data:    opts,
		context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var loggers []*types.Logger
	if err := json.NewDecoder(resp.Body).Decode(&loggers); err != nil {
		return nil, err
	}
	return loggers, nil
}
