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

package validation

import (
	"strconv"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/apis/apps"
	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestValidateStatefulSet(t *testing.T) {
	validLabels := map[string]string{"a": "b"}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
			},
		},
	}

	invalidLabels := map[string]string{"NoUppercaseOrSpecialCharsLike=Equals": "b"}
	invalidPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
			},
			ObjectMeta: metav1.ObjectMeta{
				Labels: invalidLabels,
			},
		},
	}

	invalidTime := int64(60)
	invalidPodTemplate2 := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"foo": "bar"},
			},
			Spec: api.PodSpec{
				RestartPolicy:         api.RestartPolicyOnFailure,
				DNSPolicy:             api.DNSClusterFirst,
				ActiveDeadlineSeconds: &invalidTime,
			},
		},
	}

	successCases := []apps.StatefulSet{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.ParallelPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.OnDeleteStatefulSetStrategyType},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				Replicas:            3,
				UpdateStrategy: apps.StatefulSetUpdateStrategy{
					Type: apps.RollingUpdateStatefulSetStrategyType,
					RollingUpdate: func() *apps.RollingUpdateStatefulSetStrategy {
						return &apps.RollingUpdateStatefulSetStrategy{Partition: 2}
					}()},
			},
		},
	}

	for i, successCase := range successCases {
		t.Run("success case "+strconv.Itoa(i), func(t *testing.T) {
			if errs := ValidateStatefulSet(&successCase); len(errs) != 0 {
				t.Errorf("expected success: %v", errs)
			}
		})
	}

	errorCases := map[string]apps.StatefulSet{
		"zero-length ID": {
			ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"missing-namespace": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123"},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"empty selector": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"selector_doesnt_match": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid manifest": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"negative_replicas": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Replicas:            -1,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid_label": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
				Labels: map[string]string{
					"NoUppercaseOrSpecialCharsLike=Equals": "bar",
				},
			},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid_label 2": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
				Labels: map[string]string{
					"NoUppercaseOrSpecialCharsLike=Equals": "bar",
				},
			},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Template:            invalidPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid_annotation": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					"NoUppercaseOrSpecialCharsLike=Equals": "bar",
				},
			},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid restart policy 1": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template: api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyOnFailure,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent"}},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: validLabels,
					},
				},
				UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid restart policy 2": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template: api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyNever,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent"}},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: validLabels,
					},
				},
				UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid update strategy": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				Replicas:            3,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: "foo"},
			},
		},
		"empty update strategy": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				Replicas:            3,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: ""},
			},
		},
		"invalid rolling update": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				Replicas:            3,
				UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.OnDeleteStatefulSetStrategyType,
					RollingUpdate: func() *apps.RollingUpdateStatefulSetStrategy {
						return &apps.RollingUpdateStatefulSetStrategy{Partition: 1}
					}()},
			},
		},
		"negative parition": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: apps.OrderedReadyPodManagement,
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				Replicas:            3,
				UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType,
					RollingUpdate: func() *apps.RollingUpdateStatefulSetStrategy {
						return &apps.RollingUpdateStatefulSetStrategy{Partition: -1}
					}()},
			},
		},
		"empty pod management policy": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: "",
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				Replicas:            3,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"invalid pod management policy": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: "foo",
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            validPodTemplate.Template,
				Replicas:            3,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
		"set active deadline seconds": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: apps.StatefulSetSpec{
				PodManagementPolicy: "foo",
				Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
				Template:            invalidPodTemplate2.Template,
				Replicas:            3,
				UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
			},
		},
	}

	for k, v := range errorCases {
		t.Run(k, func(t *testing.T) {
			errs := ValidateStatefulSet(&v)
			if len(errs) == 0 {
				t.Errorf("expected failure for %s", k)
			}

			for i := range errs {
				field := errs[i].Field
				if !strings.HasPrefix(field, "spec.template.") &&
					field != "metadata.name" &&
					field != "metadata.namespace" &&
					field != "spec.selector" &&
					field != "spec.template" &&
					field != "GCEPersistentDisk.ReadOnly" &&
					field != "spec.replicas" &&
					field != "spec.template.labels" &&
					field != "metadata.annotations" &&
					field != "metadata.labels" &&
					field != "status.replicas" &&
					field != "spec.updateStrategy" &&
					field != "spec.updateStrategy.rollingUpdate" &&
					field != "spec.updateStrategy.rollingUpdate.partition" &&
					field != "spec.podManagementPolicy" &&
					field != "spec.template.spec.activeDeadlineSeconds" {
					t.Errorf("%s: missing prefix for: %v", k, errs[i])
				}
			}
		})
	}
}

