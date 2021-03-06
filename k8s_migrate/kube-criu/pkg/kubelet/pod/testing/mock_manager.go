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
Copyright 2017 The Kubernetes Authors.

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

// Code generated by mockery v1.0.0
package testing

import kubelettypes "k8s.io/kubernetes/pkg/kubelet/types"
import mock "github.com/stretchr/testify/mock"

import types "k8s.io/apimachinery/pkg/types"
import v1 "k8s.io/api/core/v1"

// MockManager is an autogenerated mock type for the Manager type
type MockManager struct {
	mock.Mock
}

// AddPod provides a mock function with given fields: _a0
func (_m *MockManager) AddPod(_a0 *v1.Pod) {
	_m.Called(_a0)
}

// CreateMirrorPod provides a mock function with given fields: _a0
func (_m *MockManager) CreateMirrorPod(_a0 *v1.Pod) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.Pod) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMirrorPod provides a mock function with given fields: podFullName
func (_m *MockManager) DeleteMirrorPod(podFullName string) error {
	ret := _m.Called(podFullName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(podFullName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOrphanedMirrorPods provides a mock function with given fields:
func (_m *MockManager) DeleteOrphanedMirrorPods() {
	_m.Called()
}

// DeletePod provides a mock function with given fields: _a0
func (_m *MockManager) DeletePod(_a0 *v1.Pod) {
	_m.Called(_a0)
}

// GetMirrorPodByPod provides a mock function with given fields: _a0
func (_m *MockManager) GetMirrorPodByPod(_a0 *v1.Pod) (*v1.Pod, bool) {
	ret := _m.Called(_a0)

	var r0 *v1.Pod
	if rf, ok := ret.Get(0).(func(*v1.Pod) *v1.Pod); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Pod)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(*v1.Pod) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetPodByFullName provides a mock function with given fields: podFullName
func (_m *MockManager) GetPodByFullName(podFullName string) (*v1.Pod, bool) {
	ret := _m.Called(podFullName)

	var r0 *v1.Pod
	if rf, ok := ret.Get(0).(func(string) *v1.Pod); ok {
		r0 = rf(podFullName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Pod)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(podFullName)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetPodByMirrorPod provides a mock function with given fields: _a0
func (_m *MockManager) GetPodByMirrorPod(_a0 *v1.Pod) (*v1.Pod, bool) {
	ret := _m.Called(_a0)

	var r0 *v1.Pod
	if rf, ok := ret.Get(0).(func(*v1.Pod) *v1.Pod); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Pod)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(*v1.Pod) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetPodByName provides a mock function with given fields: namespace, name
func (_m *MockManager) GetPodByName(namespace string, name string) (*v1.Pod, bool) {
	ret := _m.Called(namespace, name)

	var r0 *v1.Pod
	if rf, ok := ret.Get(0).(func(string, string) *v1.Pod); ok {
		r0 = rf(namespace, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Pod)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string, string) bool); ok {
		r1 = rf(namespace, name)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetPodByUID provides a mock function with given fields: _a0
func (_m *MockManager) GetPodByUID(_a0 types.UID) (*v1.Pod, bool) {
	ret := _m.Called(_a0)

	var r0 *v1.Pod
	if rf, ok := ret.Get(0).(func(types.UID) *v1.Pod); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Pod)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(types.UID) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetPods provides a mock function with given fields:
func (_m *MockManager) GetPods() []*v1.Pod {
	ret := _m.Called()

	var r0 []*v1.Pod
	if rf, ok := ret.Get(0).(func() []*v1.Pod); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.Pod)
		}
	}

	return r0
}

// GetPodsAndMirrorPods provides a mock function with given fields:
func (_m *MockManager) GetPodsAndMirrorPods() ([]*v1.Pod, []*v1.Pod) {
	ret := _m.Called()

	var r0 []*v1.Pod
	if rf, ok := ret.Get(0).(func() []*v1.Pod); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.Pod)
		}
	}

	var r1 []*v1.Pod
	if rf, ok := ret.Get(1).(func() []*v1.Pod); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*v1.Pod)
		}
	}

	return r0, r1
}

// GetUIDTranslations provides a mock function with given fields:
func (_m *MockManager) GetUIDTranslations() (map[kubelettypes.ResolvedPodUID]kubelettypes.MirrorPodUID, map[kubelettypes.MirrorPodUID]kubelettypes.ResolvedPodUID) {
	ret := _m.Called()

	var r0 map[kubelettypes.ResolvedPodUID]kubelettypes.MirrorPodUID
	if rf, ok := ret.Get(0).(func() map[kubelettypes.ResolvedPodUID]kubelettypes.MirrorPodUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[kubelettypes.ResolvedPodUID]kubelettypes.MirrorPodUID)
		}
	}

	var r1 map[kubelettypes.MirrorPodUID]kubelettypes.ResolvedPodUID
	if rf, ok := ret.Get(1).(func() map[kubelettypes.MirrorPodUID]kubelettypes.ResolvedPodUID); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[kubelettypes.MirrorPodUID]kubelettypes.ResolvedPodUID)
		}
	}

	return r0, r1
}

// IsMirrorPodOf provides a mock function with given fields: mirrorPod, _a1
func (_m *MockManager) IsMirrorPodOf(mirrorPod *v1.Pod, _a1 *v1.Pod) bool {
	ret := _m.Called(mirrorPod, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*v1.Pod, *v1.Pod) bool); ok {
		r0 = rf(mirrorPod, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetPods provides a mock function with given fields: pods
func (_m *MockManager) SetPods(pods []*v1.Pod) {
	_m.Called(pods)
}

// TranslatePodUID provides a mock function with given fields: uid
func (_m *MockManager) TranslatePodUID(uid types.UID) kubelettypes.ResolvedPodUID {
	ret := _m.Called(uid)

	var r0 kubelettypes.ResolvedPodUID
	if rf, ok := ret.Get(0).(func(types.UID) kubelettypes.ResolvedPodUID); ok {
		r0 = rf(uid)
	} else {
		r0 = ret.Get(0).(kubelettypes.ResolvedPodUID)
	}

	return r0
}

// UpdatePod provides a mock function with given fields: _a0
func (_m *MockManager) UpdatePod(_a0 *v1.Pod) {
	_m.Called(_a0)
}
