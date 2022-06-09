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
// source: google/api/monitored_resource.proto
// DO NOT EDIT!

/*
Package monitoredres is a generated protocol buffer package.

It is generated from these files:
	google/api/monitored_resource.proto

It has these top-level messages:
	MonitoredResourceDescriptor
	MonitoredResource
*/
package monitoredres

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_api "google.golang.org/genproto/googleapis/api/label"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// An object that describes the schema of a [MonitoredResource][google.api.MonitoredResource] object using a
// type name and a set of labels.  For example, the monitored resource
// descriptor for Google Compute Engine VM instances has a type of
// `"gce_instance"` and specifies the use of the labels `"instance_id"` and
// `"zone"` to identify particular VM instances.
//
// Different APIs can support different monitored resource types. APIs generally
// provide a `list` method that returns the monitored resource descriptors used
// by the API.
type MonitoredResourceDescriptor struct {
	// Optional. The resource name of the monitored resource descriptor:
	// `"projects/{project_id}/monitoredResourceDescriptors/{type}"` where
	// {type} is the value of the `type` field in this object and
	// {project_id} is a project ID that provides API-specific context for
	// accessing the type.  APIs that do not use project information can use the
	// resource name format `"monitoredResourceDescriptors/{type}"`.
	Name string `protobuf:"bytes,5,opt,name=name" json:"name,omitempty"`
	// Required. The monitored resource type. For example, the type
	// `"cloudsql_database"` represents databases in Google Cloud SQL.
	// The maximum length of this value is 256 characters.
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	// Optional. A concise name for the monitored resource type that might be
	// displayed in user interfaces. It should be a Title Cased Noun Phrase,
	// without any article or other determiners. For example,
	// `"Google Cloud SQL Database"`.
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName" json:"display_name,omitempty"`
	// Optional. A detailed description of the monitored resource type that might
	// be used in documentation.
	Description string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
	// Required. A set of labels used to describe instances of this monitored
	// resource type. For example, an individual Google Cloud SQL database is
	// identified by values for the labels `"database_id"` and `"zone"`.
	Labels []*google_api.LabelDescriptor `protobuf:"bytes,4,rep,name=labels" json:"labels,omitempty"`
}

