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
// Copyright 2017 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapter

import (
	"github.com/coreos/etcd/etcdserver/api/v3lock/v3lockpb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ls2lsc struct{ ls v3lockpb.LockServer }

func LockServerToLockClient(ls v3lockpb.LockServer) v3lockpb.LockClient {
	return &ls2lsc{ls}
}

func (s *ls2lsc) Lock(ctx context.Context, r *v3lockpb.LockRequest, opts ...grpc.CallOption) (*v3lockpb.LockResponse, error) {
	return s.ls.Lock(ctx, r)
}

func (s *ls2lsc) Unlock(ctx context.Context, r *v3lockpb.UnlockRequest, opts ...grpc.CallOption) (*v3lockpb.UnlockResponse, error) {
	return s.ls.Unlock(ctx, r)
}
