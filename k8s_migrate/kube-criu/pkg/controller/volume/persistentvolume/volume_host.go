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

package persistentvolume

import (
	"fmt"
	"net"

	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/util/mount"
	vol "k8s.io/kubernetes/pkg/volume"
)

// VolumeHost interface implementation for PersistentVolumeController.

var _ vol.VolumeHost = &PersistentVolumeController{}

func (ctrl *PersistentVolumeController) GetPluginDir(pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetVolumeDevicePluginDir(pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodsDir() string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodVolumeDir(podUID types.UID, pluginName string, volumeName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodPluginDir(podUID types.UID, pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodVolumeDeviceDir(ppodUID types.UID, pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetKubeClient() clientset.Interface {
	return ctrl.kubeClient
}

func (ctrl *PersistentVolumeController) NewWrapperMounter(volName string, spec vol.Spec, pod *v1.Pod, opts vol.VolumeOptions) (vol.Mounter, error) {
	return nil, fmt.Errorf("PersistentVolumeController.NewWrapperMounter is not implemented")
}

func (ctrl *PersistentVolumeController) NewWrapperUnmounter(volName string, spec vol.Spec, podUID types.UID) (vol.Unmounter, error) {
	return nil, fmt.Errorf("PersistentVolumeController.NewWrapperMounter is not implemented")
}

func (ctrl *PersistentVolumeController) GetCloudProvider() cloudprovider.Interface {
	return ctrl.cloud
}

func (ctrl *PersistentVolumeController) GetMounter(pluginName string) mount.Interface {
	return nil
}

func (ctrl *PersistentVolumeController) GetHostName() string {
	return ""
}

func (ctrl *PersistentVolumeController) GetHostIP() (net.IP, error) {
	return nil, fmt.Errorf("PersistentVolumeController.GetHostIP() is not implemented")
}

func (ctrl *PersistentVolumeController) GetNodeAllocatable() (v1.ResourceList, error) {
	return v1.ResourceList{}, nil
}

func (ctrl *PersistentVolumeController) GetSecretFunc() func(namespace, name string) (*v1.Secret, error) {
	return func(_, _ string) (*v1.Secret, error) {
		return nil, fmt.Errorf("GetSecret unsupported in PersistentVolumeController")
	}
}

func (ctrl *PersistentVolumeController) GetConfigMapFunc() func(namespace, name string) (*v1.ConfigMap, error) {
	return func(_, _ string) (*v1.ConfigMap, error) {
		return nil, fmt.Errorf("GetConfigMap unsupported in PersistentVolumeController")
	}
}

func (ctrl *PersistentVolumeController) GetServiceAccountTokenFunc() func(_, _ string, _ *authenticationv1.TokenRequest) (*authenticationv1.TokenRequest, error) {
	return func(_, _ string, _ *authenticationv1.TokenRequest) (*authenticationv1.TokenRequest, error) {
		return nil, fmt.Errorf("GetServiceAccountToken unsupported in PersistentVolumeController")
	}
}

func (adc *PersistentVolumeController) GetExec(pluginName string) mount.Exec {
	return mount.NewOsExec()
}

func (ctrl *PersistentVolumeController) GetNodeLabels() (map[string]string, error) {
	return nil, fmt.Errorf("GetNodeLabels() unsupported in PersistentVolumeController")
}

func (ctrl *PersistentVolumeController) GetNodeName() types.NodeName {
	return ""
}

func (ctrl *PersistentVolumeController) GetEventRecorder() record.EventRecorder {
	return ctrl.eventRecorder
}
