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

package defaults

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
)

func TestCopyAndReplace(t *testing.T) {
	testCases := []struct {
		set         sets.String
		replaceWhat string
		replaceWith string
		expected    sets.String
	}{
		{
			set:         sets.String{"A": sets.Empty{}, "B": sets.Empty{}},
			replaceWhat: "A",
			replaceWith: "C",
			expected:    sets.String{"B": sets.Empty{}, "C": sets.Empty{}},
		},
		{
			set:         sets.String{"A": sets.Empty{}, "B": sets.Empty{}},
			replaceWhat: "D",
			replaceWith: "C",
			expected:    sets.String{"A": sets.Empty{}, "B": sets.Empty{}},
		},
	}
	for _, testCase := range testCases {
		result := copyAndReplace(testCase.set, testCase.replaceWhat, testCase.replaceWith)
		if !result.Equal(testCase.expected) {
			t.Errorf("expected %v got %v", testCase.expected, result)
		}
	}
}

func TestDefaultPriorities(t *testing.T) {
	result := sets.NewString(
		"SelectorSpreadPriority",
		"InterPodAffinityPriority",
		"LeastRequestedPriority",
		"BalancedResourceAllocation",
		"NodePreferAvoidPodsPriority",
		"NodeAffinityPriority",
		"TaintTolerationPriority")
	if expected := defaultPriorities(); !result.Equal(expected) {
		t.Errorf("expected %v got %v", expected, result)
	}
}

func TestDefaultPredicates(t *testing.T) {
	result := sets.NewString(
		predicates.NoVolumeZoneConflictPred,
		predicates.MaxEBSVolumeCountPred,
		predicates.MaxGCEPDVolumeCountPred,
		predicates.MaxAzureDiskVolumeCountPred,
		predicates.MatchInterPodAffinityPred,
		predicates.NoDiskConflictPred,
		predicates.GeneralPred,
		predicates.CheckNodeMemoryPressurePred,
		predicates.CheckNodeDiskPressurePred,
		predicates.CheckNodePIDPressurePred,
		predicates.CheckNodeConditionPred,
		predicates.PodToleratesNodeTaintsPred,
		predicates.CheckVolumeBindingPred,
	)

	if expected := defaultPredicates(); !result.Equal(expected) {
		t.Errorf("expected %v got %v", expected, result)
	}
}
