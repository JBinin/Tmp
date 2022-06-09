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

// gRPC Prometheus monitoring interceptors for server-side gRPC.

package grpc_prometheus

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// PreregisterServices takes a gRPC server and pre-initializes all counters to 0.
// This allows for easier monitoring in Prometheus (no missing metrics), and should be called *after* all services have
// been registered with the server.
func Register(server *grpc.Server) {
	serviceInfo := server.GetServiceInfo()
	for serviceName, info := range serviceInfo {
		for _, mInfo := range info.Methods {
			preRegisterMethod(serviceName, &mInfo)
		}
	}
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	monitor := newServerReporter(Unary, info.FullMethod)
	monitor.ReceivedMessage()
	resp, err := handler(ctx, req)
	monitor.Handled(grpc.Code(err))
	if err == nil {
		monitor.SentMessage()
	}
	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	monitor := newServerReporter(streamRpcType(info), info.FullMethod)
	err := handler(srv, &monitoredServerStream{ss, monitor})
	monitor.Handled(grpc.Code(err))
	return err
}

func streamRpcType(info *grpc.StreamServerInfo) grpcType {
	if info.IsClientStream && !info.IsServerStream {
		return ClientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return ServerStream
	}
	return BidiStream
}

// monitoredStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type monitoredServerStream struct {
	grpc.ServerStream
	monitor *serverReporter
}

func (s *monitoredServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.monitor.SentMessage()
	}
	return err
}

func (s *monitoredServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.monitor.ReceivedMessage()
	}
	return err
}
