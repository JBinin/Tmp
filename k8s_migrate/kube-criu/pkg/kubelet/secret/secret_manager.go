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

package secret

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	corev1 "k8s.io/kubernetes/pkg/apis/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/util/manager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
)

type Manager interface {
	// Get secret by secret namespace and name.
	GetSecret(namespace, name string) (*v1.Secret, error)

	// WARNING: Register/UnregisterPod functions should be efficient,
	// i.e. should not block on network operations.

	// RegisterPod registers all secrets from a given pod.
	RegisterPod(pod *v1.Pod)

	// UnregisterPod unregisters secrets from a given pod that are not
	// used by any other registered pod.
	UnregisterPod(pod *v1.Pod)
}

// simpleSecretManager implements SecretManager interfaces with
// simple operations to apiserver.
type simpleSecretManager struct {
	kubeClient clientset.Interface
}

func NewSimpleSecretManager(kubeClient clientset.Interface) Manager {
	return &simpleSecretManager{kubeClient: kubeClient}
}

func (s *simpleSecretManager) GetSecret(namespace, name string) (*v1.Secret, error) {
	return s.kubeClient.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
}

func (s *simpleSecretManager) RegisterPod(pod *v1.Pod) {
}

func (s *simpleSecretManager) UnregisterPod(pod *v1.Pod) {
}

// secretManager keeps a store with secrets necessary
// for registered pods. Different implementations of the store
// may result in different semantics for freshness of secrets
// (e.g. ttl-based implementation vs watch-based implementation).
type secretManager struct {
	manager manager.Manager
}

func (s *secretManager) GetSecret(namespace, name string) (*v1.Secret, error) {
	object, err := s.manager.GetObject(namespace, name)
	if err != nil {
		return nil, err
	}
	if secret, ok := object.(*v1.Secret); ok {
		return secret, nil
	}
	return nil, fmt.Errorf("unexpected object type: %v", object)
}

func (s *secretManager) RegisterPod(pod *v1.Pod) {
	s.manager.RegisterPod(pod)
}

func (s *secretManager) UnregisterPod(pod *v1.Pod) {
	s.manager.UnregisterPod(pod)
}

func getSecretNames(pod *v1.Pod) sets.String {
	result := sets.NewString()
	podutil.VisitPodSecretNames(pod, func(name string) bool {
		result.Insert(name)
		return true
	})
	return result
}

const (
	defaultTTL = time.Minute
)

// NewCachingSecretManager creates a manager that keeps a cache of all secrets
// necessary for registered pods.
// It implements the following logic:
// - whenever a pod is created or updated, the cached versions of all secrets
//   are invalidated
// - every GetObject() call tries to fetch the value from local cache; if it is
//   not there, invalidated or too old, we fetch it from apiserver and refresh the
//   value in cache; otherwise it is just fetched from cache
func NewCachingSecretManager(kubeClient clientset.Interface, getTTL manager.GetObjectTTLFunc) Manager {
	getSecret := func(namespace, name string, opts metav1.GetOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Secrets(namespace).Get(name, opts)
	}
	secretStore := manager.NewObjectStore(getSecret, clock.RealClock{}, getTTL, defaultTTL)
	return &secretManager{
		manager: manager.NewCacheBasedManager(secretStore, getSecretNames),
	}
}

// NewWatchingSecretManager creates a manager that keeps a cache of all secrets
// necessary for registered pods.
// It implements the following logic:
// - whenever a pod is created or updated, we start inidvidual watches for all
//   referenced objects that aren't referenced from other registered pods
// - every GetObject() returns a value from local cache propagated via watches
func NewWatchingSecretManager(kubeClient clientset.Interface) Manager {
	listSecret := func(namespace string, opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Secrets(namespace).List(opts)
	}
	watchSecret := func(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.CoreV1().Secrets(namespace).Watch(opts)
	}
	newSecret := func() runtime.Object {
		return &v1.Secret{}
	}
	gr := corev1.Resource("secret")
	return &secretManager{
		manager: manager.NewWatchBasedManager(listSecret, watchSecret, newSecret, gr, getSecretNames),
	}
}
