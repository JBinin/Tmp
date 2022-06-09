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
// Copyright 2016 The etcd Authors
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
	pb "github.com/coreos/etcd/etcdserver/etcdserverpb"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type kvs2kvc struct{ kvs pb.KVServer }

func KvServerToKvClient(kvs pb.KVServer) pb.KVClient {
	return &kvs2kvc{kvs}
}

func (s *kvs2kvc) Range(ctx context.Context, in *pb.RangeRequest, opts ...grpc.CallOption) (*pb.RangeResponse, error) {
	return s.kvs.Range(ctx, in)
}

func (s *kvs2kvc) Put(ctx context.Context, in *pb.PutRequest, opts ...grpc.CallOption) (*pb.PutResponse, error) {
	return s.kvs.Put(ctx, in)
}

func (s *kvs2kvc) DeleteRange(ctx context.Context, in *pb.DeleteRangeRequest, opts ...grpc.CallOption) (*pb.DeleteRangeResponse, error) {
	return s.kvs.DeleteRange(ctx, in)
}

func (s *kvs2kvc) Txn(ctx context.Context, in *pb.TxnRequest, opts ...grpc.CallOption) (*pb.TxnResponse, error) {
	return s.kvs.Txn(ctx, in)
}

func (s *kvs2kvc) Compact(ctx context.Context, in *pb.CompactionRequest, opts ...grpc.CallOption) (*pb.CompactionResponse, error) {
	return s.kvs.Compact(ctx, in)
}
