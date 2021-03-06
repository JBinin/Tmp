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
Copyright 2015 The Kubernetes Authors.

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

package framework

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

func NewFakeControllerSource() *FakeControllerSource {
	return &FakeControllerSource{
		Items:       map[nnu]runtime.Object{},
		Broadcaster: watch.NewBroadcaster(100, watch.WaitIfChannelFull),
	}
}

func NewFakePVControllerSource() *FakePVControllerSource {
	return &FakePVControllerSource{
		FakeControllerSource{
			Items:       map[nnu]runtime.Object{},
			Broadcaster: watch.NewBroadcaster(100, watch.WaitIfChannelFull),
		}}
}

func NewFakePVCControllerSource() *FakePVCControllerSource {
	return &FakePVCControllerSource{
		FakeControllerSource{
			Items:       map[nnu]runtime.Object{},
			Broadcaster: watch.NewBroadcaster(100, watch.WaitIfChannelFull),
		}}
}

// FakeControllerSource implements listing/watching for testing.
type FakeControllerSource struct {
	lock        sync.RWMutex
	Items       map[nnu]runtime.Object
	changes     []watch.Event // one change per resourceVersion
	Broadcaster *watch.Broadcaster
}

type FakePVControllerSource struct {
	FakeControllerSource
}

type FakePVCControllerSource struct {
	FakeControllerSource
}

// namespace, name, uid to be used as a key.
type nnu struct {
	namespace, name string
	uid             types.UID
}

// Add adds an object to the set and sends an add event to watchers.
// obj's ResourceVersion is set.
func (f *FakeControllerSource) Add(obj runtime.Object) {
	f.Change(watch.Event{Type: watch.Added, Object: obj}, 1)
}

// Modify updates an object in the set and sends a modified event to watchers.
// obj's ResourceVersion is set.
func (f *FakeControllerSource) Modify(obj runtime.Object) {
	f.Change(watch.Event{Type: watch.Modified, Object: obj}, 1)
}

// Delete deletes an object from the set and sends a delete event to watchers.
// obj's ResourceVersion is set.
func (f *FakeControllerSource) Delete(lastValue runtime.Object) {
	f.Change(watch.Event{Type: watch.Deleted, Object: lastValue}, 1)
}

// AddDropWatch adds an object to the set but forgets to send an add event to
// watchers.
// obj's ResourceVersion is set.
func (f *FakeControllerSource) AddDropWatch(obj runtime.Object) {
	f.Change(watch.Event{Type: watch.Added, Object: obj}, 0)
}

// ModifyDropWatch updates an object in the set but forgets to send a modify
// event to watchers.
// obj's ResourceVersion is set.
func (f *FakeControllerSource) ModifyDropWatch(obj runtime.Object) {
	f.Change(watch.Event{Type: watch.Modified, Object: obj}, 0)
}

// DeleteDropWatch deletes an object from the set but forgets to send a delete
// event to watchers.
// obj's ResourceVersion is set.
func (f *FakeControllerSource) DeleteDropWatch(lastValue runtime.Object) {
	f.Change(watch.Event{Type: watch.Deleted, Object: lastValue}, 0)
}

func (f *FakeControllerSource) key(accessor metav1.Object) nnu {
	return nnu{accessor.GetNamespace(), accessor.GetName(), accessor.GetUID()}
}

// Change records the given event (setting the object's resource version) and
// sends a watch event with the specified probability.
func (f *FakeControllerSource) Change(e watch.Event, watchProbability float64) {
	f.lock.Lock()
	defer f.lock.Unlock()

	accessor, err := meta.Accessor(e.Object)
	if err != nil {
		panic(err) // this is test code only
	}

	resourceVersion := len(f.changes) + 1
	accessor.SetResourceVersion(strconv.Itoa(resourceVersion))
	f.changes = append(f.changes, e)
	key := f.key(accessor)
	switch e.Type {
	case watch.Added, watch.Modified:
		f.Items[key] = e.Object
	case watch.Deleted:
		delete(f.Items, key)
	}

	if rand.Float64() < watchProbability {
		f.Broadcaster.Action(e.Type, e.Object)
	}
}

func (f *FakeControllerSource) getListItemsLocked() ([]runtime.Object, error) {
	list := make([]runtime.Object, 0, len(f.Items))
	for _, obj := range f.Items {
		// Must make a copy to allow clients to modify the object.
		// Otherwise, if they make a change and write it back, they
		// will inadvertently change our canonical copy (in
		// addition to racing with other clients).
		list = append(list, obj.DeepCopyObject())
	}
	return list, nil
}

// List returns a list object, with its resource version set.
func (f *FakeControllerSource) List(options metav1.ListOptions) (runtime.Object, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	list, err := f.getListItemsLocked()
	if err != nil {
		return nil, err
	}
	listObj := &v1.List{}
	if err := meta.SetList(listObj, list); err != nil {
		return nil, err
	}
	listAccessor, err := meta.ListAccessor(listObj)
	if err != nil {
		return nil, err
	}
	resourceVersion := len(f.changes)
	listAccessor.SetResourceVersion(strconv.Itoa(resourceVersion))
	return listObj, nil
}

// List returns a list object, with its resource version set.
func (f *FakePVControllerSource) List(options metav1.ListOptions) (runtime.Object, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	list, err := f.FakeControllerSource.getListItemsLocked()
	if err != nil {
		return nil, err
	}
	listObj := &v1.PersistentVolumeList{}
	if err := meta.SetList(listObj, list); err != nil {
		return nil, err
	}
	listAccessor, err := meta.ListAccessor(listObj)
	if err != nil {
		return nil, err
	}
	resourceVersion := len(f.changes)
	listAccessor.SetResourceVersion(strconv.Itoa(resourceVersion))
	return listObj, nil
}

// List returns a list object, with its resource version set.
func (f *FakePVCControllerSource) List(options metav1.ListOptions) (runtime.Object, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	list, err := f.FakeControllerSource.getListItemsLocked()
	if err != nil {
		return nil, err
	}
	listObj := &v1.PersistentVolumeClaimList{}
	if err := meta.SetList(listObj, list); err != nil {
		return nil, err
	}
	listAccessor, err := meta.ListAccessor(listObj)
	if err != nil {
		return nil, err
	}
	resourceVersion := len(f.changes)
	listAccessor.SetResourceVersion(strconv.Itoa(resourceVersion))
	return listObj, nil
}

// Watch returns a watch, which will be pre-populated with all changes
// after resourceVersion.
func (f *FakeControllerSource) Watch(options metav1.ListOptions) (watch.Interface, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	rc, err := strconv.Atoi(options.ResourceVersion)
	if err != nil {
		return nil, err
	}
	if rc < len(f.changes) {
		changes := []watch.Event{}
		for _, c := range f.changes[rc:] {
			// Must make a copy to allow clients to modify the
			// object.  Otherwise, if they make a change and write
			// it back, they will inadvertently change the our
			// canonical copy (in addition to racing with other
			// clients).
			changes = append(changes, watch.Event{Type: c.Type, Object: c.Object.DeepCopyObject()})
		}
		return f.Broadcaster.WatchWithPrefix(changes), nil
	} else if rc > len(f.changes) {
		return nil, errors.New("resource version in the future not supported by this fake")
	}
	return f.Broadcaster.Watch(), nil
}

// Shutdown closes the underlying broadcaster, waiting for events to be
// delivered. It's an error to call any method after calling shutdown. This is
// enforced by Shutdown() leaving f locked.
func (f *FakeControllerSource) Shutdown() {
	f.lock.Lock() // Purposely no unlock.
	f.Broadcaster.Shutdown()
}
