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

package validation

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/validation/field"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func TestValidateDaemonSetStatusUpdate(t *testing.T) {
	type dsUpdateTest struct {
		old    extensions.DaemonSet
		update extensions.DaemonSet
	}

	successCases := []dsUpdateTest{
		{
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
		},
	}

	for _, successCase := range successCases {
		successCase.old.ObjectMeta.ResourceVersion = "1"
		successCase.update.ObjectMeta.ResourceVersion = "1"
		if errs := ValidateDaemonSetStatusUpdate(&successCase.update, &successCase.old); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}
	errorCases := map[string]dsUpdateTest{
		"negative values": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: -1,
					NumberMisscheduled:     -1,
					DesiredNumberScheduled: -3,
					NumberReady:            -1,
					ObservedGeneration:     -3,
					UpdatedNumberScheduled: -1,
					NumberAvailable:        -1,
					NumberUnavailable:      -2,
				},
			},
		},
		"negative CurrentNumberScheduled": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: -1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
		},
		"negative NumberMisscheduled": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     -1,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
		},
		"negative DesiredNumberScheduled": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: -3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
		},
		"negative NumberReady": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
					NumberReady:            -1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
		},
		"negative ObservedGeneration": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     -3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
		},
		"negative UpdatedNumberScheduled": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: -1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
		},
		"negative NumberAvailable": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        -1,
					NumberUnavailable:      2,
				},
			},
		},
		"negative NumberUnavailable": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      2,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "abc",
					Namespace:       metav1.NamespaceDefault,
					ResourceVersion: "10",
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					ObservedGeneration:     3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        1,
					NumberUnavailable:      -2,
				},
			},
		},
	}

	for testName, errorCase := range errorCases {
		if errs := ValidateDaemonSetStatusUpdate(&errorCase.update, &errorCase.old); len(errs) == 0 {
			t.Errorf("expected failure: %s", testName)
		}
	}
}

