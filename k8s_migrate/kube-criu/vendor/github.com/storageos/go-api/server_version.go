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
	"net/http"

	"github.com/storageos/go-api/types"
)

// ServerVersion returns the server's version and runtime info.
func (c *Client) ServerVersion(ctx context.Context) (*types.VersionInfo, error) {

	// Send as unversioned
	resp, err := c.do("GET", "version", doOptions{context: ctx, unversioned: true})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, newError(resp)
	}
	defer resp.Body.Close()
	var version types.VersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return nil, err
	}
	return &version, nil
}
