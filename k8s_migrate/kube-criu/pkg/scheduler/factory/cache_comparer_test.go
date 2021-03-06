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

package factory

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

func TestCompareNodes(t *testing.T) {
	tests := []struct {
		name      string
		actual    []string
		cached    []string
		missing   []string
		redundant []string
	}{
		{
			name:      "redundant cached value",
			actual:    []string{"foo", "bar"},
			cached:    []string{"bar", "foo", "foobar"},
			missing:   []string{},
			redundant: []string{"foobar"},
		},
		{
			name:      "missing cached value",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"bar", "foo"},
			missing:   []string{"foobar"},
			redundant: []string{},
		},
		{
			name:      "proper cache set",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"bar", "foobar", "foo"},
			missing:   []string{},
			redundant: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testCompareNodes(test.actual, test.cached, test.missing, test.redundant, t)
		})
	}
}

func testCompareNodes(actual, cached, missing, redundant []string, t *testing.T) {
	compare := compareStrategy{}
	nodes := []*v1.Node{}
	for _, nodeName := range actual {
		node := &v1.Node{}
		node.Name = nodeName
		nodes = append(nodes, node)
	}

	nodeInfo := make(map[string]*schedulercache.NodeInfo)
	for _, nodeName := range cached {
		nodeInfo[nodeName] = &schedulercache.NodeInfo{}
	}

	m, r := compare.CompareNodes(nodes, nodeInfo)

	if !reflect.DeepEqual(m, missing) {
		t.Errorf("missing expected to be %s; got %s", missing, m)
	}

	if !reflect.DeepEqual(r, redundant) {
		t.Errorf("redundant expected to be %s; got %s", redundant, r)
	}
}

func TestComparePods(t *testing.T) {
	tests := []struct {
		name      string
		actual    []string
		cached    []string
		queued    []string
		missing   []string
		redundant []string
	}{
		{
			name:      "redundant cached value",
			actual:    []string{"foo", "bar"},
			cached:    []string{"bar", "foo", "foobar"},
			queued:    []string{},
			missing:   []string{},
			redundant: []string{"foobar"},
		},
		{
			name:      "redundant and queued values",
			actual:    []string{"foo", "bar"},
			cached:    []string{"foo", "foobar"},
			queued:    []string{"bar"},
			missing:   []string{},
			redundant: []string{"foobar"},
		},
		{
			name:      "missing cached value",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"bar", "foo"},
			queued:    []string{},
			missing:   []string{"foobar"},
			redundant: []string{},
		},
		{
			name:      "missing and queued values",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"foo"},
			queued:    []string{"bar"},
			missing:   []string{"foobar"},
			redundant: []string{},
		},
		{
			name:      "correct cache set",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"bar", "foobar", "foo"},
			queued:    []string{},
			missing:   []string{},
			redundant: []string{},
		},
		{
			name:      "queued cache value",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"foobar", "foo"},
			queued:    []string{"bar"},
			missing:   []string{},
			redundant: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testComparePods(test.actual, test.cached, test.queued, test.missing, test.redundant, t)
		})
	}
}

func testComparePods(actual, cached, queued, missing, redundant []string, t *testing.T) {
	compare := compareStrategy{}
	pods := []*v1.Pod{}
	for _, uid := range actual {
		pod := &v1.Pod{}
		pod.UID = types.UID(uid)
		pods = append(pods, pod)
	}

	queuedPods := []*v1.Pod{}
	for _, uid := range queued {
		pod := &v1.Pod{}
		pod.UID = types.UID(uid)
		queuedPods = append(queuedPods, pod)
	}

	nodeInfo := make(map[string]*schedulercache.NodeInfo)
	for _, uid := range cached {
		pod := &v1.Pod{}
		pod.UID = types.UID(uid)
		pod.Namespace = "ns"
		pod.Name = uid

		nodeInfo[uid] = schedulercache.NewNodeInfo(pod)
	}

	m, r := compare.ComparePods(pods, queuedPods, nodeInfo)

	if !reflect.DeepEqual(m, missing) {
		t.Errorf("missing expected to be %s; got %s", missing, m)
	}

	if !reflect.DeepEqual(r, redundant) {
		t.Errorf("redundant expected to be %s; got %s", redundant, r)
	}
}

func TestComparePdbs(t *testing.T) {
	tests := []struct {
		name      string
		actual    []string
		cached    []string
		missing   []string
		redundant []string
	}{
		{
			name:      "redundant cache value",
			actual:    []string{"foo", "bar"},
			cached:    []string{"bar", "foo", "foobar"},
			missing:   []string{},
			redundant: []string{"foobar"},
		},
		{
			name:      "missing cache value",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"bar", "foo"},
			missing:   []string{"foobar"},
			redundant: []string{},
		},
		{
			name:      "correct cache",
			actual:    []string{"foo", "bar", "foobar"},
			cached:    []string{"bar", "foobar", "foo"},
			missing:   []string{},
			redundant: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testComparePdbs(test.actual, test.cached, test.missing, test.redundant, t)
		})
	}
}

func testComparePdbs(actual, cached, missing, redundant []string, t *testing.T) {
	compare := compareStrategy{}
	pdbs := []*policy.PodDisruptionBudget{}
	for _, uid := range actual {
		pdb := &policy.PodDisruptionBudget{}
		pdb.UID = types.UID(uid)
		pdbs = append(pdbs, pdb)
	}

	cache := make(map[string]*policy.PodDisruptionBudget)
	for _, uid := range cached {
		pdb := &policy.PodDisruptionBudget{}
		pdb.UID = types.UID(uid)
		cache[uid] = pdb
	}

	m, r := compare.ComparePdbs(pdbs, cache)

	if !reflect.DeepEqual(m, missing) {
		t.Errorf("missing expected to be %s; got %s", missing, m)
	}

	if !reflect.DeepEqual(r, redundant) {
		t.Errorf("redundant expected to be %s; got %s", redundant, r)
	}
}
