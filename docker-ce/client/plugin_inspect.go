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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// PluginInspectWithRaw inspects an existing plugin
func (cli *Client) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error) {
	resp, err := cli.get(ctx, "/plugins/"+name+"/json", nil, nil)
	if err != nil {
		if resp.statusCode == http.StatusNotFound {
			return nil, nil, pluginNotFoundError{name}
		}
		return nil, nil, err
	}

	defer ensureReaderClosed(resp)
	body, err := ioutil.ReadAll(resp.body)
	if err != nil {
		return nil, nil, err
	}
	var p types.Plugin
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&p)
	return &p, body, err
}
