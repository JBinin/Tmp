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
package raftselector

import (
	"errors"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// ConnProvider is basic interface for connecting API package(raft proxy in particular)
// to manager/state/raft package without import cycles. It provides only one
// method for obtaining connection to leader.
type ConnProvider interface {
	LeaderConn(ctx context.Context) (*grpc.ClientConn, error)
}

// ErrIsLeader is returned from LeaderConn method when current machine is leader.
// It's just shim between packages to avoid import cycles.
var ErrIsLeader = errors.New("current node is leader")
