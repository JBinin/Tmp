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

package v1beta1

import (
	v1beta1 "k8s.io/api/coordination/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// LeaseLister helps list Leases.
type LeaseLister interface {
	// List lists all Leases in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.Lease, err error)
	// Leases returns an object that can list and get Leases.
	Leases(namespace string) LeaseNamespaceLister
	LeaseListerExpansion
}

// leaseLister implements the LeaseLister interface.
type leaseLister struct {
	indexer cache.Indexer
}

// NewLeaseLister returns a new LeaseLister.
func NewLeaseLister(indexer cache.Indexer) LeaseLister {
	return &leaseLister{indexer: indexer}
}

// List lists all Leases in the indexer.
func (s *leaseLister) List(selector labels.Selector) (ret []*v1beta1.Lease, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.Lease))
	})
	return ret, err
}

// Leases returns an object that can list and get Leases.
func (s *leaseLister) Leases(namespace string) LeaseNamespaceLister {
	return leaseNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// LeaseNamespaceLister helps list and get Leases.
type LeaseNamespaceLister interface {
	// List lists all Leases in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.Lease, err error)
	// Get retrieves the Lease from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.Lease, error)
	LeaseNamespaceListerExpansion
}

// leaseNamespaceLister implements the LeaseNamespaceLister
// interface.
type leaseNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Leases in the indexer for a given namespace.
func (s leaseNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.Lease, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.Lease))
	})
	return ret, err
}

// Get retrieves the Lease from the indexer for a given namespace and name.
func (s leaseNamespaceLister) Get(name string) (*v1beta1.Lease, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("lease"), name)
	}
	return obj.(*v1beta1.Lease), nil
}
