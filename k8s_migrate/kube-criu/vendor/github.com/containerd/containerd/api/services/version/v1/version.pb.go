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
// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/containerd/containerd/api/services/version/v1/version.proto

/*
	Package version is a generated protocol buffer package.

	It is generated from these files:
		github.com/containerd/containerd/api/services/version/v1/version.proto

	It has these top-level messages:
		VersionResponse
*/
package version

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/types"

// skipping weak import gogoproto "github.com/gogo/protobuf/gogoproto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type VersionResponse struct {
	Version  string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Revision string `protobuf:"bytes,2,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (m *VersionResponse) Reset()                    { *m = VersionResponse{} }
func (*VersionResponse) ProtoMessage()               {}
func (*VersionResponse) Descriptor() ([]byte, []int) { return fileDescriptorVersion, []int{0} }

func init() {
	proto.RegisterType((*VersionResponse)(nil), "containerd.services.version.v1.VersionResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Version service

type VersionClient interface {
	Version(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
}

type versionClient struct {
	cc *grpc.ClientConn
}

func NewVersionClient(cc *grpc.ClientConn) VersionClient {
	return &versionClient{cc}
}

func (c *versionClient) Version(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := grpc.Invoke(ctx, "/containerd.services.version.v1.Version/Version", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Version service

type VersionServer interface {
	Version(context.Context, *google_protobuf.Empty) (*VersionResponse, error)
}

func RegisterVersionServer(s *grpc.Server, srv VersionServer) {
	s.RegisterService(&_Version_serviceDesc, srv)
}

func _Version_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.version.v1.Version/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionServer).Version(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Version_serviceDesc = grpc.ServiceDesc{
	ServiceName: "containerd.services.version.v1.Version",
	HandlerType: (*VersionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Version",
			Handler:    _Version_Version_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/containerd/containerd/api/services/version/v1/version.proto",
}

func (m *VersionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VersionResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Version) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintVersion(dAtA, i, uint64(len(m.Version)))
		i += copy(dAtA[i:], m.Version)
	}
	if len(m.Revision) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintVersion(dAtA, i, uint64(len(m.Revision)))
		i += copy(dAtA[i:], m.Revision)
	}
	return i, nil
}

func encodeVarintVersion(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *VersionResponse) Size() (n int) {
	var l int
	_ = l
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovVersion(uint64(l))
	}
	l = len(m.Revision)
	if l > 0 {
		n += 1 + l + sovVersion(uint64(l))
	}
	return n
}

func sovVersion(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozVersion(x uint64) (n int) {
	return sovVersion(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *VersionResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&VersionResponse{`,
		`Version:` + fmt.Sprintf("%v", this.Version) + `,`,
		`Revision:` + fmt.Sprintf("%v", this.Revision) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringVersion(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *VersionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVersion
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: VersionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VersionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVersion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthVersion
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Revision", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVersion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthVersion
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Revision = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVersion(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthVersion
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipVersion(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVersion
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowVersion
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowVersion
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthVersion
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowVersion
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipVersion(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthVersion = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVersion   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/containerd/containerd/api/services/version/v1/version.proto", fileDescriptorVersion)
}

var fileDescriptorVersion = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0x4b, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0xcf, 0x2b, 0x49, 0xcc, 0xcc, 0x4b, 0x2d,
	0x4a, 0x41, 0x66, 0x26, 0x16, 0x64, 0xea, 0x17, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x16, 0xeb,
	0x97, 0xa5, 0x16, 0x15, 0x67, 0xe6, 0xe7, 0xe9, 0x97, 0x19, 0xc2, 0x98, 0x7a, 0x05, 0x45, 0xf9,
	0x25, 0xf9, 0x42, 0x72, 0x08, 0x1d, 0x7a, 0x30, 0xd5, 0x7a, 0x30, 0x25, 0x65, 0x86, 0x52, 0xd2,
	0xe9, 0xf9, 0xf9, 0xe9, 0x39, 0xa9, 0xfa, 0x60, 0xd5, 0x49, 0xa5, 0x69, 0xfa, 0xa9, 0xb9, 0x05,
	0x25, 0x95, 0x10, 0xcd, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0xa6, 0x3e, 0x88, 0x05, 0x11,
	0x55, 0x72, 0xe7, 0xe2, 0x0f, 0x83, 0x18, 0x10, 0x94, 0x5a, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0x2a,
	0x24, 0xc1, 0xc5, 0x0e, 0x35, 0x53, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc6, 0x15, 0x92,
	0xe2, 0xe2, 0x28, 0x4a, 0x2d, 0xcb, 0x04, 0x4b, 0x31, 0x81, 0xa5, 0xe0, 0x7c, 0xa3, 0x58, 0x2e,
	0x76, 0xa8, 0x41, 0x42, 0x41, 0x08, 0xa6, 0x98, 0x1e, 0xc4, 0x49, 0x7a, 0x30, 0x27, 0xe9, 0xb9,
	0x82, 0x9c, 0x24, 0xa5, 0xaf, 0x87, 0xdf, 0x2b, 0x7a, 0x68, 0x8e, 0x72, 0x8a, 0x3a, 0xf1, 0x50,
	0x8e, 0xe1, 0xc6, 0x43, 0x39, 0x86, 0x86, 0x47, 0x72, 0x8c, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78,
	0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0x63, 0x94, 0x03, 0xb9, 0x81, 0x6b, 0x0d, 0x65, 0x46, 0x30,
	0x26, 0xb1, 0x81, 0x9d, 0x67, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x95, 0x0d, 0x52, 0x23, 0xa9,
	0x01, 0x00, 0x00,
}