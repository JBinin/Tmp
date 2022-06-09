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
// source: extension.proto
// DO NOT EDIT!

/*
Package openapiextension_v1 is a generated protocol buffer package.

It is generated from these files:
	extension.proto

It has these top-level messages:
	Version
	ExtensionHandlerRequest
	ExtensionHandlerResponse
	Wrapper
*/
package openapiextension_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The version number of OpenAPI compiler.
type Version struct {
	Major int32 `protobuf:"varint,1,opt,name=major" json:"major,omitempty"`
	Minor int32 `protobuf:"varint,2,opt,name=minor" json:"minor,omitempty"`
	Patch int32 `protobuf:"varint,3,opt,name=patch" json:"patch,omitempty"`
	// A suffix for alpha, beta or rc release, e.g., "alpha-1", "rc2". It should
	// be empty for mainline stable releases.
	Suffix string `protobuf:"bytes,4,opt,name=suffix" json:"suffix,omitempty"`
}

func (m *Version) Reset()                    { *m = Version{} }
func (m *Version) String() string            { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()               {}
func (*Version) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Version) GetMajor() int32 {
	if m != nil {
		return m.Major
	}
	return 0
}

func (m *Version) GetMinor() int32 {
	if m != nil {
		return m.Minor
	}
	return 0
}

func (m *Version) GetPatch() int32 {
	if m != nil {
		return m.Patch
	}
	return 0
}

func (m *Version) GetSuffix() string {
	if m != nil {
		return m.Suffix
	}
	return ""
}

// An encoded Request is written to the ExtensionHandler's stdin.
type ExtensionHandlerRequest struct {
	// The OpenAPI descriptions that were explicitly listed on the command line.
	// The specifications will appear in the order they are specified to openapic.
	Wrapper *Wrapper `protobuf:"bytes,1,opt,name=wrapper" json:"wrapper,omitempty"`
	// The version number of openapi compiler.
	CompilerVersion *Version `protobuf:"bytes,3,opt,name=compiler_version,json=compilerVersion" json:"compiler_version,omitempty"`
}

func (m *ExtensionHandlerRequest) Reset()                    { *m = ExtensionHandlerRequest{} }
func (m *ExtensionHandlerRequest) String() string            { return proto.CompactTextString(m) }
func (*ExtensionHandlerRequest) ProtoMessage()               {}
func (*ExtensionHandlerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ExtensionHandlerRequest) GetWrapper() *Wrapper {
	if m != nil {
		return m.Wrapper
	}
	return nil
}

func (m *ExtensionHandlerRequest) GetCompilerVersion() *Version {
	if m != nil {
		return m.CompilerVersion
	}
	return nil
}

