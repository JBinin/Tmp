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
package state

import (
	"github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

// A Proposer can propose actions to a cluster.
type Proposer interface {
	// ProposeValue adds storeAction to the distributed log. If this
	// completes successfully, ProposeValue calls cb to commit the
	// proposed changes. The callback is necessary for the Proposer to make
	// sure that the changes are committed before it interacts further
	// with the store.
	ProposeValue(ctx context.Context, storeAction []*api.StoreAction, cb func()) error
	GetVersion() *api.Version
}