func TestValidateDaemonSetUpdate(t *testing.T) {
	validSelector := map[string]string{"a": "b"}
	validSelector2 := map[string]string{"c": "d"}
	invalidSelector := map[string]string{"NoUppercaseOrSpecialCharsLike=Equals": "b"}

	validPodSpecAbc := api.PodSpec{
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
	}
	validPodSpecDef := api.PodSpec{
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "def", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
	}
	validPodSpecNodeSelector := api.PodSpec{
		NodeSelector:  validSelector,
		NodeName:      "xyz",
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
	}
	validPodSpecVolume := api.PodSpec{
		Volumes:       []api.Volume{{Name: "gcepd", VolumeSource: api.VolumeSource{GCEPersistentDisk: &api.GCEPersistentDiskVolumeSource{PDName: "my-PD", FSType: "ext4", Partition: 1, ReadOnly: false}}}},
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
	}

	validPodTemplateAbc := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector,
			},
			Spec: validPodSpecAbc,
		},
	}
	validPodTemplateAbcSemanticallyEqual := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector,
			},
			Spec: validPodSpecAbc,
		},
	}
	validPodTemplateAbcSemanticallyEqual.Template.Spec.ImagePullSecrets = []api.LocalObjectReference{}
	validPodTemplateNodeSelector := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector,
			},
			Spec: validPodSpecNodeSelector,
		},
	}
	validPodTemplateAbc2 := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector2,
			},
			Spec: validPodSpecAbc,
		},
	}
	validPodTemplateDef := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector2,
			},
			Spec: validPodSpecDef,
		},
	}
	invalidPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			Spec: api.PodSpec{
				// no containers specified
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
			},
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector,
			},
		},
	}
	readWriteVolumePodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector,
			},
			Spec: validPodSpecVolume,
		},
	}

	type dsUpdateTest struct {
		old            extensions.DaemonSet
		update         extensions.DaemonSet
		expectedErrNum int
	}
	successCases := map[string]dsUpdateTest{
		"no change": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
		},
		"change template and selector": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 2,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector2},
					TemplateGeneration: 3,
					Template:           validPodTemplateAbc2.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
		},
		"change template": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 3,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 4,
					Template:           validPodTemplateNodeSelector.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
		},
		"change container image name": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector2},
					TemplateGeneration: 2,
					Template:           validPodTemplateDef.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
		},
		"change update strategy": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 4,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 4,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.RollingUpdateDaemonSetStrategyType,
						RollingUpdate: &extensions.RollingUpdateDaemonSet{
							MaxUnavailable: intstr.FromInt(1),
						},
					},
				},
			},
		},
		"unchanged templateGeneration upon semantically equal template update": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 4,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 4,
					Template:           validPodTemplateAbcSemanticallyEqual.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.RollingUpdateDaemonSetStrategyType,
						RollingUpdate: &extensions.RollingUpdateDaemonSet{
							MaxUnavailable: intstr.FromInt(1),
						},
					},
				},
			},
		},
	}
	for testName, successCase := range successCases {
		// ResourceVersion is required for updates.
		successCase.old.ObjectMeta.ResourceVersion = "1"
		successCase.update.ObjectMeta.ResourceVersion = "2"
		// Check test setup
		if successCase.expectedErrNum > 0 {
			t.Errorf("%q has incorrect test setup with expectedErrNum %d, expected no error", testName, successCase.expectedErrNum)
		}
		if len(successCase.old.ObjectMeta.ResourceVersion) == 0 || len(successCase.update.ObjectMeta.ResourceVersion) == 0 {
			t.Errorf("%q has incorrect test setup with no resource version set", testName)
		}
		if errs := ValidateDaemonSetUpdate(&successCase.update, &successCase.old); len(errs) != 0 {
			t.Errorf("%q expected no error, but got: %v", testName, errs)
		}
	}
	errorCases := map[string]dsUpdateTest{
		"change daemon name": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			expectedErrNum: 1,
		},
		"invalid selector": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: invalidSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			expectedErrNum: 1,
		},
		"invalid pod": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 2,
					Template:           invalidPodTemplate.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			expectedErrNum: 1,
		},
		"invalid read-write volume": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 2,
					Template:           readWriteVolumePodTemplate.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			expectedErrNum: 1,
		},
		"invalid update strategy": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: 1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: "Random",
					},
				},
			},
			expectedErrNum: 1,
		},
		"negative templateGeneration": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: -1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					TemplateGeneration: -1,
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			expectedErrNum: 1,
		},
		"decreased templateGeneration": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					TemplateGeneration: 2,
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					TemplateGeneration: 1,
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			expectedErrNum: 1,
		},
		"unchanged templateGeneration upon template update": {
			old: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					TemplateGeneration: 2,
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector},
					Template:           validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					TemplateGeneration: 2,
					Selector:           &metav1.LabelSelector{MatchLabels: validSelector2},
					Template:           validPodTemplateAbc2.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type: extensions.OnDeleteDaemonSetStrategyType,
					},
				},
			},
			expectedErrNum: 1,
		},
	}
	for testName, errorCase := range errorCases {
		// ResourceVersion is required for updates.
		errorCase.old.ObjectMeta.ResourceVersion = "1"
		errorCase.update.ObjectMeta.ResourceVersion = "2"
		// Check test setup
		if errorCase.expectedErrNum <= 0 {
			t.Errorf("%q has incorrect test setup with expectedErrNum %d, expected at least one error", testName, errorCase.expectedErrNum)
		}
		if len(errorCase.old.ObjectMeta.ResourceVersion) == 0 || len(errorCase.update.ObjectMeta.ResourceVersion) == 0 {
			t.Errorf("%q has incorrect test setup with no resource version set", testName)
		}
		// Run the tests
		if errs := ValidateDaemonSetUpdate(&errorCase.update, &errorCase.old); len(errs) != errorCase.expectedErrNum {
			t.Errorf("%q expected %d errors, but got %d error: %v", testName, errorCase.expectedErrNum, len(errs), errs)
		} else {
			t.Logf("(PASS) %q got errors %v", testName, errs)
		}
	}
}

func TestValidateDaemonSet(t *testing.T) {
	validSelector := map[string]string{"a": "b"}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validSelector,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
			},
		},
	}
	invalidSelector := map[string]string{"NoUppercaseOrSpecialCharsLike=Equals": "b"}
	invalidPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
			},
			ObjectMeta: metav1.ObjectMeta{
				Labels: invalidSelector,
			},
		},
	}
	successCases := []extensions.DaemonSet{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type: extensions.OnDeleteDaemonSetStrategyType,
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type: extensions.OnDeleteDaemonSetStrategyType,
				},
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateDaemonSet(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]extensions.DaemonSet{
		"zero-length ID": {
			ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: validPodTemplate.Template,
			},
		},
		"missing-namespace": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123"},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: validPodTemplate.Template,
			},
		},
		"nil selector": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Template: validPodTemplate.Template,
			},
		},
		"empty selector": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{},
				Template: validPodTemplate.Template,
			},
		},
		"selector_doesnt_match": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				Template: validPodTemplate.Template,
			},
		},
		"invalid template": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
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
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: validPodTemplate.Template,
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
			Spec: extensions.DaemonSetSpec{
				Template: invalidPodTemplate.Template,
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
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: validPodTemplate.Template,
			},
		},
		"invalid restart policy 1": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyOnFailure,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: validSelector,
					},
				},
			},
		},
		"invalid restart policy 2": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validSelector},
				Template: api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyNever,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: validSelector,
					},
				},
			},
		},
	}
	for k, v := range errorCases {
		errs := ValidateDaemonSet(&v)
		if len(errs) == 0 {
			t.Errorf("expected failure for %s", k)
		}
		for i := range errs {
			field := errs[i].Field
			if !strings.HasPrefix(field, "spec.template.") &&
				!strings.HasPrefix(field, "spec.updateStrategy") &&
				field != "metadata.name" &&
				field != "metadata.namespace" &&
				field != "spec.selector" &&
				field != "spec.template" &&
				field != "GCEPersistentDisk.ReadOnly" &&
				field != "spec.template.labels" &&
				field != "metadata.annotations" &&
				field != "metadata.labels" {
				t.Errorf("%s: missing prefix for: %v", k, errs[i])
			}
		}
	}
}

