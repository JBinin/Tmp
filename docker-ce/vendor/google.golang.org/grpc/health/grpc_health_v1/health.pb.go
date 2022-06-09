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
// Code generated by protoc-gen-go.
// source: health.proto
// DO NOT EDIT!

/*
Package grpc_health_v1 is a generated protocol buffer package.

It is generated from these files:
	health.proto

It has these top-level messages:
	HealthCheckRequest
	HealthCheckResponse
*/
package grpc_health_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type HealthCheckResponse_ServingStatus int32

const (
	HealthCheckResponse_UNKNOWN     HealthCheckResponse_ServingStatus = 0
	HealthCheckResponse_SERVING     HealthCheckResponse_ServingStatus = 1
	HealthCheckResponse_NOT_SERVING HealthCheckResponse_ServingStatus = 2
)

var HealthCheckResponse_ServingStatus_name = map[int32]string{
	0: "UNKNOWN",
	1: "SERVING",
	2: "NOT_SERVING",
}
var HealthCheckResponse_ServingStatus_value = map[string]int32{
	"UNKNOWN":     0,
	"SERVING":     1,
	"NOT_SERVING": 2,
}

func (x HealthCheckResponse_ServingStatus) String() string {
	return proto.EnumName(HealthCheckResponse_ServingStatus_name, int32(x))
}
func (HealthCheckResponse_ServingStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1, 0}
}

type HealthCheckRequest struct {
	Service string `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *HealthCheckRequest) Reset()                    { *m = HealthCheckRequest{} }
func (m *HealthCheckRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthCheckRequest) ProtoMessage()               {}
func (*HealthCheckRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type HealthCheckResponse struct {
	Status HealthCheckResponse_ServingStatus `protobuf:"varint,1,opt,name=status,enum=grpc.health.v1.HealthCheckResponse_ServingStatus" json:"status,omitempty"`
}

func (m *HealthCheckResponse) Reset()                    { *m = HealthCheckResponse{} }
func (m *HealthCheckResponse) String() string            { return proto.CompactTextString(m) }
func (*HealthCheckResponse) ProtoMessage()               {}
func (*HealthCheckResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*HealthCheckRequest)(nil), "grpc.health.v1.HealthCheckRequest")
	proto.RegisterType((*HealthCheckResponse)(nil), "grpc.health.v1.HealthCheckResponse")
	proto.RegisterEnum("grpc.health.v1.HealthCheckResponse_ServingStatus", HealthCheckResponse_ServingStatus_name, HealthCheckResponse_ServingStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Health service

type HealthClient interface {
	Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type healthClient struct {
	cc *grpc.ClientConn
}

func NewHealthClient(cc *grpc.ClientConn) HealthClient {
	return &healthClient{cc}
}

func (c *healthClient) Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := grpc.Invoke(ctx, "/grpc.health.v1.Health/Check", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Health service

type HealthServer interface {
	Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

func RegisterHealthServer(s *grpc.Server, srv HealthServer) {
	s.RegisterService(&_Health_serviceDesc, srv)
}

func _Health_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.health.v1.Health/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthServer).Check(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Health_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.health.v1.Health",
	HandlerType: (*HealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _Health_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "health.proto",
}

func init() { proto.RegisterFile("health.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xc9, 0x48, 0x4d, 0xcc,
	0x29, 0xc9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4b, 0x2f, 0x2a, 0x48, 0xd6, 0x83,
	0x0a, 0x95, 0x19, 0x2a, 0xe9, 0x71, 0x09, 0x79, 0x80, 0x39, 0xce, 0x19, 0xa9, 0xc9, 0xd9, 0x41,
	0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x12, 0x5c, 0xec, 0xc5, 0xa9, 0x45, 0x65, 0x99, 0xc9,
	0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x30, 0xae, 0xd2, 0x1c, 0x46, 0x2e, 0x61, 0x14,
	0x0d, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x9e, 0x5c, 0x6c, 0xc5, 0x25, 0x89, 0x25, 0xa5,
	0xc5, 0x60, 0x0d, 0x7c, 0x46, 0x86, 0x7a, 0xa8, 0x16, 0xe9, 0x61, 0xd1, 0xa4, 0x17, 0x0c, 0x32,
	0x34, 0x2f, 0x3d, 0x18, 0xac, 0x31, 0x08, 0x6a, 0x80, 0x92, 0x15, 0x17, 0x2f, 0x8a, 0x84, 0x10,
	0x37, 0x17, 0x7b, 0xa8, 0x9f, 0xb7, 0x9f, 0x7f, 0xb8, 0x9f, 0x00, 0x03, 0x88, 0x13, 0xec, 0x1a,
	0x14, 0xe6, 0xe9, 0xe7, 0x2e, 0xc0, 0x28, 0xc4, 0xcf, 0xc5, 0xed, 0xe7, 0x1f, 0x12, 0x0f, 0x13,
	0x60, 0x32, 0x8a, 0xe2, 0x62, 0x83, 0x58, 0x24, 0x14, 0xc0, 0xc5, 0x0a, 0xb6, 0x4c, 0x48, 0x09,
	0xaf, 0x4b, 0xc0, 0xfe, 0x95, 0x52, 0x26, 0xc2, 0xb5, 0x49, 0x6c, 0xe0, 0x10, 0x34, 0x06, 0x04,
	0x00, 0x00, 0xff, 0xff, 0xac, 0x56, 0x2a, 0xcb, 0x51, 0x01, 0x00, 0x00,
}
