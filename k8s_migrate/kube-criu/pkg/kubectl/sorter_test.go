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
Copyright 2014 The Kubernetes Authors.

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

package kubectl

import (
	"reflect"
	"strings"
	"testing"

	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

func encodeOrDie(obj runtime.Object) []byte {
	data, err := runtime.Encode(legacyscheme.Codecs.LegacyCodec(api.SchemeGroupVersion), obj)
	if err != nil {
		panic(err.Error())
	}
	return data
}

func TestSortingPrinter(t *testing.T) {
	intPtr := func(val int32) *int32 { return &val }

	a := &api.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "a",
		},
	}

	b := &api.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "b",
		},
	}

	c := &api.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "c",
		},
	}

	tests := []struct {
		obj         runtime.Object
		sort        runtime.Object
		field       string
		name        string
		expectedErr string
	}{
		{
			name: "in-order-already",
			obj: &api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "a",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "b",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "c",
						},
					},
				},
			},
			sort: &api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "a",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "b",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "c",
						},
					},
				},
			},
			field: "{.metadata.name}",
		},
		{
			name: "reverse-order",
			obj: &api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "b",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "c",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "a",
						},
					},
				},
			},
			sort: &api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "a",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "b",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "c",
						},
					},
				},
			},
			field: "{.metadata.name}",
		},
		{
			name: "random-order-timestamp",
			obj: &api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							CreationTimestamp: metav1.Unix(300, 0),
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							CreationTimestamp: metav1.Unix(100, 0),
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							CreationTimestamp: metav1.Unix(200, 0),
						},
					},
				},
			},
			sort: &api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							CreationTimestamp: metav1.Unix(100, 0),
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							CreationTimestamp: metav1.Unix(200, 0),
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							CreationTimestamp: metav1.Unix(300, 0),
						},
					},
				},
			},
			field: "{.metadata.creationTimestamp}",
		},
		{
			name: "random-order-numbers",
			obj: &api.ReplicationControllerList{
				Items: []api.ReplicationController{
					{
						Spec: api.ReplicationControllerSpec{
							Replicas: intPtr(5),
						},
					},
					{
						Spec: api.ReplicationControllerSpec{
							Replicas: intPtr(1),
						},
					},
					{
						Spec: api.ReplicationControllerSpec{
							Replicas: intPtr(9),
						},
					},
				},
			},
			sort: &api.ReplicationControllerList{
				Items: []api.ReplicationController{
					{
						Spec: api.ReplicationControllerSpec{
							Replicas: intPtr(1),
						},
					},
					{
						Spec: api.ReplicationControllerSpec{
							Replicas: intPtr(5),
						},
					},
					{
						Spec: api.ReplicationControllerSpec{
							Replicas: intPtr(9),
						},
					},
				},
			},
			field: "{.spec.replicas}",
		},
		{
			name: "v1.List in order",
			obj: &api.List{
				Items: []runtime.RawExtension{
					{Raw: encodeOrDie(a)},
					{Raw: encodeOrDie(b)},
					{Raw: encodeOrDie(c)},
				},
			},
			sort: &api.List{
				Items: []runtime.RawExtension{
					{Raw: encodeOrDie(a)},
					{Raw: encodeOrDie(b)},
					{Raw: encodeOrDie(c)},
				},
			},
			field: "{.metadata.name}",
		},
		{
			name: "v1.List in reverse",
			obj: &api.List{
				Items: []runtime.RawExtension{
					{Raw: encodeOrDie(c)},
					{Raw: encodeOrDie(b)},
					{Raw: encodeOrDie(a)},
				},
			},
			sort: &api.List{
				Items: []runtime.RawExtension{
					{Raw: encodeOrDie(a)},
					{Raw: encodeOrDie(b)},
					{Raw: encodeOrDie(c)},
				},
			},
			field: "{.metadata.name}",
		},
		{
			name: "some-missing-fields",
			obj: &unstructured.UnstructuredList{
				Object: map[string]interface{}{
					"kind":       "List",
					"apiVersion": "v1",
				},
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status": map[string]interface{}{
								"availableReplicas": 2,
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status":     map[string]interface{}{},
						},
					},
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status": map[string]interface{}{
								"availableReplicas": 1,
							},
						},
					},
				},
			},
			sort: &unstructured.UnstructuredList{
				Object: map[string]interface{}{
					"kind":       "List",
					"apiVersion": "v1",
				},
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status":     map[string]interface{}{},
						},
					},
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status": map[string]interface{}{
								"availableReplicas": 1,
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status": map[string]interface{}{
								"availableReplicas": 2,
							},
						},
					},
				},
			},
			field: "{.status.availableReplicas}",
		},
		{
			name: "all-missing-fields",
			obj: &unstructured.UnstructuredList{
				Object: map[string]interface{}{
					"kind":       "List",
					"apiVersion": "v1",
				},
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status": map[string]interface{}{
								"replicas": 0,
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind":       "ReplicationController",
							"apiVersion": "v1",
							"status": map[string]interface{}{
								"replicas": 0,
							},
						},
					},
				},
			},
			field:       "{.status.availableReplicas}",
			expectedErr: "couldn't find any field with path \"{.status.availableReplicas}\" in the list of objects",
		},
		{
			name: "model-invalid-fields",
			obj: &api.ReplicationControllerList{
				Items: []api.ReplicationController{
					{
						Status: api.ReplicationControllerStatus{},
					},
					{
						Status: api.ReplicationControllerStatus{},
					},
					{
						Status: api.ReplicationControllerStatus{},
					},
				},
			},
			field:       "{.invalid}",
			expectedErr: "couldn't find any field with path \"{.invalid}\" in the list of objects",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort := &SortingPrinter{SortField: tt.field, Decoder: legacyscheme.Codecs.UniversalDecoder()}
			err := sort.sortObj(tt.obj)
			if err != nil {
				if len(tt.expectedErr) > 0 {
					if strings.Contains(err.Error(), tt.expectedErr) {
						return
					}
					t.Fatalf("%s: expected error containing: %q, got: \"%v\"", tt.name, tt.expectedErr, err)
				}
				t.Fatalf("%s: unexpected error: %v", tt.name, err)
			}
			if len(tt.expectedErr) > 0 {
				t.Fatalf("%s: expected error containing: %q, got none", tt.name, tt.expectedErr)
			}
			if !reflect.DeepEqual(tt.obj, tt.sort) {
				t.Errorf("[%s]\nexpected:\n%v\nsaw:\n%v", tt.name, tt.sort, tt.obj)
			}
		})
	}
}

func TestRuntimeSortLess(t *testing.T) {
	var testobj runtime.Object

	testobj = &api.PodList{
		Items: []api.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "b",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "a",
				},
			},
		},
	}

	testobjs, err := meta.ExtractList(testobj)
	if err != nil {
		t.Fatalf("ExtractList testobj got unexpected error: %v", err)
	}

	testfield := "{.metadata.name}"
	testruntimeSort := NewRuntimeSort(testfield, testobjs)
	tests := []struct {
		name         string
		runtimeSort  *RuntimeSort
		i            int
		j            int
		expectResult bool
		expectErr    bool
	}{
		{
			name:         "test less true",
			runtimeSort:  testruntimeSort,
			i:            0,
			j:            1,
			expectResult: true,
		},
		{
			name:         "test less false",
			runtimeSort:  testruntimeSort,
			i:            1,
			j:            2,
			expectResult: false,
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.runtimeSort.Less(test.i, test.j)
			if result != test.expectResult {
				t.Errorf("case[%d]:%s Expected result: %v, Got result: %v", i, test.name, test.expectResult, result)
			}
		})
	}
}
