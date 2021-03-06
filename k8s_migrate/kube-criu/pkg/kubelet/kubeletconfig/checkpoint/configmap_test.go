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

package checkpoint

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	apiv1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utiltest "k8s.io/kubernetes/pkg/kubelet/kubeletconfig/util/test"
)

func TestNewConfigMapPayload(t *testing.T) {
	cases := []struct {
		desc string
		cm   *apiv1.ConfigMap
		err  string
	}{
		{
			desc: "nil",
			cm:   nil,
			err:  "ConfigMap must be non-nil",
		},
		{
			desc: "missing uid",
			cm: &apiv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "name",
					ResourceVersion: "rv",
				},
			},
			err: "ConfigMap must have a non-empty UID",
		},
		{
			desc: "missing resourceVersion",
			cm: &apiv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "name",
					UID:  "uid",
				},
			},
			err: "ConfigMap must have a non-empty ResourceVersion",
		},
		{
			desc: "populated v1/ConfigMap",
			cm: &apiv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "name",
					UID:             "uid",
					ResourceVersion: "rv",
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			err: "",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			payload, err := NewConfigMapPayload(c.cm)
			utiltest.ExpectError(t, err, c.err)
			if err != nil {
				return
			}
			// underlying object should match the object passed in
			if !apiequality.Semantic.DeepEqual(c.cm, payload.object()) {
				t.Errorf("expect %s but got %s", spew.Sdump(c.cm), spew.Sdump(payload))
			}
		})
	}
}

func TestConfigMapPayloadUID(t *testing.T) {
	const expect = "uid"
	payload, err := NewConfigMapPayload(&apiv1.ConfigMap{ObjectMeta: metav1.ObjectMeta{UID: expect, ResourceVersion: "rv"}})
	if err != nil {
		t.Fatalf("error constructing payload: %v", err)
	}
	uid := payload.UID()
	if expect != uid {
		t.Errorf("expect %q, but got %q", expect, uid)
	}
}

func TestConfigMapPayloadResourceVersion(t *testing.T) {
	const expect = "rv"
	payload, err := NewConfigMapPayload(&apiv1.ConfigMap{ObjectMeta: metav1.ObjectMeta{UID: "uid", ResourceVersion: expect}})
	if err != nil {
		t.Fatalf("error constructing payload: %v", err)
	}
	resourceVersion := payload.ResourceVersion()
	if expect != resourceVersion {
		t.Errorf("expect %q, but got %q", expect, resourceVersion)
	}
}

func TestConfigMapPayloadFiles(t *testing.T) {
	cases := []struct {
		desc   string
		data   map[string]string
		expect map[string]string
	}{
		{"nil", nil, nil},
		{"empty", map[string]string{}, map[string]string{}},
		{"populated",
			map[string]string{
				"foo": "1",
				"bar": "2",
			},
			map[string]string{
				"foo": "1",
				"bar": "2",
			}},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			payload, err := NewConfigMapPayload(&apiv1.ConfigMap{ObjectMeta: metav1.ObjectMeta{UID: "uid", ResourceVersion: "rv"}, Data: c.data})
			if err != nil {
				t.Fatalf("error constructing payload: %v", err)
			}
			files := payload.Files()
			if !reflect.DeepEqual(c.expect, files) {
				t.Errorf("expected %v, but got %v", c.expect, files)
			}
		})
	}
}
