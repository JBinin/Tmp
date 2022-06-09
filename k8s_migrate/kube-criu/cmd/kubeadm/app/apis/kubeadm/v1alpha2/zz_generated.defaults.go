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
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by defaulter-gen. DO NOT EDIT.

package v1alpha2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	v1beta1 "k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig/v1beta1"
	v1alpha1 "k8s.io/kubernetes/pkg/proxy/apis/kubeproxyconfig/v1alpha1"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&InitConfiguration{}, func(obj interface{}) { SetObjectDefaults_InitConfiguration(obj.(*InitConfiguration)) })
	scheme.AddTypeDefaultingFunc(&JoinConfiguration{}, func(obj interface{}) { SetObjectDefaults_JoinConfiguration(obj.(*JoinConfiguration)) })
	return nil
}

func SetObjectDefaults_InitConfiguration(in *InitConfiguration) {
	SetDefaults_InitConfiguration(in)
	for i := range in.BootstrapTokens {
		a := &in.BootstrapTokens[i]
		SetDefaults_BootstrapToken(a)
	}
	SetDefaults_NodeRegistrationOptions(&in.NodeRegistration)
	if in.KubeProxy.Config != nil {
		v1alpha1.SetDefaults_KubeProxyConfiguration(in.KubeProxy.Config)
	}
	if in.KubeletConfiguration.BaseConfig != nil {
		v1beta1.SetDefaults_KubeletConfiguration(in.KubeletConfiguration.BaseConfig)
	}
}

func SetObjectDefaults_JoinConfiguration(in *JoinConfiguration) {
	SetDefaults_JoinConfiguration(in)
	SetDefaults_NodeRegistrationOptions(&in.NodeRegistration)
}