func validDeployment() *extensions.Deployment {
	return &extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "abc",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: extensions.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "abc",
				},
			},
			Strategy: extensions.DeploymentStrategy{
				Type: extensions.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &extensions.RollingUpdateDeployment{
					MaxSurge:       intstr.FromInt(1),
					MaxUnavailable: intstr.FromInt(1),
				},
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "abc",
					Namespace: metav1.NamespaceDefault,
					Labels: map[string]string{
						"name": "abc",
					},
				},
				Spec: api.PodSpec{
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSDefault,
					Containers: []api.Container{
						{
							Name:                     "nginx",
							Image:                    "image",
							ImagePullPolicy:          api.PullNever,
							TerminationMessagePolicy: api.TerminationMessageReadFile,
						},
					},
				},
			},
			RollbackTo: &extensions.RollbackConfig{
				Revision: 1,
			},
		},
	}
}

func TestValidateDeployment(t *testing.T) {
	successCases := []*extensions.Deployment{
		validDeployment(),
	}
	for _, successCase := range successCases {
		if errs := ValidateDeployment(successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]*extensions.Deployment{}
	errorCases["metadata.name: Required value"] = &extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
		},
	}
	// selector should match the labels in pod template.
	invalidSelectorDeployment := validDeployment()
	invalidSelectorDeployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"name": "def",
		},
	}
	errorCases["`selector` does not match template `labels`"] = invalidSelectorDeployment

	// RestartPolicy should be always.
	invalidRestartPolicyDeployment := validDeployment()
	invalidRestartPolicyDeployment.Spec.Template.Spec.RestartPolicy = api.RestartPolicyNever
	errorCases["Unsupported value: \"Never\""] = invalidRestartPolicyDeployment

	// must have valid strategy type
	invalidStrategyDeployment := validDeployment()
	invalidStrategyDeployment.Spec.Strategy.Type = extensions.DeploymentStrategyType("randomType")
	errorCases[`supported values: "Recreate", "RollingUpdate"`] = invalidStrategyDeployment

	// rollingUpdate should be nil for recreate.
	invalidRecreateDeployment := validDeployment()
	invalidRecreateDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type:          extensions.RecreateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{},
	}
	errorCases["may not be specified when strategy `type` is 'Recreate'"] = invalidRecreateDeployment

	// MaxSurge should be in the form of 20%.
	invalidMaxSurgeDeployment := validDeployment()
	invalidMaxSurgeDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type: extensions.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{
			MaxSurge: intstr.FromString("20Percent"),
		},
	}
	errorCases["a valid percent string must be"] = invalidMaxSurgeDeployment

	// MaxSurge and MaxUnavailable cannot both be zero.
	invalidRollingUpdateDeployment := validDeployment()
	invalidRollingUpdateDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type: extensions.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{
			MaxSurge:       intstr.FromString("0%"),
			MaxUnavailable: intstr.FromInt(0),
		},
	}
	errorCases["may not be 0 when `maxSurge` is 0"] = invalidRollingUpdateDeployment

	// MaxUnavailable should not be more than 100%.
	invalidMaxUnavailableDeployment := validDeployment()
	invalidMaxUnavailableDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type: extensions.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{
			MaxUnavailable: intstr.FromString("110%"),
		},
	}
	errorCases["must not be greater than 100%"] = invalidMaxUnavailableDeployment

	// Rollback.Revision must be non-negative
	invalidRollbackRevisionDeployment := validDeployment()
	invalidRollbackRevisionDeployment.Spec.RollbackTo.Revision = -3
	errorCases["must be greater than or equal to 0"] = invalidRollbackRevisionDeployment

	// ProgressDeadlineSeconds should be greater than MinReadySeconds
	invalidProgressDeadlineDeployment := validDeployment()
	seconds := int32(600)
	invalidProgressDeadlineDeployment.Spec.ProgressDeadlineSeconds = &seconds
	invalidProgressDeadlineDeployment.Spec.MinReadySeconds = seconds
	errorCases["must be greater than minReadySeconds"] = invalidProgressDeadlineDeployment

	for k, v := range errorCases {
		errs := ValidateDeployment(v)
		if len(errs) == 0 {
			t.Errorf("[%s] expected failure", k)
		} else if !strings.Contains(errs[0].Error(), k) {
			t.Errorf("unexpected error: %q, expected: %q", errs[0].Error(), k)
		}
	}
}

