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
Copyright 2018 The Kubernetes Authors.

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

package componentconfigs

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeletvalidation "k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig/validation"
	proxyvalidation "k8s.io/kubernetes/pkg/proxy/apis/kubeproxyconfig/validation"
)

// ValidateKubeProxyConfiguration validates proxy configuration and collects all encountered errors
func ValidateKubeProxyConfiguration(internalcfg *kubeadmapi.InitConfiguration, _ *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if internalcfg.ComponentConfigs.KubeProxy == nil {
		return allErrs
	}
	return proxyvalidation.Validate(internalcfg.ComponentConfigs.KubeProxy)
}

// ValidateKubeletConfiguration validates kubelet configuration and collects all encountered errors
func ValidateKubeletConfiguration(internalcfg *kubeadmapi.InitConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if internalcfg.ComponentConfigs.Kubelet == nil {
		return allErrs
	}

	if err := kubeletvalidation.ValidateKubeletConfiguration(internalcfg.ComponentConfigs.Kubelet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, "", err.Error()))
	}
	return allErrs
}
