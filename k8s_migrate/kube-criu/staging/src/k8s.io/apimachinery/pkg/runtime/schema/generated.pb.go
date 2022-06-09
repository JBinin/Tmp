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
Copyright The Kubernetes Authors.

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

// Code generated by protoc-gen-gogo.
// source: k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/runtime/schema/generated.proto
// DO NOT EDIT!

/*
Package schema is a generated protocol buffer package.

It is generated from these files:
	k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/runtime/schema/generated.proto

It has these top-level messages:
*/
package schema

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func init() {
	proto.RegisterFile("k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/runtime/schema/generated.proto", fileDescriptorGenerated)
}

var fileDescriptorGenerated = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xcc, 0xaf, 0x6e, 0xc3, 0x30,
	0x10, 0xc7, 0x71, 0x9b, 0x0c, 0x0c, 0x0e, 0x0e, 0x1c, 0x1c, 0xda, 0x7c, 0x74, 0xb8, 0x2f, 0x50,
	0x5e, 0xe6, 0x24, 0x57, 0xc7, 0xb2, 0xfc, 0x47, 0x8e, 0x5d, 0xa9, 0xac, 0x8f, 0xd0, 0xc7, 0x0a,
	0x0c, 0x0c, 0x6c, 0xdc, 0x17, 0xa9, 0x64, 0x07, 0x94, 0xdd, 0x4f, 0xa7, 0xcf, 0xf7, 0xf3, 0x68,
	0xfe, 0x27, 0xa1, 0x3d, 0x9a, 0xdc, 0x51, 0x74, 0x94, 0x68, 0xc2, 0x0b, 0xb9, 0xc1, 0x47, 0xdc,
	0x1f, 0x32, 0x68, 0x2b, 0xfb, 0x51, 0x3b, 0x8a, 0x57, 0x0c, 0x46, 0x61, 0xcc, 0x2e, 0x69, 0x4b,
	0x38, 0xf5, 0x23, 0x59, 0x89, 0x8a, 0x1c, 0x45, 0x99, 0x68, 0x10, 0x21, 0xfa, 0xe4, 0xbf, 0x7e,
	0x9a, 0x13, 0xef, 0x4e, 0x04, 0xa3, 0xc4, 0xee, 0x44, 0x73, 0xdf, 0x7f, 0x4a, 0xa7, 0x31, 0x77,
	0xa2, 0xf7, 0x16, 0x95, 0x57, 0x1e, 0x2b, 0xef, 0xf2, 0xb9, 0xae, 0x3a, 0xea, 0xd5, 0xb2, 0x87,
	0xdf, 0x79, 0x03, 0xb6, 0x6c, 0xc0, 0xd6, 0x0d, 0xd8, 0xad, 0x00, 0x9f, 0x0b, 0xf0, 0xa5, 0x00,
	0x5f, 0x0b, 0xf0, 0x47, 0x01, 0x7e, 0x7f, 0x02, 0x3b, 0x7d, 0xb4, 0xf8, 0x2b, 0x00, 0x00, 0xff,
	0xff, 0xba, 0x7e, 0x65, 0xf4, 0xd6, 0x00, 0x00, 0x00,
}