func TestValidateDeploymentStatus(t *testing.T) {
	collisionCount := int32(-3)
	tests := []struct {
		name string

		replicas           int32
		updatedReplicas    int32
		readyReplicas      int32
		availableReplicas  int32
		observedGeneration int64
		collisionCount     *int32

		expectedErr bool
	}{
		{
			name:               "valid status",
			replicas:           3,
			updatedReplicas:    3,
			readyReplicas:      2,
			availableReplicas:  1,
			observedGeneration: 2,
			expectedErr:        false,
		},
		{
			name:               "invalid replicas",
			replicas:           -1,
			updatedReplicas:    2,
			readyReplicas:      2,
			availableReplicas:  1,
			observedGeneration: 2,
			expectedErr:        true,
		},
		{
			name:               "invalid updatedReplicas",
			replicas:           2,
			updatedReplicas:    -1,
			readyReplicas:      2,
			availableReplicas:  1,
			observedGeneration: 2,
			expectedErr:        true,
		},
		{
			name:               "invalid readyReplicas",
			replicas:           3,
			readyReplicas:      -1,
			availableReplicas:  1,
			observedGeneration: 2,
			expectedErr:        true,
		},
		{
			name:               "invalid availableReplicas",
			replicas:           3,
			readyReplicas:      3,
			availableReplicas:  -1,
			observedGeneration: 2,
			expectedErr:        true,
		},
		{
			name:               "invalid observedGeneration",
			replicas:           3,
			readyReplicas:      3,
			availableReplicas:  3,
			observedGeneration: -1,
			expectedErr:        true,
		},
		{
			name:               "updatedReplicas greater than replicas",
			replicas:           3,
			updatedReplicas:    4,
			readyReplicas:      3,
			availableReplicas:  3,
			observedGeneration: 1,
			expectedErr:        true,
		},
		{
			name:               "readyReplicas greater than replicas",
			replicas:           3,
			readyReplicas:      4,
			availableReplicas:  3,
			observedGeneration: 1,
			expectedErr:        true,
		},
		{
			name:               "availableReplicas greater than replicas",
			replicas:           3,
			readyReplicas:      3,
			availableReplicas:  4,
			observedGeneration: 1,
			expectedErr:        true,
		},
		{
			name:               "availableReplicas greater than readyReplicas",
			replicas:           3,
			readyReplicas:      2,
			availableReplicas:  3,
			observedGeneration: 1,
			expectedErr:        true,
		},
		{
			name:               "invalid collisionCount",
			replicas:           3,
			observedGeneration: 1,
			collisionCount:     &collisionCount,
			expectedErr:        true,
		},
	}

	for _, test := range tests {
		status := extensions.DeploymentStatus{
			Replicas:           test.replicas,
			UpdatedReplicas:    test.updatedReplicas,
			ReadyReplicas:      test.readyReplicas,
			AvailableReplicas:  test.availableReplicas,
			ObservedGeneration: test.observedGeneration,
			CollisionCount:     test.collisionCount,
		}

		errs := ValidateDeploymentStatus(&status, field.NewPath("status"))
		if hasErr := len(errs) > 0; hasErr != test.expectedErr {
			errString := spew.Sprintf("%#v", errs)
			t.Errorf("%s: expected error: %t, got error: %t\nerrors: %s", test.name, test.expectedErr, hasErr, errString)
		}
	}
}

