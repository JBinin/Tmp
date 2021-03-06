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
// source: google/logging/v2/logging_config.proto
// DO NOT EDIT!

package logging

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf5 "github.com/golang/protobuf/ptypes/empty"
import google_protobuf4 "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Available log entry formats. Log entries can be written to Stackdriver
// Logging in either format and can be exported in either format.
// Version 2 is the preferred format.
type LogSink_VersionFormat int32

const (
	// An unspecified format version that will default to V2.
	LogSink_VERSION_FORMAT_UNSPECIFIED LogSink_VersionFormat = 0
	// `LogEntry` version 2 format.
	LogSink_V2 LogSink_VersionFormat = 1
	// `LogEntry` version 1 format.
	LogSink_V1 LogSink_VersionFormat = 2
)

var LogSink_VersionFormat_name = map[int32]string{
	0: "VERSION_FORMAT_UNSPECIFIED",
	1: "V2",
	2: "V1",
}
var LogSink_VersionFormat_value = map[string]int32{
	"VERSION_FORMAT_UNSPECIFIED": 0,
	"V2": 1,
	"V1": 2,
}

func (x LogSink_VersionFormat) String() string {
	return proto.EnumName(LogSink_VersionFormat_name, int32(x))
}
func (LogSink_VersionFormat) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0, 0} }

// Describes a sink used to export log entries to one of the following
// destinations in any project: a Cloud Storage bucket, a BigQuery dataset, or a
// Cloud Pub/Sub topic.  A logs filter controls which log entries are
// exported. The sink must be created within a project or organization.
type LogSink struct {
	// Required. The client-assigned sink identifier, unique within the
	// project. Example: `"my-syslog-errors-to-pubsub"`.  Sink identifiers are
	// limited to 100 characters and can include only the following characters:
	// upper and lower-case alphanumeric characters, underscores, hyphens, and
	// periods.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Required. The export destination:
	//
	//     "storage.googleapis.com/[GCS_BUCKET]"
	//     "bigquery.googleapis.com/projects/[PROJECT_ID]/datasets/[DATASET]"
	//     "pubsub.googleapis.com/projects/[PROJECT_ID]/topics/[TOPIC_ID]"
	//
	// The sink's `writer_identity`, set when the sink is created, must
	// have permission to write to the destination or else the log
	// entries are not exported.  For more information, see
	// [Exporting Logs With Sinks](/logging/docs/api/tasks/exporting-logs).
	Destination string `protobuf:"bytes,3,opt,name=destination" json:"destination,omitempty"`
	// Optional.
	// An [advanced logs filter](/logging/docs/view/advanced_filters).  The only
	// exported log entries are those that are in the resource owning the sink and
	// that match the filter. The filter must use the log entry format specified
	// by the `output_version_format` parameter.  For example, in the v2 format:
	//
	//     logName="projects/[PROJECT_ID]/logs/[LOG_ID]" AND severity>=ERROR
	Filter string `protobuf:"bytes,5,opt,name=filter" json:"filter,omitempty"`
	// Optional. The log entry format to use for this sink's exported log
	// entries.  The v2 format is used by default.
	// **The v1 format is deprecated** and should be used only as part of a
	// migration effort to v2.
	// See [Migration to the v2 API](/logging/docs/api/v2/migration-to-v2).
	OutputVersionFormat LogSink_VersionFormat `protobuf:"varint,6,opt,name=output_version_format,json=outputVersionFormat,enum=google.logging.v2.LogSink_VersionFormat" json:"output_version_format,omitempty"`
	// Output only. An IAM identity&mdash;a service account or group&mdash;under
	// which Stackdriver Logging writes the exported log entries to the sink's
	// destination.  This field is set by
	// [sinks.create](/logging/docs/api/reference/rest/v2/projects.sinks/create)
	// and
	// [sinks.update](/logging/docs/api/reference/rest/v2/projects.sinks/update),
	// based on the setting of `unique_writer_identity` in those methods.
	//
	// Until you grant this identity write-access to the destination, log entry
	// exports from this sink will fail. For more information,
	// see [Granting access for a
	// resource](/iam/docs/granting-roles-to-service-accounts#granting_access_to_a_service_account_for_a_resource).
	// Consult the destination service's documentation to determine the
	// appropriate IAM roles to assign to the identity.
	WriterIdentity string `protobuf:"bytes,8,opt,name=writer_identity,json=writerIdentity" json:"writer_identity,omitempty"`
	// Optional. The time at which this sink will begin exporting log entries.
	// Log entries are exported only if their timestamp is not earlier than the
	// start time.  The default value of this field is the time the sink is
	// created or updated.
	StartTime *google_protobuf4.Timestamp `protobuf:"bytes,10,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	// Optional. The time at which this sink will stop exporting log entries.  Log
	// entries are exported only if their timestamp is earlier than the end time.
	// If this field is not supplied, there is no end time.  If both a start time
	// and an end time are provided, then the end time must be later than the
	// start time.
	EndTime *google_protobuf4.Timestamp `protobuf:"bytes,11,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
}

