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

package polymorphichelpers

import (
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func portsForObject(object runtime.Object) ([]string, error) {
	switch t := object.(type) {
	case *api.ReplicationController:
		return getPortsInternal(t.Spec.Template.Spec), nil
	case *corev1.ReplicationController:
		return getPorts(t.Spec.Template.Spec), nil

	case *api.Pod:
		return getPortsInternal(t.Spec), nil
	case *corev1.Pod:
		return getPorts(t.Spec), nil

	case *api.Service:
		return getServicePortsInternal(t.Spec), nil
	case *corev1.Service:
		return getServicePorts(t.Spec), nil

	case *extensions.Deployment:
		return getPortsInternal(t.Spec.Template.Spec), nil
	case *extensionsv1beta1.Deployment:
		return getPorts(t.Spec.Template.Spec), nil
	case *appsv1.Deployment:
		return getPorts(t.Spec.Template.Spec), nil
	case *appsv1beta2.Deployment:
		return getPorts(t.Spec.Template.Spec), nil
	case *appsv1beta1.Deployment:
		return getPorts(t.Spec.Template.Spec), nil

	case *extensions.ReplicaSet:
		return getPortsInternal(t.Spec.Template.Spec), nil
	case *extensionsv1beta1.ReplicaSet:
		return getPorts(t.Spec.Template.Spec), nil
	case *appsv1.ReplicaSet:
		return getPorts(t.Spec.Template.Spec), nil
	case *appsv1beta2.ReplicaSet:
		return getPorts(t.Spec.Template.Spec), nil
	default:
		return nil, fmt.Errorf("cannot extract ports from %T", object)
	}
}

func getPortsInternal(spec api.PodSpec) []string {
	result := []string{}
	for _, container := range spec.Containers {
		for _, port := range container.Ports {
			result = append(result, strconv.Itoa(int(port.ContainerPort)))
		}
	}
	return result
}

// Extracts the ports exposed by a service from the given service spec.
func getServicePortsInternal(spec api.ServiceSpec) []string {
	result := []string{}
	for _, servicePort := range spec.Ports {
		result = append(result, strconv.Itoa(int(servicePort.Port)))
	}
	return result
}

func getPorts(spec corev1.PodSpec) []string {
	result := []string{}
	for _, container := range spec.Containers {
		for _, port := range container.Ports {
			result = append(result, strconv.Itoa(int(port.ContainerPort)))
		}
	}
	return result
}

func getServicePorts(spec corev1.ServiceSpec) []string {
	result := []string{}
	for _, servicePort := range spec.Ports {
		result = append(result, strconv.Itoa(int(servicePort.Port)))
	}
	return result
}