func TestValidateStatefulSetStatus(t *testing.T) {
	observedGenerationMinusOne := int64(-1)
	collisionCountMinusOne := int32(-1)
	tests := []struct {
		name               string
		replicas           int32
		readyReplicas      int32
		currentReplicas    int32
		updatedReplicas    int32
		observedGeneration *int64
		collisionCount     *int32
		expectedErr        bool
	}{
		{
			name:            "valid status",
			replicas:        3,
			readyReplicas:   3,
			currentReplicas: 2,
			updatedReplicas: 1,
			expectedErr:     false,
		},
		{
			name:            "invalid replicas",
			replicas:        -1,
			readyReplicas:   3,
			currentReplicas: 2,
			updatedReplicas: 1,
			expectedErr:     true,
		},
		{
			name:            "invalid readyReplicas",
			replicas:        3,
			readyReplicas:   -1,
			currentReplicas: 2,
			updatedReplicas: 1,
			expectedErr:     true,
		},
		{
			name:            "invalid currentReplicas",
			replicas:        3,
			readyReplicas:   3,
			currentReplicas: -1,
			updatedReplicas: 1,
			expectedErr:     true,
		},
		{
			name:            "invalid updatedReplicas",
			replicas:        3,
			readyReplicas:   3,
			currentReplicas: 2,
			updatedReplicas: -1,
			expectedErr:     true,
		},
		{
			name:               "invalid observedGeneration",
			replicas:           3,
			readyReplicas:      3,
			currentReplicas:    2,
			updatedReplicas:    1,
			observedGeneration: &observedGenerationMinusOne,
			expectedErr:        true,
		},
		{
			name:            "invalid collisionCount",
			replicas:        3,
			readyReplicas:   3,
			currentReplicas: 2,
			updatedReplicas: 1,
			collisionCount:  &collisionCountMinusOne,
			expectedErr:     true,
		},
		{
			name:            "readyReplicas greater than replicas",
			replicas:        3,
			readyReplicas:   4,
			currentReplicas: 2,
			updatedReplicas: 1,
			expectedErr:     true,
		},
		{
			name:            "currentReplicas greater than replicas",
			replicas:        3,
			readyReplicas:   3,
			currentReplicas: 4,
			updatedReplicas: 1,
			expectedErr:     true,
		},
		{
			name:            "updatedReplicas greater than replicas",
			replicas:        3,
			readyReplicas:   3,
			currentReplicas: 2,
			updatedReplicas: 4,
			expectedErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := apps.StatefulSetStatus{
				Replicas:           test.replicas,
				ReadyReplicas:      test.readyReplicas,
				CurrentReplicas:    test.currentReplicas,
				UpdatedReplicas:    test.updatedReplicas,
				ObservedGeneration: test.observedGeneration,
				CollisionCount:     test.collisionCount,
			}

			errs := ValidateStatefulSetStatus(&status, field.NewPath("status"))
			if hasErr := len(errs) > 0; hasErr != test.expectedErr {
				t.Errorf("%s: expected error: %t, got error: %t\nerrors: %s", test.name, test.expectedErr, hasErr, errs.ToAggregate().Error())
			}
		})
	}
}