func (m *LogSink) Reset()                    { *m = LogSink{} }
func (m *LogSink) String() string            { return proto.CompactTextString(m) }
func (*LogSink) ProtoMessage()               {}
func (*LogSink) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *LogSink) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LogSink) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *LogSink) GetFilter() string {
	if m != nil {
		return m.Filter
	}
	return ""
}

func (m *LogSink) GetOutputVersionFormat() LogSink_VersionFormat {
	if m != nil {
		return m.OutputVersionFormat
	}
	return LogSink_VERSION_FORMAT_UNSPECIFIED
}

func (m *LogSink) GetWriterIdentity() string {
	if m != nil {
		return m.WriterIdentity
	}
	return ""
}

func (m *LogSink) GetStartTime() *google_protobuf4.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *LogSink) GetEndTime() *google_protobuf4.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

// The parameters to `ListSinks`.
type ListSinksRequest struct {
	// Required. The parent resource whose sinks are to be listed.
	// Examples: `"projects/my-logging-project"`, `"organizations/123456789"`.
	Parent string `protobuf:"bytes,1,opt,name=parent" json:"parent,omitempty"`
	// Optional. If present, then retrieve the next batch of results from the
	// preceding call to this method.  `pageToken` must be the value of
	// `nextPageToken` from the previous response.  The values of other method
	// parameters should be identical to those in the previous call.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken" json:"page_token,omitempty"`
	// Optional. The maximum number of results to return from this request.
	// Non-positive values are ignored.  The presence of `nextPageToken` in the
	// response indicates that more results might be available.
	PageSize int32 `protobuf:"varint,3,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
}

func (m *ListSinksRequest) Reset()                    { *m = ListSinksRequest{} }
func (m *ListSinksRequest) String() string            { return proto.CompactTextString(m) }
func (*ListSinksRequest) ProtoMessage()               {}
func (*ListSinksRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ListSinksRequest) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *ListSinksRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

func (m *ListSinksRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

// Result returned from `ListSinks`.
type ListSinksResponse struct {
	// A list of sinks.
	Sinks []*LogSink `protobuf:"bytes,1,rep,name=sinks" json:"sinks,omitempty"`
	// If there might be more results than appear in this response, then
	// `nextPageToken` is included.  To get the next set of results, call the same
	// method again using the value of `nextPageToken` as `pageToken`.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken" json:"next_page_token,omitempty"`
}

func (m *ListSinksResponse) Reset()                    { *m = ListSinksResponse{} }
func (m *ListSinksResponse) String() string            { return proto.CompactTextString(m) }
func (*ListSinksResponse) ProtoMessage()               {}
func (*ListSinksResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *ListSinksResponse) GetSinks() []*LogSink {
	if m != nil {
		return m.Sinks
	}
	return nil
}