func TestValidateDeploymentStatusUpdate(t *testing.T) {
	collisionCount := int32(1)
	otherCollisionCount := int32(2)
	tests := []struct {
		name string

		from, to extensions.DeploymentStatus

		expectedErr bool
	}{
		{
			name: "increase: valid update",
			from: extensions.DeploymentStatus{
				CollisionCount: nil,
			},
			to: extensions.DeploymentStatus{
				CollisionCount: &collisionCount,
			},
			expectedErr: false,
		},
		{
			name: "stable: valid update",
			from: extensions.DeploymentStatus{
				CollisionCount: &collisionCount,
			},
			to: extensions.DeploymentStatus{
				CollisionCount: &collisionCount,
			},
			expectedErr: false,
		},
		{
			name: "unset: invalid update",
			from: extensions.DeploymentStatus{
				CollisionCount: &collisionCount,
			},
			to: extensions.DeploymentStatus{
				CollisionCount: nil,
			},
			expectedErr: true,
		},
		{
			name: "decrease: invalid update",
			from: extensions.DeploymentStatus{
				CollisionCount: &otherCollisionCount,
			},
			to: extensions.DeploymentStatus{
				CollisionCount: &collisionCount,
			},
			expectedErr: true,
		},
	}

	for _, test := range tests {
		meta := metav1.ObjectMeta{Name: "foo", Namespace: metav1.NamespaceDefault, ResourceVersion: "1"}
		from := &extensions.Deployment{
			ObjectMeta: meta,
			Status:     test.from,
		}
		to := &extensions.Deployment{
			ObjectMeta: meta,
			Status:     test.to,
		}

		errs := ValidateDeploymentStatusUpdate(to, from)
		if hasErr := len(errs) > 0; hasErr != test.expectedErr {
			errString := spew.Sprintf("%#v", errs)
			t.Errorf("%s: expected error: %t, got error: %t\nerrors: %s", test.name, test.expectedErr, hasErr, errString)
		}
	}
}

func validDeploymentRollback() *extensions.DeploymentRollback {
	return &extensions.DeploymentRollback{
		Name: "abc",
		UpdatedAnnotations: map[string]string{
			"created-by": "abc",
		},
		RollbackTo: extensions.RollbackConfig{
			Revision: 1,
		},
	}
}

func TestValidateDeploymentRollback(t *testing.T) {
	noAnnotation := validDeploymentRollback()
	noAnnotation.UpdatedAnnotations = nil
	successCases := []*extensions.DeploymentRollback{
		validDeploymentRollback(),
		noAnnotation,
	}
	for _, successCase := range successCases {
		if errs := ValidateDeploymentRollback(successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]*extensions.DeploymentRollback{}
	invalidNoName := validDeploymentRollback()
	invalidNoName.Name = ""
	errorCases["name: Required value"] = invalidNoName

	for k, v := range errorCases {
		errs := ValidateDeploymentRollback(v)
		if len(errs) == 0 {
			t.Errorf("[%s] expected failure", k)
		} else if !strings.Contains(errs[0].Error(), k) {
			t.Errorf("unexpected error: %q, expected: %q", errs[0].Error(), k)
		}
	}
}

func TestValidateIngress(t *testing.T) {
	defaultBackend := extensions.IngressBackend{
		ServiceName: "default-backend",
		ServicePort: intstr.FromInt(80),
	}

	newValid := func() extensions.Ingress {
		return extensions.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: extensions.IngressSpec{
				Backend: &extensions.IngressBackend{
					ServiceName: "default-backend",
					ServicePort: intstr.FromInt(80),
				},
				Rules: []extensions.IngressRule{
					{
						Host: "foo.bar.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path:    "/foo",
										Backend: defaultBackend,
									},
								},
							},
						},
					},
				},
			},
			Status: extensions.IngressStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{IP: "127.0.0.1"},
					},
				},
			},
		}
	}
	servicelessBackend := newValid()
	servicelessBackend.Spec.Backend.ServiceName = ""
	invalidNameBackend := newValid()
	invalidNameBackend.Spec.Backend.ServiceName = "defaultBackend"
	noPortBackend := newValid()
	noPortBackend.Spec.Backend = &extensions.IngressBackend{ServiceName: defaultBackend.ServiceName}
	noForwardSlashPath := newValid()
	noForwardSlashPath.Spec.Rules[0].IngressRuleValue.HTTP.Paths = []extensions.HTTPIngressPath{
		{
			Path:    "invalid",
			Backend: defaultBackend,
		},
	}
	noPaths := newValid()
	noPaths.Spec.Rules[0].IngressRuleValue.HTTP.Paths = []extensions.HTTPIngressPath{}
	badHost := newValid()
	badHost.Spec.Rules[0].Host = "foobar:80"
	badRegexPath := newValid()
	badPathExpr := "/invalid["
	badRegexPath.Spec.Rules[0].IngressRuleValue.HTTP.Paths = []extensions.HTTPIngressPath{
		{
			Path:    badPathExpr,
			Backend: defaultBackend,
		},
	}
	badPathErr := fmt.Sprintf("spec.rules[0].http.paths[0].path: Invalid value: '%v'", badPathExpr)
	hostIP := "127.0.0.1"
	badHostIP := newValid()
	badHostIP.Spec.Rules[0].Host = hostIP
	badHostIPErr := fmt.Sprintf("spec.rules[0].host: Invalid value: '%v'", hostIP)

	errorCases := map[string]extensions.Ingress{
		"spec.backend.serviceName: Required value":        servicelessBackend,
		"spec.backend.serviceName: Invalid value":         invalidNameBackend,
		"spec.backend.servicePort: Invalid value":         noPortBackend,
		"spec.rules[0].host: Invalid value":               badHost,
		"spec.rules[0].http.paths: Required value":        noPaths,
		"spec.rules[0].http.paths[0].path: Invalid value": noForwardSlashPath,
	}
	errorCases[badPathErr] = badRegexPath
	errorCases[badHostIPErr] = badHostIP

	wildcardHost := "foo.*.bar.com"
	badWildcard := newValid()
	badWildcard.Spec.Rules[0].Host = wildcardHost
	badWildcardErr := fmt.Sprintf("spec.rules[0].host: Invalid value: '%v'", wildcardHost)
	errorCases[badWildcardErr] = badWildcard

	for k, v := range errorCases {
		errs := ValidateIngress(&v)
		if len(errs) == 0 {
			t.Errorf("expected failure for %q", k)
		} else {
			s := strings.Split(k, ":")
			err := errs[0]
			if err.Field != s[0] || !strings.Contains(err.Error(), s[1]) {
				t.Errorf("unexpected error: %q, expected: %q", err, k)
			}
		}
	}
}