func TestValidateStatefulSetUpdate(t *testing.T) {
	validLabels := map[string]string{"a": "b"}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
			},
		},
	}

	addContainersValidTemplate := validPodTemplate.DeepCopy()
	addContainersValidTemplate.Template.Spec.Containers = append(addContainersValidTemplate.Template.Spec.Containers,
		api.Container{Name: "def", Image: "image2", ImagePullPolicy: "IfNotPresent"})
	if len(addContainersValidTemplate.Template.Spec.Containers) != len(validPodTemplate.Template.Spec.Containers)+1 {
		t.Errorf("failure during test setup: template %v should have more containers than template %v", addContainersValidTemplate, validPodTemplate)
	}

	readWriteVolumePodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
				Volumes:       []api.Volume{{Name: "gcepd", VolumeSource: api.VolumeSource{GCEPersistentDisk: &api.GCEPersistentDiskVolumeSource{PDName: "my-PD", FSType: "ext4", Partition: 1, ReadOnly: false}}}},
			},
		},
	}
	invalidLabels := map[string]string{"NoUppercaseOrSpecialCharsLike=Equals": "b"}
	invalidPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
			},
			ObjectMeta: metav1.ObjectMeta{
				Labels: invalidLabels,
			},
		},
	}

	type psUpdateTest struct {
		old    apps.StatefulSet
		update apps.StatefulSet
	}
	successCases := []psUpdateTest{
		{
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Replicas:            3,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
		{
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            addContainersValidTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
		{
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            addContainersValidTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
	}

	for i, successCase := range successCases {
		t.Run("success case "+strconv.Itoa(i), func(t *testing.T) {
			successCase.old.ObjectMeta.ResourceVersion = "1"
			successCase.update.ObjectMeta.ResourceVersion = "1"
			if errs := ValidateStatefulSetUpdate(&successCase.update, &successCase.old); len(errs) != 0 {
				t.Errorf("expected success: %v", errs)
			}
		})
	}

	errorCases := map[string]psUpdateTest{
		"more than one read/write": {
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					Selector:       &metav1.LabelSelector{MatchLabels: validLabels},
					Template:       validPodTemplate.Template,
					UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Replicas:            2,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            readWriteVolumePodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
		"empty pod creation policy": {
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					Selector:       &metav1.LabelSelector{MatchLabels: validLabels},
					Template:       validPodTemplate.Template,
					UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					Replicas:       3,
					Selector:       &metav1.LabelSelector{MatchLabels: validLabels},
					Template:       validPodTemplate.Template,
					UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
		"invalid pod creation policy": {
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					Selector:       &metav1.LabelSelector{MatchLabels: validLabels},
					Template:       validPodTemplate.Template,
					UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.PodManagementPolicyType("Other"),
					Replicas:            3,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
		"invalid selector": {
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					Selector:       &metav1.LabelSelector{MatchLabels: validLabels},
					Template:       validPodTemplate.Template,
					UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Replicas:            2,
					Selector:            &metav1.LabelSelector{MatchLabels: invalidLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
		"invalid pod": {
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					Replicas:       2,
					Selector:       &metav1.LabelSelector{MatchLabels: validLabels},
					Template:       invalidPodTemplate.Template,
					UpdateStrategy: apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
		"negative replicas": {
			old: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
			update: apps.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: apps.StatefulSetSpec{
					PodManagementPolicy: apps.OrderedReadyPodManagement,
					Replicas:            -1,
					Selector:            &metav1.LabelSelector{MatchLabels: validLabels},
					Template:            validPodTemplate.Template,
					UpdateStrategy:      apps.StatefulSetUpdateStrategy{Type: apps.RollingUpdateStatefulSetStrategyType},
				},
			},
		},
	}

	for testName, errorCase := range errorCases {
		t.Run(testName, func(t *testing.T) {
			if errs := ValidateStatefulSetUpdate(&errorCase.update, &errorCase.old); len(errs) == 0 {
				t.Errorf("expected failure: %s", testName)
			}
		})
	}
}

func TestValidateControllerRevision(t *testing.T) {
	newControllerRevision := func(name, namespace string, data runtime.Object, revision int64) apps.ControllerRevision {
		return apps.ControllerRevision{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data:     data,
			Revision: revision,
		}
	}

	ss := apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
		Spec: apps.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
			Template: api.PodTemplateSpec{
				Spec: api.PodSpec{
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSClusterFirst,
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"foo": "bar"},
				},
			},
		},
	}

	var (
		valid       = newControllerRevision("validname", "validns", &ss, 0)
		badRevision = newControllerRevision("validname", "validns", &ss, -1)
		emptyName   = newControllerRevision("", "validns", &ss, 0)
		invalidName = newControllerRevision("NoUppercaseOrSpecialCharsLike=Equals", "validns", &ss, 0)
		emptyNs     = newControllerRevision("validname", "", &ss, 100)
		invalidNs   = newControllerRevision("validname", "NoUppercaseOrSpecialCharsLike=Equals", &ss, 100)
		nilData     = newControllerRevision("validname", "NoUppercaseOrSpecialCharsLike=Equals", nil, 100)
	)

	tests := map[string]struct {
		history apps.ControllerRevision
		isValid bool
	}{
		"valid":             {valid, true},
		"negative revision": {badRevision, false},
		"empty name":        {emptyName, false},
		"invalid name":      {invalidName, false},
		"empty namespace":   {emptyNs, false},
		"invalid namespace": {invalidNs, false},
		"nil data":          {nilData, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			errs := ValidateControllerRevision(&tc.history)
			if tc.isValid && len(errs) > 0 {
				t.Errorf("%v: unexpected error: %v", name, errs)
			}
			if !tc.isValid && len(errs) == 0 {
				t.Errorf("%v: unexpected non-error", name)
			}
		})
	}
}

func TestValidateControllerRevisionUpdate(t *testing.T) {
	newControllerRevision := func(version, name, namespace string, data runtime.Object, revision int64) apps.ControllerRevision {
		return apps.ControllerRevision{
			ObjectMeta: metav1.ObjectMeta{
				Name:            name,
				Namespace:       namespace,
				ResourceVersion: version,
			},
			Data:     data,
			Revision: revision,
		}
	}

	ss := apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
		Spec: apps.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
			Template: api.PodTemplateSpec{
				Spec: api.PodSpec{
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSClusterFirst,
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"foo": "bar"},
				},
			},
		},
	}
	modifiedss := apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "cdf", Namespace: metav1.NamespaceDefault},
		Spec: apps.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
			Template: api.PodTemplateSpec{
				Spec: api.PodSpec{
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSClusterFirst,
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"foo": "bar"},
				},
			},
		},
	}

	var (
		valid           = newControllerRevision("1", "validname", "validns", &ss, 0)
		noVersion       = newControllerRevision("", "validname", "validns", &ss, 0)
		changedData     = newControllerRevision("1", "validname", "validns", &modifiedss, 0)
		changedRevision = newControllerRevision("1", "validname", "validns", &ss, 1)
	)

	cases := []struct {
		name       string
		newHistory apps.ControllerRevision
		oldHistory apps.ControllerRevision
		isValid    bool
	}{
		{
			name:       "valid",
			newHistory: valid,
			oldHistory: valid,
			isValid:    true,
		},
		{
			name:       "invalid",
			newHistory: noVersion,
			oldHistory: valid,
			isValid:    false,
		},
		{
			name:       "changed data",
			newHistory: changedData,
			oldHistory: valid,
			isValid:    false,
		},
		{
			name:       "changed revision",
			newHistory: changedRevision,
			oldHistory: valid,
			isValid:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateControllerRevisionUpdate(&tc.newHistory, &tc.oldHistory)
			if tc.isValid && len(errs) > 0 {
				t.Errorf("%v: unexpected error: %v", tc.name, errs)
			}
			if !tc.isValid && len(errs) == 0 {
				t.Errorf("%v: unexpected non-error", tc.name)
			}
		})
	}
}
