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
 * Copyright 2016 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package internal contains gRPC-internal code for testing, to avoid polluting
// the godoc of the top-level grpc package.
package internal

// TestingCloseConns closes all existing transports but keeps
// grpcServer.lis accepting new connections.
//
// The provided grpcServer must be of type *grpc.Server. It is untyped
// for circular dependency reasons.
var TestingCloseConns func(grpcServer interface{})

// TestingUseHandlerImpl enables the http.Handler-based server implementation.
// It must be called before Serve and requires TLS credentials.
//
// The provided grpcServer must be of type *grpc.Server. It is untyped
// for circular dependency reasons.
var TestingUseHandlerImpl func(grpcServer interface{})
