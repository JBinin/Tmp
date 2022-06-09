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
// source: k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/api/resource/generated.proto
// DO NOT EDIT!

/*
Package resource is a generated protocol buffer package.

It is generated from these files:
	k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/api/resource/generated.proto

It has these top-level messages:
	Quantity
*/
package resource

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

func (m *Quantity) Reset()                    { *m = Quantity{} }
func (*Quantity) ProtoMessage()               {}
func (*Quantity) Descriptor() ([]byte, []int) { return fileDescriptorGenerated, []int{0} }

func init() {
	proto.RegisterType((*Quantity)(nil), "k8s.io.apimachinery.pkg.api.resource.Quantity")
}

func init() {
	proto.RegisterFile("k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/api/resource/generated.proto", fileDescriptorGenerated)
}

var fileDescriptorGenerated = []byte{
	// 237 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8e, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x40, 0xcf, 0x0b, 0x2a, 0x19, 0x2b, 0x84, 0x10, 0xc3, 0xa5, 0x42, 0x0c, 0x2c, 0xd8, 0x6b,
	0xc5, 0xc8, 0xce, 0x00, 0x23, 0x5b, 0x92, 0x1e, 0xae, 0x15, 0xd5, 0x8e, 0x2e, 0x36, 0x52, 0xb7,
	0x8e, 0x8c, 0x1d, 0x19, 0x9b, 0xbf, 0xe9, 0xd8, 0xb1, 0x03, 0x03, 0x31, 0x3f, 0x82, 0xea, 0x36,
	0x52, 0xb7, 0x7b, 0xef, 0xf4, 0x4e, 0x97, 0xbd, 0xd4, 0xd3, 0x56, 0x1a, 0xa7, 0xea, 0x50, 0x12,
	0x5b, 0xf2, 0xd4, 0xaa, 0x4f, 0xb2, 0x33, 0xc7, 0xea, 0xb4, 0x28, 0x1a, 0xb3, 0x28, 0xaa, 0xb9,
	0xb1, 0xc4, 0x4b, 0xd5, 0xd4, 0xfa, 0x20, 0x14, 0x53, 0xeb, 0x02, 0x57, 0xa4, 0x34, 0x59, 0xe2,
	0xc2, 0xd3, 0x4c, 0x36, 0xec, 0xbc, 0x1b, 0xdf, 0x1f, 0x2b, 0x79, 0x5e, 0xc9, 0xa6, 0xd6, 0x07,
	0x21, 0x87, 0xea, 0xf6, 0x51, 0x1b, 0x3f, 0x0f, 0xa5, 0xac, 0xdc, 0x42, 0x69, 0xa7, 0x9d, 0x4a,
	0x71, 0x19, 0x3e, 0x12, 0x25, 0x48, 0xd3, 0xf1, 0xe8, 0xdd, 0x34, 0x1b, 0xbd, 0x86, 0xc2, 0x7a,
	0xe3, 0x97, 0xe3, 0xeb, 0xec, 0xa2, 0xf5, 0x6c, 0xac, 0xbe, 0x11, 0x13, 0xf1, 0x70, 0xf9, 0x76,
	0xa2, 0xa7, 0xab, 0xef, 0x4d, 0x0e, 0x5f, 0x5d, 0x0e, 0xeb, 0x2e, 0x87, 0x4d, 0x97, 0xc3, 0xea,
	0x67, 0x02, 0xcf, 0x72, 0xdb, 0x23, 0xec, 0x7a, 0x84, 0x7d, 0x8f, 0xb0, 0x8a, 0x28, 0xb6, 0x11,
	0xc5, 0x2e, 0xa2, 0xd8, 0x47, 0x14, 0xbf, 0x11, 0xc5, 0xfa, 0x0f, 0xe1, 0x7d, 0x34, 0x3c, 0xf6,
	0x1f, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x08, 0x88, 0x49, 0x0e, 0x01, 0x00, 0x00,
}
