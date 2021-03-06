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
Copyright 2016 The Kubernetes Authors.

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

package remote

import (
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/dockershim"
	"k8s.io/kubernetes/pkg/kubelet/util"
)

// maxMsgSize use 8MB as the default message size limit.
// grpc library default is 4MB
const maxMsgSize = 1024 * 1024 * 8

// DockerServer is the grpc server of dockershim.
type DockerServer struct {
	// endpoint is the endpoint to serve on.
	endpoint string
	// service is the docker service which implements runtime and image services.
	service dockershim.CRIService
	// server is the grpc server.
	server *grpc.Server
}

// NewDockerServer creates the dockershim grpc server.
func NewDockerServer(endpoint string, s dockershim.CRIService) *DockerServer {
	return &DockerServer{
		endpoint: endpoint,
		service:  s,
	}
}

// Start starts the dockershim grpc server.
func (s *DockerServer) Start() error {
	// Start the internal service.
	if err := s.service.Start(); err != nil {
		glog.Errorf("Unable to start docker service")
		return err
	}

	glog.V(2).Infof("Start dockershim grpc server")
	l, err := util.CreateListener(s.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen on %q: %v", s.endpoint, err)
	}
	// Create the grpc server and register runtime and image services.
	s.server = grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	runtimeapi.RegisterRuntimeServiceServer(s.server, s.service)
	runtimeapi.RegisterImageServiceServer(s.server, s.service)
	go func() {
		if err := s.server.Serve(l); err != nil {
			glog.Fatalf("Failed to serve connections: %v", err)
		}
	}()
	return nil
}