func TestValidateIngressTLS(t *testing.T) {
	defaultBackend := extensions.IngressBackend{
		ServiceName: "default-backend",
		ServicePort: intstr.FromInt(80),
	}

	newValid := func() extensions.Ingress {
		return extensions.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: extensions.IngressSpec{
				Backend: &extensions.IngressBackend{
					ServiceName: "default-backend",
					ServicePort: intstr.FromInt(80),
				},
				Rules: []extensions.IngressRule{
					{
						Host: "foo.bar.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path:    "/foo",
										Backend: defaultBackend,
									},
								},
							},
						},
					},
				},
			},
			Status: extensions.IngressStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{IP: "127.0.0.1"},
					},
				},
			},
		}
	}

	errorCases := map[string]extensions.Ingress{}

	wildcardHost := "foo.*.bar.com"
	badWildcardTLS := newValid()
	badWildcardTLS.Spec.Rules[0].Host = "*.foo.bar.com"
	badWildcardTLS.Spec.TLS = []extensions.IngressTLS{
		{
			Hosts: []string{wildcardHost},
		},
	}
	badWildcardTLSErr := fmt.Sprintf("spec.tls[0].hosts: Invalid value: '%v'", wildcardHost)
	errorCases[badWildcardTLSErr] = badWildcardTLS

	for k, v := range errorCases {
		errs := ValidateIngress(&v)
		if len(errs) == 0 {
			t.Errorf("expected failure for %q", k)
		} else {
			s := strings.Split(k, ":")
			err := errs[0]
			if err.Field != s[0] || !strings.Contains(err.Error(), s[1]) {
				t.Errorf("unexpected error: %q, expected: %q", err, k)
			}
		}
	}
}

func TestValidateIngressStatusUpdate(t *testing.T) {
	defaultBackend := extensions.IngressBackend{
		ServiceName: "default-backend",
		ServicePort: intstr.FromInt(80),
	}

	newValid := func() extensions.Ingress {
		return extensions.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "foo",
				Namespace:       metav1.NamespaceDefault,
				ResourceVersion: "9",
			},
			Spec: extensions.IngressSpec{
				Backend: &extensions.IngressBackend{
					ServiceName: "default-backend",
					ServicePort: intstr.FromInt(80),
				},
				Rules: []extensions.IngressRule{
					{
						Host: "foo.bar.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path:    "/foo",
										Backend: defaultBackend,
									},
								},
							},
						},
					},
				},
			},
			Status: extensions.IngressStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{IP: "127.0.0.1", Hostname: "foo.bar.com"},
					},
				},
			},
		}
	}
	oldValue := newValid()
	newValue := newValid()
	newValue.Status = extensions.IngressStatus{
		LoadBalancer: api.LoadBalancerStatus{
			Ingress: []api.LoadBalancerIngress{
				{IP: "127.0.0.2", Hostname: "foo.com"},
			},
		},
	}
	invalidIP := newValid()
	invalidIP.Status = extensions.IngressStatus{
		LoadBalancer: api.LoadBalancerStatus{
			Ingress: []api.LoadBalancerIngress{
				{IP: "abcd", Hostname: "foo.com"},
			},
		},
	}
	invalidHostname := newValid()
	invalidHostname.Status = extensions.IngressStatus{
		LoadBalancer: api.LoadBalancerStatus{
			Ingress: []api.LoadBalancerIngress{
				{IP: "127.0.0.1", Hostname: "127.0.0.1"},
			},
		},
	}

	errs := ValidateIngressStatusUpdate(&newValue, &oldValue)
	if len(errs) != 0 {
		t.Errorf("Unexpected error %v", errs)
	}

	errorCases := map[string]extensions.Ingress{
		"status.loadBalancer.ingress[0].ip: Invalid value":       invalidIP,
		"status.loadBalancer.ingress[0].hostname: Invalid value": invalidHostname,
	}
	for k, v := range errorCases {
		errs := ValidateIngressStatusUpdate(&v, &oldValue)
		if len(errs) == 0 {
			t.Errorf("expected failure for %s", k)
		} else {
			s := strings.Split(k, ":")
			err := errs[0]
			if err.Field != s[0] || !strings.Contains(err.Error(), s[1]) {
				t.Errorf("unexpected error: %q, expected: %q", err, k)
			}
		}
	}
}