func (m *ListSinksResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

// The parameters to `GetSink`.
type GetSinkRequest struct {
	// Required. The parent resource name of the sink:
	//
	//     "projects/[PROJECT_ID]/sinks/[SINK_ID]"
	//     "organizations/[ORGANIZATION_ID]/sinks/[SINK_ID]"
	//
	// Example: `"projects/my-project-id/sinks/my-sink-id"`.
	SinkName string `protobuf:"bytes,1,opt,name=sink_name,json=sinkName" json:"sink_name,omitempty"`
}

func (m *GetSinkRequest) Reset()                    { *m = GetSinkRequest{} }
func (m *GetSinkRequest) String() string            { return proto.CompactTextString(m) }
func (*GetSinkRequest) ProtoMessage()               {}
func (*GetSinkRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *GetSinkRequest) GetSinkName() string {
	if m != nil {
		return m.SinkName
	}
	return ""
}

// The parameters to `CreateSink`.
type CreateSinkRequest struct {
	// Required. The resource in which to create the sink:
	//
	//     "projects/[PROJECT_ID]"
	//     "organizations/[ORGANIZATION_ID]"
	//
	// Examples: `"projects/my-logging-project"`, `"organizations/123456789"`.
	Parent string `protobuf:"bytes,1,opt,name=parent" json:"parent,omitempty"`
	// Required. The new sink, whose `name` parameter is a sink identifier that
	// is not already in use.
	Sink *LogSink `protobuf:"bytes,2,opt,name=sink" json:"sink,omitempty"`
	// Optional. Determines the kind of IAM identity returned as `writer_identity`
	// in the new sink.  If this value is omitted or set to false, and if the
	// sink's parent is a project, then the value returned as `writer_identity` is
	// `cloud-logs@google.com`, the same identity used before the addition of
	// writer identities to this API. The sink's destination must be in the same
	// project as the sink itself.
	//
	// If this field is set to true, or if the sink is owned by a non-project
	// resource such as an organization, then the value of `writer_identity` will
	// be a unique service account used only for exports from the new sink.  For
	// more information, see `writer_identity` in [LogSink][google.logging.v2.LogSink].
	UniqueWriterIdentity bool `protobuf:"varint,3,opt,name=unique_writer_identity,json=uniqueWriterIdentity" json:"unique_writer_identity,omitempty"`
}

func (m *CreateSinkRequest) Reset()                    { *m = CreateSinkRequest{} }
func (m *CreateSinkRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateSinkRequest) ProtoMessage()               {}
func (*CreateSinkRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *CreateSinkRequest) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *CreateSinkRequest) GetSink() *LogSink {
	if m != nil {
		return m.Sink
	}
	return nil
}

func (m *CreateSinkRequest) GetUniqueWriterIdentity() bool {
	if m != nil {
		return m.UniqueWriterIdentity
	}
	return false
}

// The parameters to `UpdateSink`.
type UpdateSinkRequest struct {
	// Required. The full resource name of the sink to update, including the
	// parent resource and the sink identifier:
	//
	//     "projects/[PROJECT_ID]/sinks/[SINK_ID]"
	//     "organizations/[ORGANIZATION_ID]/sinks/[SINK_ID]"
	//
	// Example: `"projects/my-project-id/sinks/my-sink-id"`.
	SinkName string `protobuf:"bytes,1,opt,name=sink_name,json=sinkName" json:"sink_name,omitempty"`
	// Required. The updated sink, whose name is the same identifier that appears
	// as part of `sink_name`.  If `sink_name` does not exist, then
	// this method creates a new sink.
	Sink *LogSink `protobuf:"bytes,2,opt,name=sink" json:"sink,omitempty"`
	// Optional. See
	// [sinks.create](/logging/docs/api/reference/rest/v2/projects.sinks/create)
	// for a description of this field.  When updating a sink, the effect of this
	// field on the value of `writer_identity` in the updated sink depends on both
	// the old and new values of this field:
	//
	// +   If the old and new values of this field are both false or both true,
	//     then there is no change to the sink's `writer_identity`.
	// +   If the old value was false and the new value is true, then
	//     `writer_identity` is changed to a unique service account.
	// +   It is an error if the old value was true and the new value is false.
	UniqueWriterIdentity bool `protobuf:"varint,3,opt,name=unique_writer_identity,json=uniqueWriterIdentity" json:"unique_writer_identity,omitempty"`
}

