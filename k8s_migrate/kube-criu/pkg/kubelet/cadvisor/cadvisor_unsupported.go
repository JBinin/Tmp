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
// +build !linux,!windows linux,!cgo

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
	"errors"

	"github.com/google/cadvisor/events"
	cadvisorapi "github.com/google/cadvisor/info/v1"
	cadvisorapiv2 "github.com/google/cadvisor/info/v2"
)

type cadvisorUnsupported struct {
}

var _ Interface = new(cadvisorUnsupported)

func New(imageFsInfoProvider ImageFsInfoProvider, rootPath string, usingLegacyStats bool) (Interface, error) {
	return &cadvisorUnsupported{}, nil
}

var unsupportedErr = errors.New("cAdvisor is unsupported in this build")

func (cu *cadvisorUnsupported) Start() error {
	return unsupportedErr
}

func (cu *cadvisorUnsupported) DockerContainer(name string, req *cadvisorapi.ContainerInfoRequest) (cadvisorapi.ContainerInfo, error) {
	return cadvisorapi.ContainerInfo{}, unsupportedErr
}

func (cu *cadvisorUnsupported) ContainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (*cadvisorapi.ContainerInfo, error) {
	return nil, unsupportedErr
}

func (cu *cadvisorUnsupported) ContainerInfoV2(name string, options cadvisorapiv2.RequestOptions) (map[string]cadvisorapiv2.ContainerInfo, error) {
	return nil, unsupportedErr
}

func (cu *cadvisorUnsupported) SubcontainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (map[string]*cadvisorapi.ContainerInfo, error) {
	return nil, unsupportedErr
}

func (cu *cadvisorUnsupported) MachineInfo() (*cadvisorapi.MachineInfo, error) {
	return nil, unsupportedErr
}

func (cu *cadvisorUnsupported) VersionInfo() (*cadvisorapi.VersionInfo, error) {
	return nil, unsupportedErr
}

func (cu *cadvisorUnsupported) ImagesFsInfo() (cadvisorapiv2.FsInfo, error) {
	return cadvisorapiv2.FsInfo{}, unsupportedErr
}

func (cu *cadvisorUnsupported) RootFsInfo() (cadvisorapiv2.FsInfo, error) {
	return cadvisorapiv2.FsInfo{}, unsupportedErr
}

func (cu *cadvisorUnsupported) WatchEvents(request *events.Request) (*events.EventChannel, error) {
	return nil, unsupportedErr
}

func (cu *cadvisorUnsupported) GetDirFsInfo(path string) (cadvisorapiv2.FsInfo, error) {
	return cadvisorapiv2.FsInfo{}, nil
}