func TestValidateReplicaSetStatus(t *testing.T) {
	tests := []struct {
		name string

		replicas             int32
		fullyLabeledReplicas int32
		readyReplicas        int32
		availableReplicas    int32
		observedGeneration   int64

		expectedErr bool
	}{
		{
			name:                 "valid status",
			replicas:             3,
			fullyLabeledReplicas: 3,
			readyReplicas:        2,
			availableReplicas:    1,
			observedGeneration:   2,
			expectedErr:          false,
		},
		{
			name:                 "invalid replicas",
			replicas:             -1,
			fullyLabeledReplicas: 3,
			readyReplicas:        2,
			availableReplicas:    1,
			observedGeneration:   2,
			expectedErr:          true,
		},
		{
			name:                 "invalid fullyLabeledReplicas",
			replicas:             3,
			fullyLabeledReplicas: -1,
			readyReplicas:        2,
			availableReplicas:    1,
			observedGeneration:   2,
			expectedErr:          true,
		},
		{
			name:                 "invalid readyReplicas",
			replicas:             3,
			fullyLabeledReplicas: 3,
			readyReplicas:        -1,
			availableReplicas:    1,
			observedGeneration:   2,
			expectedErr:          true,
		},
		{
			name:                 "invalid availableReplicas",
			replicas:             3,
			fullyLabeledReplicas: 3,
			readyReplicas:        3,
			availableReplicas:    -1,
			observedGeneration:   2,
			expectedErr:          true,
		},
		{
			name:                 "invalid observedGeneration",
			replicas:             3,
			fullyLabeledReplicas: 3,
			readyReplicas:        3,
			availableReplicas:    3,
			observedGeneration:   -1,
			expectedErr:          true,
		},
		{
			name:                 "fullyLabeledReplicas greater than replicas",
			replicas:             3,
			fullyLabeledReplicas: 4,
			readyReplicas:        3,
			availableReplicas:    3,
			observedGeneration:   1,
			expectedErr:          true,
		},
		{
			name:                 "readyReplicas greater than replicas",
			replicas:             3,
			fullyLabeledReplicas: 3,
			readyReplicas:        4,
			availableReplicas:    3,
			observedGeneration:   1,
			expectedErr:          true,
		},
		{
			name:                 "availableReplicas greater than replicas",
			replicas:             3,
			fullyLabeledReplicas: 3,
			readyReplicas:        3,
			availableReplicas:    4,
			observedGeneration:   1,
			expectedErr:          true,
		},
		{
			name:                 "availableReplicas greater than readyReplicas",
			replicas:             3,
			fullyLabeledReplicas: 3,
			readyReplicas:        2,
			availableReplicas:    3,
			observedGeneration:   1,
			expectedErr:          true,
		},
	}

	for _, test := range tests {
		status := extensions.ReplicaSetStatus{
			Replicas:             test.replicas,
			FullyLabeledReplicas: test.fullyLabeledReplicas,
			ReadyReplicas:        test.readyReplicas,
			AvailableReplicas:    test.availableReplicas,
			ObservedGeneration:   test.observedGeneration,
		}

		if hasErr := len(ValidateReplicaSetStatus(status, field.NewPath("status"))) > 0; hasErr != test.expectedErr {
			t.Errorf("%s: expected error: %t, got error: %t", test.name, test.expectedErr, hasErr)
		}
	}
}

