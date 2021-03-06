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
// +build windows

/*
Copyright 2015 The Kubernetes Authors.

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

package cadvisor

import (
	"github.com/google/cadvisor/events"
	cadvisorapi "github.com/google/cadvisor/info/v1"
	cadvisorapiv2 "github.com/google/cadvisor/info/v2"
	"k8s.io/kubernetes/pkg/kubelet/winstats"
)

type cadvisorClient struct {
	rootPath       string
	winStatsClient winstats.Client
}

var _ Interface = new(cadvisorClient)

// New creates a cAdvisor and exports its API on the specified port if port > 0.
func New(imageFsInfoProvider ImageFsInfoProvider, rootPath string, usingLegacyStats bool) (Interface, error) {
	client, err := winstats.NewPerfCounterClient()
	return &cadvisorClient{
		rootPath:       rootPath,
		winStatsClient: client,
	}, err
}

func (cu *cadvisorClient) Start() error {
	return nil
}

func (cu *cadvisorClient) DockerContainer(name string, req *cadvisorapi.ContainerInfoRequest) (cadvisorapi.ContainerInfo, error) {
	return cadvisorapi.ContainerInfo{}, nil
}

func (cu *cadvisorClient) ContainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (*cadvisorapi.ContainerInfo, error) {
	return &cadvisorapi.ContainerInfo{}, nil
}

func (cu *cadvisorClient) ContainerInfoV2(name string, options cadvisorapiv2.RequestOptions) (map[string]cadvisorapiv2.ContainerInfo, error) {
	return cu.winStatsClient.WinContainerInfos()
}

func (cu *cadvisorClient) SubcontainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (map[string]*cadvisorapi.ContainerInfo, error) {
	return nil, nil
}

func (cu *cadvisorClient) MachineInfo() (*cadvisorapi.MachineInfo, error) {
	return cu.winStatsClient.WinMachineInfo()
}

func (cu *cadvisorClient) VersionInfo() (*cadvisorapi.VersionInfo, error) {
	return cu.winStatsClient.WinVersionInfo()
}

func (cu *cadvisorClient) ImagesFsInfo() (cadvisorapiv2.FsInfo, error) {
	return cadvisorapiv2.FsInfo{}, nil
}

func (cu *cadvisorClient) RootFsInfo() (cadvisorapiv2.FsInfo, error) {
	return cu.GetDirFsInfo(cu.rootPath)
}

func (cu *cadvisorClient) WatchEvents(request *events.Request) (*events.EventChannel, error) {
	return &events.EventChannel{}, nil
}

func (cu *cadvisorClient) GetDirFsInfo(path string) (cadvisorapiv2.FsInfo, error) {
	return cu.winStatsClient.GetDirFsInfo(path)
}