func (m *MonitoredResourceDescriptor) Reset()                    { *m = MonitoredResourceDescriptor{} }
func (m *MonitoredResourceDescriptor) String() string            { return proto.CompactTextString(m) }
func (*MonitoredResourceDescriptor) ProtoMessage()               {}
func (*MonitoredResourceDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MonitoredResourceDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetLabels() []*google_api.LabelDescriptor {
	if m != nil {
		return m.Labels
	}
	return nil
}

// An object representing a resource that can be used for monitoring, logging,
// billing, or other purposes. Examples include virtual machine instances,
// databases, and storage devices such as disks. The `type` field identifies a
// [MonitoredResourceDescriptor][google.api.MonitoredResourceDescriptor] object that describes the resource's
// schema. Information in the `labels` field identifies the actual resource and
// its attributes according to the schema. For example, a particular Compute
// Engine VM instance could be represented by the following object, because the
// [MonitoredResourceDescriptor][google.api.MonitoredResourceDescriptor] for `"gce_instance"` has labels
// `"instance_id"` and `"zone"`:
//
//     { "type": "gce_instance",
//       "labels": { "instance_id": "12345678901234",
//                   "zone": "us-central1-a" }}
type MonitoredResource struct {
	// Required. The monitored resource type. This field must match
	// the `type` field of a [MonitoredResourceDescriptor][google.api.MonitoredResourceDescriptor] object. For
	// example, the type of a Cloud SQL database is `"cloudsql_database"`.
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	// Required. Values for all of the labels listed in the associated monitored
	// resource descriptor. For example, Cloud SQL databases use the labels
	// `"database_id"` and `"zone"`.
	Labels map[string]string `protobuf:"bytes,2,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *MonitoredResource) Reset()                    { *m = MonitoredResource{} }
func (m *MonitoredResource) String() string            { return proto.CompactTextString(m) }
func (*MonitoredResource) ProtoMessage()               {}
func (*MonitoredResource) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MonitoredResource) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *MonitoredResource) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func init() {
	proto.RegisterType((*MonitoredResourceDescriptor)(nil), "google.api.MonitoredResourceDescriptor")
	proto.RegisterType((*MonitoredResource)(nil), "google.api.MonitoredResource")
}

func init() { proto.RegisterFile("google/api/monitored_resource.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0x4b, 0x4b, 0x3b, 0x31,
	0x10, 0x27, 0xdb, 0x07, 0xfc, 0x67, 0xff, 0x88, 0x06, 0x29, 0x4b, 0x7b, 0xa9, 0xf5, 0x52, 0x2f,
	0xbb, 0x60, 0x2f, 0x3e, 0x4e, 0xad, 0x8a, 0x08, 0x2a, 0xa5, 0x47, 0x2f, 0x25, 0x6d, 0xc3, 0x12,
	0xdc, 0x66, 0x42, 0xb2, 0x15, 0xf6, 0xeb, 0x08, 0x7e, 0x0e, 0xbf, 0x96, 0x47, 0xc9, 0xa3, 0x76,
	0xa5, 0xde, 0x26, 0xbf, 0xf9, 0x3d, 0x66, 0x32, 0x70, 0x9a, 0x23, 0xe6, 0x05, 0xcf, 0x98, 0x12,
	0xd9, 0x1a, 0xa5, 0x28, 0x51, 0xf3, 0xd5, 0x5c, 0x73, 0x83, 0x1b, 0xbd, 0xe4, 0xa9, 0xd2, 0x58,
	0x22, 0x05, 0x4f, 0x4a, 0x99, 0x12, 0xdd, 0x4e, 0x4d, 0x50, 0xb0, 0x05, 0x2f, 0x3c, 0x67, 0xf0,
	0x49, 0xa0, 0xf7, 0xb4, 0x35, 0x98, 0x05, 0xfd, 0x2d, 0x37, 0x4b, 0x2d, 0x54, 0x89, 0x9a, 0x52,
	0x68, 0x4a, 0xb6, 0xe6, 0x49, 0xab, 0x4f, 0x86, 0xff, 0x66, 0xae, 0xb6, 0x58, 0x59, 0x29, 0x9e,
	0x10, 0x8f, 0xd9, 0x9a, 0x9e, 0xc0, 0xff, 0x95, 0x30, 0xaa, 0x60, 0xd5, 0xdc, 0xf1, 0x23, 0xd7,
	0x8b, 0x03, 0xf6, 0x6c, 0x65, 0x7d, 0x88, 0x57, 0xc1, 0x58, 0xa0, 0x4c, 0x1a, 0x81, 0xb1, 0x83,
	0xe8, 0x08, 0xda, 0x6e, 0x36, 0x93, 0x34, 0xfb, 0x8d, 0x61, 0x7c, 0xde, 0x4b, 0x77, 0x1b, 0xa4,
	0x8f, 0xb6, 0xb3, 0x9b, 0x6c, 0x16, 0xa8, 0x83, 0x0f, 0x02, 0x47, 0x7b, 0x1b, 0xfc, 0x39, 0xe3,
	0xf8, 0xc7, 0x3e, 0x72, 0xf6, 0x67, 0x75, 0xfb, 0x3d, 0x0b, 0x1f, 0x68, 0xee, 0x64, 0xa9, 0xab,
	0x6d, 0x58, 0xf7, 0x12, 0xe2, 0x1a, 0x4c, 0x0f, 0xa1, 0xf1, 0xca, 0xab, 0x10, 0x62, 0x4b, 0x7a,
	0x0c, 0xad, 0x37, 0x56, 0x6c, 0xb6, 0x1f, 0xe0, 0x1f, 0x57, 0xd1, 0x05, 0x99, 0x54, 0x70, 0xb0,
	0xc4, 0x75, 0x2d, 0x72, 0xd2, 0xd9, 0xcb, 0x9c, 0xda, 0x9b, 0x4c, 0xc9, 0xcb, 0x4d, 0x60, 0xe5,
	0x58, 0x30, 0x99, 0xa7, 0xa8, 0xf3, 0x2c, 0xe7, 0xd2, 0x5d, 0x2c, 0xf3, 0x2d, 0xa6, 0x84, 0xf9,
	0x7d, 0x7d, 0xcd, 0xcd, 0x75, 0xfd, 0xf1, 0x45, 0xc8, 0x7b, 0xd4, 0xbc, 0x1f, 0x4f, 0x1f, 0x16,
	0x6d, 0xa7, 0x1c, 0x7d, 0x07, 0x00, 0x00, 0xff, 0xff, 0xf8, 0xfb, 0xfb, 0x11, 0x36, 0x02, 0x00,
	0x00,
}
