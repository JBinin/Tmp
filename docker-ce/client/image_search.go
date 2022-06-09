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

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"golang.org/x/net/context"
)

// ImageSearch makes the docker host to search by a term in a remote registry.
// The list of results is not sorted in any fashion.
func (cli *Client) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	var results []registry.SearchResult
	query := url.Values{}
	query.Set("term", term)
	query.Set("limit", fmt.Sprintf("%d", options.Limit))

	if options.Filters.Len() > 0 {
		filterJSON, err := filters.ToParam(options.Filters)
		if err != nil {
			return results, err
		}
		query.Set("filters", filterJSON)
	}

	resp, err := cli.tryImageSearch(ctx, query, options.RegistryAuth)
	if resp.statusCode == http.StatusUnauthorized && options.PrivilegeFunc != nil {
		newAuthHeader, privilegeErr := options.PrivilegeFunc()
		if privilegeErr != nil {
			return results, privilegeErr
		}
		resp, err = cli.tryImageSearch(ctx, query, newAuthHeader)
	}
	if err != nil {
		return results, err
	}

	err = json.NewDecoder(resp.body).Decode(&results)
	ensureReaderClosed(resp)
	return results, err
}

func (cli *Client) tryImageSearch(ctx context.Context, query url.Values, registryAuth string) (serverResponse, error) {
	headers := map[string][]string{"X-Registry-Auth": {registryAuth}}
	return cli.get(ctx, "/images/search", query, headers)
}
