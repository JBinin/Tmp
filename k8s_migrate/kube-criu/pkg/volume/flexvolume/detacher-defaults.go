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
Copyright 2017 The Kubernetes Authors.

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

package flexvolume

import (
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/volume/util"
)

type detacherDefaults flexVolumeDetacher

// Detach is part of the volume.Detacher interface.
func (d *detacherDefaults) Detach(volumeName string, hostName types.NodeName) error {
	glog.Warning(logPrefix(d.plugin.flexVolumePlugin), "using default Detach for volume ", volumeName, ", host ", hostName)
	return nil
}

// WaitForDetach is part of the volume.Detacher interface.
func (d *detacherDefaults) WaitForDetach(devicePath string, timeout time.Duration) error {
	glog.Warning(logPrefix(d.plugin.flexVolumePlugin), "using default WaitForDetach for device ", devicePath)
	return nil
}

// UnmountDevice is part of the volume.Detacher interface.
func (d *detacherDefaults) UnmountDevice(deviceMountPath string) error {
	glog.Warning(logPrefix(d.plugin.flexVolumePlugin), "using default UnmountDevice for device mount path ", deviceMountPath)
	return util.UnmountPath(deviceMountPath, d.plugin.host.GetMounter(d.plugin.GetPluginName()))
}
