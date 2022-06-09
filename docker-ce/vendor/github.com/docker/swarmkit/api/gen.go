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
package api

//go:generate protoc -I.:../protobuf:../vendor:../vendor/github.com/gogo/protobuf --gogoswarm_out=plugins=grpc+deepcopy+raftproxy+authenticatedwrapper,import_path=github.com/docker/swarmkit/api,Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,Mplugin/plugin.proto=github.com/docker/swarmkit/protobuf/plugin,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. types.proto specs.proto objects.proto control.proto dispatcher.proto ca.proto snapshot.proto raft.proto health.proto resource.proto logbroker.proto
