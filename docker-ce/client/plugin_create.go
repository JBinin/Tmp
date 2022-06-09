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
	"io"
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// PluginCreate creates a plugin
func (cli *Client) PluginCreate(ctx context.Context, createContext io.Reader, createOptions types.PluginCreateOptions) error {
	headers := http.Header(make(map[string][]string))
	headers.Set("Content-Type", "application/x-tar")

	query := url.Values{}
	query.Set("name", createOptions.RepoName)

	resp, err := cli.postRaw(ctx, "/plugins/create", query, createContext, headers)
	if err != nil {
		return err
	}
	ensureReaderClosed(resp)
	return err
}
