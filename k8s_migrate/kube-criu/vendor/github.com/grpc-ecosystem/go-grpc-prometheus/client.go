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
// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for client-side gRPC.

package grpc_prometheus

import (
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	monitor := newClientReporter(Unary, method)
	monitor.SentMessage()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		monitor.ReceivedMessage()
	}
	monitor.Handled(grpc.Code(err))
	return err
}

// StreamServerInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	monitor := newClientReporter(clientStreamType(desc), method)
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		monitor.Handled(grpc.Code(err))
		return nil, err
	}
	return &monitoredClientStream{clientStream, monitor}, nil
}

func clientStreamType(desc *grpc.StreamDesc) grpcType {
	if desc.ClientStreams && !desc.ServerStreams {
		return ClientStream
	} else if !desc.ClientStreams && desc.ServerStreams {
		return ServerStream
	}
	return BidiStream
}

// monitoredClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type monitoredClientStream struct {
	grpc.ClientStream
	monitor *clientReporter
}

func (s *monitoredClientStream) SendMsg(m interface{}) error {
	err := s.ClientStream.SendMsg(m)
	if err == nil {
		s.monitor.SentMessage()
	}
	return err
}

func (s *monitoredClientStream) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if err == nil {
		s.monitor.ReceivedMessage()
	} else if err == io.EOF {
		s.monitor.Handled(codes.OK)
	} else {
		s.monitor.Handled(grpc.Code(err))
	}
	return err
}