// The extensions writes an encoded ExtensionHandlerResponse to stdout.
type ExtensionHandlerResponse struct {
	// true if the extension is handled by the extension handler; false otherwise
	Handled bool `protobuf:"varint,1,opt,name=handled" json:"handled,omitempty"`
	// Error message.  If non-empty, the extension handling failed.
	// The extension handler process should exit with status code zero
	// even if it reports an error in this way.
	//
	// This should be used to indicate errors which prevent the extension from
	// operating as intended.  Errors which indicate a problem in gnostic
	// itself -- such as the input Document being unparseable -- should be
	// reported by writing a message to stderr and exiting with a non-zero
	// status code.
	Error []string `protobuf:"bytes,2,rep,name=error" json:"error,omitempty"`
	// text output
	Value *google_protobuf.Any `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
}

func (m *ExtensionHandlerResponse) Reset()                    { *m = ExtensionHandlerResponse{} }
func (m *ExtensionHandlerResponse) String() string            { return proto.CompactTextString(m) }
func (*ExtensionHandlerResponse) ProtoMessage()               {}
func (*ExtensionHandlerResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ExtensionHandlerResponse) GetHandled() bool {
	if m != nil {
		return m.Handled
	}
	return false
}

func (m *ExtensionHandlerResponse) GetError() []string {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *ExtensionHandlerResponse) GetValue() *google_protobuf.Any {
	if m != nil {
		return m.Value
	}
	return nil
}

type Wrapper struct {
	// version of the OpenAPI specification in which this extension was written.
	Version string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
	// Name of the extension
	ExtensionName string `protobuf:"bytes,2,opt,name=extension_name,json=extensionName" json:"extension_name,omitempty"`
	// Must be a valid yaml for the proto
	Yaml string `protobuf:"bytes,3,opt,name=yaml" json:"yaml,omitempty"`
}

func (m *Wrapper) Reset()                    { *m = Wrapper{} }
func (m *Wrapper) String() string            { return proto.CompactTextString(m) }
func (*Wrapper) ProtoMessage()               {}
func (*Wrapper) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Wrapper) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Wrapper) GetExtensionName() string {
	if m != nil {
		return m.ExtensionName
	}
	return ""
}

func (m *Wrapper) GetYaml() string {
	if m != nil {
		return m.Yaml
	}
	return ""
}

func init() {
	proto.RegisterType((*Version)(nil), "openapiextension.v1.Version")
	proto.RegisterType((*ExtensionHandlerRequest)(nil), "openapiextension.v1.ExtensionHandlerRequest")
	proto.RegisterType((*ExtensionHandlerResponse)(nil), "openapiextension.v1.ExtensionHandlerResponse")
	proto.RegisterType((*Wrapper)(nil), "openapiextension.v1.Wrapper")
}

func init() { proto.RegisterFile("extension.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x91, 0x4d, 0x4b, 0xf3, 0x40,
	0x1c, 0xc4, 0x49, 0xdf, 0xf2, 0x64, 0x1f, 0xb4, 0xb2, 0x16, 0x8d, 0xe2, 0xa1, 0x04, 0x84, 0x22,
	0xb8, 0xa5, 0x0a, 0xde, 0x5b, 0x28, 0xea, 0xc5, 0x96, 0x3d, 0xd4, 0x9b, 0x65, 0x9b, 0xfe, 0xdb,
	0x46, 0x92, 0xdd, 0x75, 0xf3, 0x62, 0xfb, 0x55, 0x3c, 0xfa, 0x49, 0x25, 0xbb, 0xd9, 0x7a, 0x50,
	0x6f, 0x99, 0x1f, 0x93, 0xfc, 0x67, 0x26, 0xa8, 0x0d, 0xdb, 0x0c, 0x78, 0x1a, 0x09, 0x4e, 0xa4,
	0x12, 0x99, 0xc0, 0xc7, 0x42, 0x02, 0x67, 0x32, 0xfa, 0xe6, 0xc5, 0xe0, 0xfc, 0x6c, 0x2d, 0xc4,
	0x3a, 0x86, 0xbe, 0xb6, 0x2c, 0xf2, 0x55, 0x9f, 0xf1, 0x9d, 0xf1, 0x07, 0x21, 0x72, 0x67, 0xa0,
	0x4a, 0x23, 0xee, 0xa0, 0x66, 0xc2, 0x5e, 0x85, 0xf2, 0x9d, 0xae, 0xd3, 0x6b, 0x52, 0x23, 0x34,
	0x8d, 0xb8, 0x50, 0x7e, 0xad, 0xa2, 0xa5, 0x28, 0xa9, 0x64, 0x59, 0xb8, 0xf1, 0xeb, 0x86, 0x6a,
	0x81, 0x4f, 0x50, 0x2b, 0xcd, 0x57, 0xab, 0x68, 0xeb, 0x37, 0xba, 0x4e, 0xcf, 0xa3, 0x95, 0x0a,
	0x3e, 0x1c, 0x74, 0x3a, 0xb6, 0x81, 0x1e, 0x18, 0x5f, 0xc6, 0xa0, 0x28, 0xbc, 0xe5, 0x90, 0x66,
	0xf8, 0x0e, 0xb9, 0xef, 0x8a, 0x49, 0x09, 0xe6, 0xee, 0xff, 0x9b, 0x0b, 0xf2, 0x4b, 0x05, 0xf2,
	0x6c, 0x3c, 0xd4, 0x9a, 0xf1, 0x3d, 0x3a, 0x0a, 0x45, 0x22, 0xa3, 0x18, 0xd4, 0xbc, 0x30, 0x0d,
	0x74, 0x98, 0xbf, 0x3e, 0x50, 0xb5, 0xa4, 0x6d, 0xfb, 0x56, 0x05, 0x82, 0x02, 0xf9, 0x3f, 0xb3,
	0xa5, 0x52, 0xf0, 0x14, 0xb0, 0x8f, 0xdc, 0x8d, 0x46, 0x4b, 0x1d, 0xee, 0x1f, 0xb5, 0xb2, 0x1c,
	0x00, 0x94, 0xd2, 0xb3, 0xd4, 0x7b, 0x1e, 0x35, 0x02, 0x5f, 0xa1, 0x66, 0xc1, 0xe2, 0x1c, 0xaa,
	0x24, 0x1d, 0x62, 0x86, 0x27, 0x76, 0x78, 0x32, 0xe4, 0x3b, 0x6a, 0x2c, 0xc1, 0x0b, 0x72, 0xab,
	0x52, 0xe5, 0x19, 0x5b, 0xc1, 0xd1, 0xc3, 0x59, 0x89, 0x2f, 0xd1, 0xe1, 0xbe, 0xc5, 0x9c, 0xb3,
	0x04, 0xf4, 0x6f, 0xf0, 0xe8, 0xc1, 0x9e, 0x3e, 0xb1, 0x04, 0x30, 0x46, 0x8d, 0x1d, 0x4b, 0x62,
	0x7d, 0xd6, 0xa3, 0xfa, 0x79, 0x74, 0x8d, 0xda, 0x42, 0xad, 0xed, 0x16, 0x21, 0x29, 0x06, 0x23,
	0x3c, 0x91, 0xc0, 0x87, 0xd3, 0xc7, 0x7d, 0xdf, 0xd9, 0x60, 0xea, 0x7c, 0xd6, 0xea, 0x93, 0xe1,
	0x78, 0xd1, 0xd2, 0x19, 0x6f, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x56, 0x40, 0x4d, 0x52,
	0x02, 0x00, 0x00,
}
