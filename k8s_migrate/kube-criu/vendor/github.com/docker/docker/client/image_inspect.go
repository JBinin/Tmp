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

// ImageInspectWithRaw returns the image information and its raw representation.
func (cli *Client) ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error) {
	serverResp, err := cli.get(ctx, "/images/"+imageID+"/json", nil, nil)
	if err != nil {
		if serverResp.statusCode == http.StatusNotFound {
			return types.ImageInspect{}, nil, imageNotFoundError{imageID}
		}
		return types.ImageInspect{}, nil, err
	}
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil {
		return types.ImageInspect{}, nil, err
	}

	var response types.ImageInspect
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
}