func (m *UpdateSinkRequest) Reset()                    { *m = UpdateSinkRequest{} }
func (m *UpdateSinkRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateSinkRequest) ProtoMessage()               {}
func (*UpdateSinkRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *UpdateSinkRequest) GetSinkName() string {
	if m != nil {
		return m.SinkName
	}
	return ""
}

func (m *UpdateSinkRequest) GetSink() *LogSink {
	if m != nil {
		return m.Sink
	}
	return nil
}

func (m *UpdateSinkRequest) GetUniqueWriterIdentity() bool {
	if m != nil {
		return m.UniqueWriterIdentity
	}
	return false
}

// The parameters to `DeleteSink`.
type DeleteSinkRequest struct {
	// Required. The full resource name of the sink to delete, including the
	// parent resource and the sink identifier:
	//
	//     "projects/[PROJECT_ID]/sinks/[SINK_ID]"
	//     "organizations/[ORGANIZATION_ID]/sinks/[SINK_ID]"
	//
	// It is an error if the sink does not exist.  Example:
	// `"projects/my-project-id/sinks/my-sink-id"`.  It is an error if
	// the sink does not exist.
	SinkName string `protobuf:"bytes,1,opt,name=sink_name,json=sinkName" json:"sink_name,omitempty"`
}

func (m *DeleteSinkRequest) Reset()                    { *m = DeleteSinkRequest{} }
func (m *DeleteSinkRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteSinkRequest) ProtoMessage()               {}
func (*DeleteSinkRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

func (m *DeleteSinkRequest) GetSinkName() string {
	if m != nil {
		return m.SinkName
	}
	return ""
}

func init() {
	proto.RegisterType((*LogSink)(nil), "google.logging.v2.LogSink")
	proto.RegisterType((*ListSinksRequest)(nil), "google.logging.v2.ListSinksRequest")
	proto.RegisterType((*ListSinksResponse)(nil), "google.logging.v2.ListSinksResponse")
	proto.RegisterType((*GetSinkRequest)(nil), "google.logging.v2.GetSinkRequest")
	proto.RegisterType((*CreateSinkRequest)(nil), "google.logging.v2.CreateSinkRequest")
	proto.RegisterType((*UpdateSinkRequest)(nil), "google.logging.v2.UpdateSinkRequest")
	proto.RegisterType((*DeleteSinkRequest)(nil), "google.logging.v2.DeleteSinkRequest")
	proto.RegisterEnum("google.logging.v2.LogSink_VersionFormat", LogSink_VersionFormat_name, LogSink_VersionFormat_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ConfigServiceV2 service

type ConfigServiceV2Client interface {
	// Lists sinks.
	ListSinks(ctx context.Context, in *ListSinksRequest, opts ...grpc.CallOption) (*ListSinksResponse, error)
	// Gets a sink.
	GetSink(ctx context.Context, in *GetSinkRequest, opts ...grpc.CallOption) (*LogSink, error)
	// Creates a sink that exports specified log entries to a destination.  The
	// export of newly-ingested log entries begins immediately, unless the current
	// time is outside the sink's start and end times or the sink's
	// `writer_identity` is not permitted to write to the destination.  A sink can
	// export log entries only from the resource owning the sink.
	CreateSink(ctx context.Context, in *CreateSinkRequest, opts ...grpc.CallOption) (*LogSink, error)
	// Updates a sink. If the named sink doesn't exist, then this method is
	// identical to
	// [sinks.create](/logging/docs/api/reference/rest/v2/projects.sinks/create).
	// If the named sink does exist, then this method replaces the following
	// fields in the existing sink with values from the new sink: `destination`,
	// `filter`, `output_version_format`, `start_time`, and `end_time`.
	// The updated filter might also have a new `writer_identity`; see the
	// `unique_writer_identity` field.
	UpdateSink(ctx context.Context, in *UpdateSinkRequest, opts ...grpc.CallOption) (*LogSink, error)
	// Deletes a sink. If the sink has a unique `writer_identity`, then that
	// service account is also deleted.
	DeleteSink(ctx context.Context, in *DeleteSinkRequest, opts ...grpc.CallOption) (*google_protobuf5.Empty, error)
}

type configServiceV2Client struct {
	cc *grpc.ClientConn
}

func NewConfigServiceV2Client(cc *grpc.ClientConn) ConfigServiceV2Client {
	return &configServiceV2Client{cc}
}

func (c *configServiceV2Client) ListSinks(ctx context.Context, in *ListSinksRequest, opts ...grpc.CallOption) (*ListSinksResponse, error) {
	out := new(ListSinksResponse)
	err := grpc.Invoke(ctx, "/google.logging.v2.ConfigServiceV2/ListSinks", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceV2Client) GetSink(ctx context.Context, in *GetSinkRequest, opts ...grpc.CallOption) (*LogSink, error) {
	out := new(LogSink)
	err := grpc.Invoke(ctx, "/google.logging.v2.ConfigServiceV2/GetSink", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceV2Client) CreateSink(ctx context.Context, in *CreateSinkRequest, opts ...grpc.CallOption) (*LogSink, error) {
	out := new(LogSink)
	err := grpc.Invoke(ctx, "/google.logging.v2.ConfigServiceV2/CreateSink", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceV2Client) UpdateSink(ctx context.Context, in *UpdateSinkRequest, opts ...grpc.CallOption) (*LogSink, error) {
	out := new(LogSink)
	err := grpc.Invoke(ctx, "/google.logging.v2.ConfigServiceV2/UpdateSink", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceV2Client) DeleteSink(ctx context.Context, in *DeleteSinkRequest, opts ...grpc.CallOption) (*google_protobuf5.Empty, error) {
	out := new(google_protobuf5.Empty)
	err := grpc.Invoke(ctx, "/google.logging.v2.ConfigServiceV2/DeleteSink", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ConfigServiceV2 service

type ConfigServiceV2Server interface {
	// Lists sinks.
	ListSinks(context.Context, *ListSinksRequest) (*ListSinksResponse, error)
	// Gets a sink.
	GetSink(context.Context, *GetSinkRequest) (*LogSink, error)
	// Creates a sink that exports specified log entries to a destination.  The
	// export of newly-ingested log entries begins immediately, unless the current
	// time is outside the sink's start and end times or the sink's
	// `writer_identity` is not permitted to write to the destination.  A sink can
	// export log entries only from the resource owning the sink.
	CreateSink(context.Context, *CreateSinkRequest) (*LogSink, error)
	// Updates a sink. If the named sink doesn't exist, then this method is
	// identical to
	// [sinks.create](/logging/docs/api/reference/rest/v2/projects.sinks/create).
	// If the named sink does exist, then this method replaces the following
	// fields in the existing sink with values from the new sink: `destination`,
	// `filter`, `output_version_format`, `start_time`, and `end_time`.
	// The updated filter might also have a new `writer_identity`; see the
	// `unique_writer_identity` field.
	UpdateSink(context.Context, *UpdateSinkRequest) (*LogSink, error)
	// Deletes a sink. If the sink has a unique `writer_identity`, then that
	// service account is also deleted.
	DeleteSink(context.Context, *DeleteSinkRequest) (*google_protobuf5.Empty, error)
}

func RegisterConfigServiceV2Server(s *grpc.Server, srv ConfigServiceV2Server) {
	s.RegisterService(&_ConfigServiceV2_serviceDesc, srv)
}

func _ConfigServiceV2_ListSinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceV2Server).ListSinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.logging.v2.ConfigServiceV2/ListSinks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceV2Server).ListSinks(ctx, req.(*ListSinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigServiceV2_GetSink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceV2Server).GetSink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.logging.v2.ConfigServiceV2/GetSink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceV2Server).GetSink(ctx, req.(*GetSinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigServiceV2_CreateSink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceV2Server).CreateSink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.logging.v2.ConfigServiceV2/CreateSink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceV2Server).CreateSink(ctx, req.(*CreateSinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigServiceV2_UpdateSink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceV2Server).UpdateSink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.logging.v2.ConfigServiceV2/UpdateSink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceV2Server).UpdateSink(ctx, req.(*UpdateSinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigServiceV2_DeleteSink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceV2Server).DeleteSink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.logging.v2.ConfigServiceV2/DeleteSink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceV2Server).DeleteSink(ctx, req.(*DeleteSinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConfigServiceV2_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.logging.v2.ConfigServiceV2",
	HandlerType: (*ConfigServiceV2Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSinks",
			Handler:    _ConfigServiceV2_ListSinks_Handler,
		},
		{
			MethodName: "GetSink",
			Handler:    _ConfigServiceV2_GetSink_Handler,
		},
		{
			MethodName: "CreateSink",
			Handler:    _ConfigServiceV2_CreateSink_Handler,
		},
		{
			MethodName: "UpdateSink",
			Handler:    _ConfigServiceV2_UpdateSink_Handler,
		},
		{
			MethodName: "DeleteSink",
			Handler:    _ConfigServiceV2_DeleteSink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/logging/v2/logging_config.proto",
}

func init() { proto.RegisterFile("google/logging/v2/logging_config.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 787 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xed, 0x4e, 0x33, 0x45,
	0x14, 0x76, 0x0b, 0x94, 0xf6, 0x90, 0x17, 0xda, 0xd1, 0x17, 0x37, 0x8b, 0xaf, 0xd4, 0x15, 0xb0,
	0xa9, 0x71, 0x17, 0x57, 0x4d, 0xfc, 0x88, 0x31, 0x52, 0x0a, 0x69, 0x82, 0xd0, 0x6c, 0xa1, 0x26,
	0xc6, 0x64, 0xb3, 0xb4, 0xd3, 0xcd, 0x48, 0x77, 0x66, 0xd9, 0x9d, 0x56, 0x81, 0x90, 0xa8, 0x77,
	0x40, 0x4c, 0xbc, 0x08, 0x6f, 0xc7, 0x5b, 0xf0, 0x3a, 0x8c, 0x99, 0x99, 0x2d, 0xb4, 0xdd, 0x5a,
	0xf9, 0xe5, 0xaf, 0x9d, 0xf3, 0x35, 0xcf, 0x73, 0xce, 0x79, 0x32, 0x0b, 0x7b, 0x01, 0x63, 0xc1,
	0x00, 0xdb, 0x03, 0x16, 0x04, 0x84, 0x06, 0xf6, 0xc8, 0x19, 0x1f, 0xbd, 0x2e, 0xa3, 0x7d, 0x12,
	0x58, 0x51, 0xcc, 0x38, 0x43, 0x65, 0x95, 0x67, 0xa5, 0x41, 0x6b, 0xe4, 0x18, 0x6f, 0xa5, 0xa5,
	0x7e, 0x44, 0x6c, 0x9f, 0x52, 0xc6, 0x7d, 0x4e, 0x18, 0x4d, 0x54, 0x81, 0xb1, 0x95, 0x46, 0xa5,
	0x75, 0x39, 0xec, 0xdb, 0x38, 0x8c, 0xf8, 0x4d, 0x1a, 0xdc, 0x9e, 0x0d, 0x72, 0x12, 0xe2, 0x84,
	0xfb, 0x61, 0xa4, 0x12, 0xcc, 0x87, 0x25, 0x58, 0x3d, 0x61, 0x41, 0x9b, 0xd0, 0x2b, 0x84, 0x60,
	0x99, 0xfa, 0x21, 0xd6, 0xb5, 0x8a, 0x56, 0x2d, 0xba, 0xf2, 0x8c, 0x2a, 0xb0, 0xd6, 0xc3, 0x09,
	0x27, 0x54, 0x62, 0xea, 0x4b, 0x32, 0x34, 0xe9, 0x42, 0x9b, 0x90, 0xef, 0x93, 0x01, 0xc7, 0xb1,
	0xbe, 0x22, 0x83, 0xa9, 0x85, 0xbe, 0x87, 0x97, 0x6c, 0xc8, 0xa3, 0x21, 0xf7, 0x46, 0x38, 0x4e,
	0x08, 0xa3, 0x5e, 0x9f, 0xc5, 0xa1, 0xcf, 0xf5, 0x7c, 0x45, 0xab, 0xae, 0x3b, 0x55, 0x2b, 0xd3,
	0xa8, 0x95, 0x12, 0xb1, 0x3a, 0xaa, 0xe0, 0x48, 0xe6, 0xbb, 0xaf, 0xab, 0x6b, 0xa6, 0x9c, 0xe8,
	0x3d, 0xd8, 0xf8, 0x31, 0x26, 0x1c, 0xc7, 0x1e, 0xe9, 0x61, 0xca, 0x09, 0xbf, 0xd1, 0x0b, 0x12,
	0x7e, 0x5d, 0xb9, 0x9b, 0xa9, 0x17, 0x7d, 0x06, 0x90, 0x70, 0x3f, 0xe6, 0x9e, 0xe8, 0x5c, 0x87,
	0x8a, 0x56, 0x5d, 0x73, 0x8c, 0x31, 0xf6, 0x78, 0x2c, 0xd6, 0xf9, 0x78, 0x2c, 0x6e, 0x51, 0x66,
	0x0b, 0x1b, 0x7d, 0x02, 0x05, 0x4c, 0x7b, 0xaa, 0x70, 0xed, 0x3f, 0x0b, 0x57, 0x31, 0xed, 0x09,
	0xcb, 0xfc, 0x0a, 0x5e, 0x4c, 0x73, 0x7d, 0x1b, 0x8c, 0x4e, 0xc3, 0x6d, 0x37, 0xcf, 0x4e, 0xbd,
	0xa3, 0x33, 0xf7, 0x9b, 0xaf, 0xcf, 0xbd, 0x8b, 0xd3, 0x76, 0xab, 0x51, 0x6f, 0x1e, 0x35, 0x1b,
	0x87, 0xa5, 0xd7, 0x50, 0x1e, 0x72, 0x1d, 0xa7, 0xa4, 0xc9, 0xef, 0x87, 0xa5, 0x9c, 0xd9, 0x87,
	0xd2, 0x09, 0x49, 0xb8, 0x18, 0x45, 0xe2, 0xe2, 0xeb, 0x21, 0x4e, 0xb8, 0x98, 0x72, 0xe4, 0xc7,
	0x98, 0xf2, 0x74, 0x3b, 0xa9, 0x85, 0x5e, 0x01, 0x44, 0x7e, 0x80, 0x3d, 0xce, 0xae, 0x30, 0xd5,
	0x73, 0x32, 0x56, 0x14, 0x9e, 0x73, 0xe1, 0x40, 0x5b, 0x20, 0x0d, 0x2f, 0x21, 0xb7, 0x58, 0x2e,
	0x6f, 0xc5, 0x2d, 0x08, 0x47, 0x9b, 0xdc, 0x62, 0x33, 0x84, 0xf2, 0x04, 0x4e, 0x12, 0x31, 0x9a,
	0x60, 0xb4, 0x0f, 0x2b, 0x89, 0x70, 0xe8, 0x5a, 0x65, 0x69, 0xb2, 0xe3, 0xec, 0x9a, 0x5c, 0x95,
	0x88, 0xf6, 0x60, 0x83, 0xe2, 0x9f, 0xb8, 0x97, 0xe1, 0xf1, 0x42, 0xb8, 0x5b, 0x63, 0x2e, 0xe6,
	0x07, 0xb0, 0x7e, 0x8c, 0x25, 0xda, 0xb8, 0xa9, 0x2d, 0x28, 0x8a, 0x2b, 0xbc, 0x09, 0xd5, 0x15,
	0x84, 0xe3, 0xd4, 0x0f, 0xb1, 0xf9, 0xa0, 0x41, 0xb9, 0x1e, 0x63, 0x9f, 0xe3, 0xc9, 0x92, 0x7f,
	0x9b, 0x83, 0x05, 0xcb, 0xa2, 0x52, 0x22, 0x2f, 0x66, 0x2d, 0xf3, 0xd0, 0xc7, 0xb0, 0x39, 0xa4,
	0xe4, 0x7a, 0x88, 0xbd, 0x59, 0x19, 0x89, 0x29, 0x15, 0xdc, 0x37, 0x54, 0xf4, 0xdb, 0x29, 0x31,
	0x99, 0xbf, 0x6b, 0x50, 0xbe, 0x88, 0x7a, 0x33, 0x9c, 0x16, 0xb5, 0xf1, 0x3f, 0x11, 0xdb, 0x87,
	0xf2, 0x21, 0x1e, 0xe0, 0xe7, 0xf3, 0x72, 0xfe, 0x5e, 0x86, 0x8d, 0xba, 0x7c, 0x78, 0xda, 0x38,
	0x1e, 0x91, 0x2e, 0xee, 0x38, 0xe8, 0x1e, 0x8a, 0x8f, 0x82, 0x40, 0xef, 0xce, 0xa3, 0x3a, 0x23,
	0x4b, 0x63, 0x67, 0x71, 0x92, 0xd2, 0x94, 0xb9, 0xfb, 0xeb, 0x9f, 0x7f, 0xfd, 0x96, 0xdb, 0x46,
	0xaf, 0xc4, 0xab, 0x77, 0xa7, 0x36, 0xf6, 0x65, 0x14, 0xb3, 0x1f, 0x70, 0x97, 0x27, 0x76, 0xed,
	0xde, 0x56, 0x42, 0xe2, 0xb0, 0x9a, 0x0a, 0x04, 0xbd, 0x33, 0xe7, 0xde, 0x69, 0xf1, 0x18, 0x0b,
	0x46, 0x69, 0xd6, 0x24, 0xe0, 0x0e, 0x32, 0x25, 0xe0, 0xe3, 0x10, 0x26, 0x30, 0x15, 0xa4, 0x5d,
	0xbb, 0x47, 0x77, 0x00, 0x4f, 0x32, 0x43, 0xf3, 0x1a, 0xca, 0xa8, 0x70, 0x21, 0xf6, 0xfb, 0x12,
	0x7b, 0xd7, 0x5c, 0xdc, 0xec, 0xe7, 0x6a, 0xdb, 0x3f, 0x6b, 0x00, 0x4f, 0x82, 0x9a, 0x8b, 0x9e,
	0xd1, 0xdb, 0x42, 0xf4, 0x7d, 0x89, 0x5e, 0x33, 0x9e, 0xd1, 0x79, 0x4a, 0x61, 0x04, 0xf0, 0x24,
	0x9d, 0xb9, 0x0c, 0x32, 0xca, 0x32, 0x36, 0x33, 0xef, 0x60, 0x43, 0xfc, 0x74, 0xc6, 0x73, 0xaf,
	0x3d, 0x03, 0xfd, 0xe0, 0x17, 0x0d, 0x5e, 0x76, 0x59, 0x98, 0xc5, 0x3b, 0x40, 0x27, 0xea, 0xac,
	0xe4, 0xd9, 0x12, 0x10, 0x2d, 0xed, 0xbb, 0x4f, 0xd3, 0xc4, 0x80, 0x0d, 0x7c, 0x1a, 0x58, 0x2c,
	0x0e, 0xec, 0x00, 0x53, 0x49, 0xc0, 0x56, 0x21, 0x3f, 0x22, 0xc9, 0xc4, 0xff, 0xf5, 0x8b, 0xf4,
	0xf8, 0x47, 0xee, 0xcd, 0x63, 0x55, 0x5a, 0x1f, 0xb0, 0x61, 0xcf, 0x4a, 0x6f, 0xb7, 0x3a, 0xce,
	0x65, 0x5e, 0x96, 0x7f, 0xf4, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdc, 0x9e, 0xa5, 0xee, 0x9d,
	0x07, 0x00, 0x00,
}
