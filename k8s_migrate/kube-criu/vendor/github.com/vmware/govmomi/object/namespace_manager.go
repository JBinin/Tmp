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
/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package object

import (
	"context"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

type DatastoreNamespaceManager struct {
	Common
}

func NewDatastoreNamespaceManager(c *vim25.Client) *DatastoreNamespaceManager {
	n := DatastoreNamespaceManager{
		Common: NewCommon(c, *c.ServiceContent.DatastoreNamespaceManager),
	}

	return &n
}

// CreateDirectory creates a top-level directory on the given vsan datastore, using
// the given user display name hint and opaque storage policy.
func (nm DatastoreNamespaceManager) CreateDirectory(ctx context.Context, ds *Datastore, displayName string, policy string) (string, error) {

	req := &types.CreateDirectory{
		This:        nm.Reference(),
		Datastore:   ds.Reference(),
		DisplayName: displayName,
		Policy:      policy,
	}

	resp, err := methods.CreateDirectory(ctx, nm.c, req)
	if err != nil {
		return "", err
	}

	return resp.Returnval, nil
}

// DeleteDirectory deletes the given top-level directory from a vsan datastore.
func (nm DatastoreNamespaceManager) DeleteDirectory(ctx context.Context, dc *Datacenter, datastorePath string) error {

	req := &types.DeleteDirectory{
		This:          nm.Reference(),
		DatastorePath: datastorePath,
	}

	if dc != nil {
		ref := dc.Reference()
		req.Datacenter = &ref
	}

	if _, err := methods.DeleteDirectory(ctx, nm.c, req); err != nil {
		return err
	}

	return nil
}
