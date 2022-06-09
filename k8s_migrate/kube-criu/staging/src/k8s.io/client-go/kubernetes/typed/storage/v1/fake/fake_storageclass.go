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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	storagev1 "k8s.io/api/storage/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeStorageClasses implements StorageClassInterface
type FakeStorageClasses struct {
	Fake *FakeStorageV1
}

var storageclassesResource = schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"}

var storageclassesKind = schema.GroupVersionKind{Group: "storage.k8s.io", Version: "v1", Kind: "StorageClass"}

// Get takes name of the storageClass, and returns the corresponding storageClass object, and an error if there is any.
func (c *FakeStorageClasses) Get(name string, options v1.GetOptions) (result *storagev1.StorageClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(storageclassesResource, name), &storagev1.StorageClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*storagev1.StorageClass), err
}

// List takes label and field selectors, and returns the list of StorageClasses that match those selectors.
func (c *FakeStorageClasses) List(opts v1.ListOptions) (result *storagev1.StorageClassList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(storageclassesResource, storageclassesKind, opts), &storagev1.StorageClassList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &storagev1.StorageClassList{ListMeta: obj.(*storagev1.StorageClassList).ListMeta}
	for _, item := range obj.(*storagev1.StorageClassList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested storageClasses.
func (c *FakeStorageClasses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(storageclassesResource, opts))
}

// Create takes the representation of a storageClass and creates it.  Returns the server's representation of the storageClass, and an error, if there is any.
func (c *FakeStorageClasses) Create(storageClass *storagev1.StorageClass) (result *storagev1.StorageClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(storageclassesResource, storageClass), &storagev1.StorageClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*storagev1.StorageClass), err
}

// Update takes the representation of a storageClass and updates it. Returns the server's representation of the storageClass, and an error, if there is any.
func (c *FakeStorageClasses) Update(storageClass *storagev1.StorageClass) (result *storagev1.StorageClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(storageclassesResource, storageClass), &storagev1.StorageClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*storagev1.StorageClass), err
}

// Delete takes name of the storageClass and deletes it. Returns an error if one occurs.
func (c *FakeStorageClasses) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(storageclassesResource, name), &storagev1.StorageClass{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeStorageClasses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(storageclassesResource, listOptions)

	_, err := c.Fake.Invokes(action, &storagev1.StorageClassList{})
	return err
}

// Patch applies the patch and returns the patched storageClass.
func (c *FakeStorageClasses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *storagev1.StorageClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(storageclassesResource, name, data, subresources...), &storagev1.StorageClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*storagev1.StorageClass), err
}
