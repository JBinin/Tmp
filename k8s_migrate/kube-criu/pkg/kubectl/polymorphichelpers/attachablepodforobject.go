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
	"sort"
	"time"

	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
	apiv1 "k8s.io/kubernetes/pkg/apis/core/v1"
	"k8s.io/kubernetes/pkg/controller"
)

// attachablePodForObject returns the pod to which to attach given an object.
func attachablePodForObject(restClientGetter genericclioptions.RESTClientGetter, object runtime.Object, timeout time.Duration) (*corev1.Pod, error) {
	switch t := object.(type) {
	case *api.Pod:
		externalPod := &corev1.Pod{}
		err := apiv1.Convert_core_Pod_To_v1_Pod(t, externalPod, nil)
		return externalPod, err
	case *corev1.Pod:
		return t, nil
	}

	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := corev1client.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	namespace, selector, err := SelectorsForObject(object)
	if err != nil {
		return nil, fmt.Errorf("cannot attach to %T: %v", object, err)
	}
	sortBy := func(pods []*v1.Pod) sort.Interface { return sort.Reverse(controller.ActivePods(pods)) }
	pod, _, err := GetFirstPod(clientset, namespace, selector.String(), timeout, sortBy)
	return pod, err
}
