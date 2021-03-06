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

package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/api/testapi"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/validation"
	"k8s.io/kubernetes/pkg/securitycontext"
)

func noDefault(*core.Pod) error { return nil }

func TestDecodeSinglePod(t *testing.T) {
	grace := int64(30)
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			UID:       "12345",
			Namespace: "mynamespace",
		},
		Spec: v1.PodSpec{
			RestartPolicy:                 v1.RestartPolicyAlways,
			DNSPolicy:                     v1.DNSClusterFirst,
			TerminationGracePeriodSeconds: &grace,
			Containers: []v1.Container{{
				Name:                     "image",
				Image:                    "test/image",
				ImagePullPolicy:          "IfNotPresent",
				TerminationMessagePath:   "/dev/termination-log",
				TerminationMessagePolicy: v1.TerminationMessageReadFile,
				SecurityContext:          securitycontext.ValidSecurityContextWithContainerDefaults(),
			}},
			SecurityContext: &v1.PodSecurityContext{},
			SchedulerName:   core.DefaultSchedulerName,
		},
	}
	json, err := runtime.Encode(testapi.Default.Codec(), pod)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	parsed, podOut, err := tryDecodeSinglePod(json, noDefault)
	if !parsed {
		t.Errorf("expected to have parsed file: (%s)", string(json))
	}
	if err != nil {
		t.Errorf("unexpected error: %v (%s)", err, string(json))
	}
	if !reflect.DeepEqual(pod, podOut) {
		t.Errorf("expected:\n%#v\ngot:\n%#v\n%s", pod, podOut, string(json))
	}

	for _, gv := range legacyscheme.Scheme.PrioritizedVersionsForGroup(v1.GroupName) {
		info, _ := runtime.SerializerInfoForMediaType(legacyscheme.Codecs.SupportedMediaTypes(), "application/yaml")
		encoder := legacyscheme.Codecs.EncoderForVersion(info.Serializer, gv)
		yaml, err := runtime.Encode(encoder, pod)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		parsed, podOut, err = tryDecodeSinglePod(yaml, noDefault)
		if !parsed {
			t.Errorf("expected to have parsed file: (%s)", string(yaml))
		}
		if err != nil {
			t.Errorf("unexpected error: %v (%s)", err, string(yaml))
		}
		if !reflect.DeepEqual(pod, podOut) {
			t.Errorf("expected:\n%#v\ngot:\n%#v\n%s", pod, podOut, string(yaml))
		}
	}
}

func TestDecodePodList(t *testing.T) {
	grace := int64(30)
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			UID:       "12345",
			Namespace: "mynamespace",
		},
		Spec: v1.PodSpec{
			RestartPolicy:                 v1.RestartPolicyAlways,
			DNSPolicy:                     v1.DNSClusterFirst,
			TerminationGracePeriodSeconds: &grace,
			Containers: []v1.Container{{
				Name:                     "image",
				Image:                    "test/image",
				ImagePullPolicy:          "IfNotPresent",
				TerminationMessagePath:   "/dev/termination-log",
				TerminationMessagePolicy: v1.TerminationMessageReadFile,

				SecurityContext: securitycontext.ValidSecurityContextWithContainerDefaults(),
			}},
			SecurityContext: &v1.PodSecurityContext{},
			SchedulerName:   core.DefaultSchedulerName,
		},
	}
	podList := &v1.PodList{
		Items: []v1.Pod{*pod},
	}
	json, err := runtime.Encode(testapi.Default.Codec(), podList)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	parsed, podListOut, err := tryDecodePodList(json, noDefault)
	if !parsed {
		t.Errorf("expected to have parsed file: (%s)", string(json))
	}
	if err != nil {
		t.Errorf("unexpected error: %v (%s)", err, string(json))
	}
	if !reflect.DeepEqual(podList, &podListOut) {
		t.Errorf("expected:\n%#v\ngot:\n%#v\n%s", podList, &podListOut, string(json))
	}

	for _, gv := range legacyscheme.Scheme.PrioritizedVersionsForGroup(v1.GroupName) {
		info, _ := runtime.SerializerInfoForMediaType(legacyscheme.Codecs.SupportedMediaTypes(), "application/yaml")
		encoder := legacyscheme.Codecs.EncoderForVersion(info.Serializer, gv)
		yaml, err := runtime.Encode(encoder, podList)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		parsed, podListOut, err = tryDecodePodList(yaml, noDefault)
		if !parsed {
			t.Errorf("expected to have parsed file: (%s): %v", string(yaml), err)
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %v (%s)", err, string(yaml))
			continue
		}
		if !reflect.DeepEqual(podList, &podListOut) {
			t.Errorf("expected:\n%#v\ngot:\n%#v\n%s", pod, &podListOut, string(yaml))
		}
	}
}

func TestGetSelfLink(t *testing.T) {
	var testCases = []struct {
		desc             string
		name             string
		namespace        string
		expectedSelfLink string
	}{
		{
			desc:             "No namespace specified",
			name:             "foo",
			namespace:        "",
			expectedSelfLink: "/api/v1/namespaces/default/pods/foo",
		},
		{
			desc:             "Namespace specified",
			name:             "foo",
			namespace:        "bar",
			expectedSelfLink: "/api/v1/namespaces/bar/pods/foo",
		},
	}
	for _, testCase := range testCases {
		selfLink := getSelfLink(testCase.name, testCase.namespace)
		if testCase.expectedSelfLink != selfLink {
			t.Errorf("%s: getSelfLink error, expected: %s, got: %s", testCase.desc, testCase.expectedSelfLink, selfLink)
		}
	}
}

func TestStaticPodNameGenerate(t *testing.T) {
	testCases := []struct {
		nodeName  types.NodeName
		podName   string
		expected  string
		overwrite string
		shouldErr bool
	}{
		{
			"node1",
			"static-pod1",
			"static-pod1-node1",
			"",
			false,
		},
		{
			"Node1",
			"static-pod1",
			"static-pod1-node1",
			"",
			false,
		},
		{
			"NODE1",
			"static-pod1",
			"static-pod1-node1",
			"static-pod1-NODE1",
			true,
		},
	}
	for _, c := range testCases {
		assert.Equal(t, c.expected, generatePodName(c.podName, c.nodeName), "wrong pod name generated")
		pod := &core.Pod{}
		pod.Name = c.podName
		if c.overwrite != "" {
			pod.Name = c.overwrite
		}
		errs := validation.ValidatePod(pod)
		if c.shouldErr {
			specNameErrored := false
			for _, err := range errs {
				if err.Field == "metadata.name" {
					specNameErrored = true
				}
			}
			assert.NotEmpty(t, specNameErrored, "expecting error")
		} else {
			for _, err := range errs {
				if err.Field == "metadata.name" {
					t.Fail()
				}
			}
		}
	}
}