func TestValidateReplicaSetStatusUpdate(t *testing.T) {
	validLabels := map[string]string{"a": "b"}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
			},
		},
	}
	type rcUpdateTest struct {
		old    extensions.ReplicaSet
		update extensions.ReplicaSet
	}
	successCases := []rcUpdateTest{
		{
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
				Status: extensions.ReplicaSetStatus{
					Replicas: 2,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 3,
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
				Status: extensions.ReplicaSetStatus{
					Replicas: 4,
				},
			},
		},
	}
	for _, successCase := range successCases {
		successCase.old.ObjectMeta.ResourceVersion = "1"
		successCase.update.ObjectMeta.ResourceVersion = "1"
		if errs := ValidateReplicaSetStatusUpdate(&successCase.update, &successCase.old); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}
	errorCases := map[string]rcUpdateTest{
		"negative replicas": {
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
				Status: extensions.ReplicaSetStatus{
					Replicas: 3,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 2,
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
				Status: extensions.ReplicaSetStatus{
					Replicas: -3,
				},
			},
		},
	}
	for testName, errorCase := range errorCases {
		if errs := ValidateReplicaSetStatusUpdate(&errorCase.update, &errorCase.old); len(errs) == 0 {
			t.Errorf("expected failure: %s", testName)
		}
	}

}

func TestValidateReplicaSetUpdate(t *testing.T) {
	validLabels := map[string]string{"a": "b"}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
			},
		},
	}
	readWriteVolumePodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
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
	type rcUpdateTest struct {
		old    extensions.ReplicaSet
		update extensions.ReplicaSet
	}
	successCases := []rcUpdateTest{
		{
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 3,
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
		},
		{
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 1,
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: readWriteVolumePodTemplate.Template,
				},
			},
		},
	}
	for _, successCase := range successCases {
		successCase.old.ObjectMeta.ResourceVersion = "1"
		successCase.update.ObjectMeta.ResourceVersion = "1"
		if errs := ValidateReplicaSetUpdate(&successCase.update, &successCase.old); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}
	errorCases := map[string]rcUpdateTest{
		"more than one read/write": {
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 2,
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: readWriteVolumePodTemplate.Template,
				},
			},
		},
		"invalid selector": {
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 2,
					Selector: &metav1.LabelSelector{MatchLabels: invalidLabels},
					Template: validPodTemplate.Template,
				},
			},
		},
		"invalid pod": {
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 2,
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: invalidPodTemplate.Template,
				},
			},
		},
		"negative replicas": {
			old: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
			update: extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
				Spec: extensions.ReplicaSetSpec{
					Replicas: -1,
					Selector: &metav1.LabelSelector{MatchLabels: validLabels},
					Template: validPodTemplate.Template,
				},
			},
		},
	}
	for testName, errorCase := range errorCases {
		if errs := ValidateReplicaSetUpdate(&errorCase.update, &errorCase.old); len(errs) == 0 {
			t.Errorf("expected failure: %s", testName)
		}
	}
}

func TestValidateReplicaSet(t *testing.T) {
	validLabels := map[string]string{"a": "b"}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
			},
		},
	}
	readWriteVolumePodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: validLabels,
			},
			Spec: api.PodSpec{
				Volumes:       []api.Volume{{Name: "gcepd", VolumeSource: api.VolumeSource{GCEPersistentDisk: &api.GCEPersistentDiskVolumeSource{PDName: "my-PD", FSType: "ext4", Partition: 1, ReadOnly: false}}}},
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
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
	successCases := []extensions.ReplicaSet{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: validPodTemplate.Template,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: validPodTemplate.Template,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Replicas: 1,
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: readWriteVolumePodTemplate.Template,
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateReplicaSet(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]extensions.ReplicaSet{
		"zero-length ID": {
			ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: validPodTemplate.Template,
			},
		},
		"missing-namespace": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc-123"},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: validPodTemplate.Template,
			},
		},
		"empty selector": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Template: validPodTemplate.Template,
			},
		},
		"selector_doesnt_match": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				Template: validPodTemplate.Template,
			},
		},
		"invalid manifest": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
			},
		},
		"read-write persistent disk with > 1 pod": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc"},
			Spec: extensions.ReplicaSetSpec{
				Replicas: 2,
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: readWriteVolumePodTemplate.Template,
			},
		},
		"negative_replicas": {
			ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
			Spec: extensions.ReplicaSetSpec{
				Replicas: -1,
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
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
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: validPodTemplate.Template,
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
			Spec: extensions.ReplicaSetSpec{
				Template: invalidPodTemplate.Template,
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
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: validPodTemplate.Template,
			},
		},
		"invalid restart policy 1": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyOnFailure,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: validLabels,
					},
				},
			},
		},
		"invalid restart policy 2": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-123",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: extensions.ReplicaSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: validLabels},
				Template: api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyNever,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent", TerminationMessagePolicy: api.TerminationMessageReadFile}},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: validLabels,
					},
				},
			},
		},
	}
	for k, v := range errorCases {
		errs := ValidateReplicaSet(&v)
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
				field != "status.replicas" {
				t.Errorf("%s: missing prefix for: %v", k, errs[i])
			}
		}
	}
}
