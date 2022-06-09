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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "k8s.io/podcheckpoint/pkg/apis/podcheckpointcontroller/v1alpha1"
)

// PodCheckpointLister helps list PodCheckpoints.
type PodCheckpointLister interface {
	// List lists all PodCheckpoints in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.PodCheckpoint, err error)
	// PodCheckpoints returns an object that can list and get PodCheckpoints.
	PodCheckpoints(namespace string) PodCheckpointNamespaceLister
	PodCheckpointListerExpansion
}

// podCheckpointLister implements the PodCheckpointLister interface.
type podCheckpointLister struct {
	indexer cache.Indexer
}

// NewPodCheckpointLister returns a new PodCheckpointLister.
func NewPodCheckpointLister(indexer cache.Indexer) PodCheckpointLister {
	return &podCheckpointLister{indexer: indexer}
}

// List lists all PodCheckpoints in the indexer.
func (s *podCheckpointLister) List(selector labels.Selector) (ret []*v1alpha1.PodCheckpoint, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PodCheckpoint))
	})
	return ret, err
}

// PodCheckpoints returns an object that can list and get PodCheckpoints.
func (s *podCheckpointLister) PodCheckpoints(namespace string) PodCheckpointNamespaceLister {
	return podCheckpointNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PodCheckpointNamespaceLister helps list and get PodCheckpoints.
type PodCheckpointNamespaceLister interface {
	// List lists all PodCheckpoints in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.PodCheckpoint, err error)
	// Get retrieves the PodCheckpoint from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.PodCheckpoint, error)
	PodCheckpointNamespaceListerExpansion
}

// podCheckpointNamespaceLister implements the PodCheckpointNamespaceLister
// interface.
type podCheckpointNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PodCheckpoints in the indexer for a given namespace.
func (s podCheckpointNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.PodCheckpoint, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PodCheckpoint))
	})
	return ret, err
}

// Get retrieves the PodCheckpoint from the indexer for a given namespace and name.
func (s podCheckpointNamespaceLister) Get(name string) (*v1alpha1.PodCheckpoint, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("podcheckpoint"), name)
	}
	return obj.(*v1alpha1.PodCheckpoint), nil
}