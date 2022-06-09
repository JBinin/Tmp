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

package core

import "testing"

func TestTaintToString(t *testing.T) {
	testCases := []struct {
		taint          *Taint
		expectedString string
	}{
		{
			taint: &Taint{
				Key:    "foo",
				Value:  "bar",
				Effect: TaintEffectNoSchedule,
			},
			expectedString: "foo=bar:NoSchedule",
		},
		{
			taint: &Taint{
				Key:    "foo",
				Effect: TaintEffectNoSchedule,
			},
			expectedString: "foo:NoSchedule",
		},
	}

	for i, tc := range testCases {
		if tc.expectedString != tc.taint.ToString() {
			t.Errorf("[%v] expected taint %v converted to %s, got %s", i, tc.taint, tc.expectedString, tc.taint.ToString())
		}
	}
}

func TestMatchTaint(t *testing.T) {
	testCases := []struct {
		description  string
		taint        *Taint
		taintToMatch Taint
		expectMatch  bool
	}{
		{
			description: "two taints with the same key,value,effect should match",
			taint: &Taint{
				Key:    "foo",
				Value:  "bar",
				Effect: TaintEffectNoSchedule,
			},
			taintToMatch: Taint{
				Key:    "foo",
				Value:  "bar",
				Effect: TaintEffectNoSchedule,
			},
			expectMatch: true,
		},
		{
			description: "two taints with the same key,effect but different value should match",
			taint: &Taint{
				Key:    "foo",
				Value:  "bar",
				Effect: TaintEffectNoSchedule,
			},
			taintToMatch: Taint{
				Key:    "foo",
				Value:  "different-value",
				Effect: TaintEffectNoSchedule,
			},
			expectMatch: true,
		},
		{
			description: "two taints with the different key cannot match",
			taint: &Taint{
				Key:    "foo",
				Value:  "bar",
				Effect: TaintEffectNoSchedule,
			},
			taintToMatch: Taint{
				Key:    "different-key",
				Value:  "bar",
				Effect: TaintEffectNoSchedule,
			},
			expectMatch: false,
		},
		{
			description: "two taints with the different effect cannot match",
			taint: &Taint{
				Key:    "foo",
				Value:  "bar",
				Effect: TaintEffectNoSchedule,
			},
			taintToMatch: Taint{
				Key:    "foo",
				Value:  "bar",
				Effect: TaintEffectPreferNoSchedule,
			},
			expectMatch: false,
		},
	}

	for _, tc := range testCases {
		if tc.expectMatch != tc.taint.MatchTaint(tc.taintToMatch) {
			t.Errorf("[%s] expect taint %s match taint %s", tc.description, tc.taint.ToString(), tc.taintToMatch.ToString())
		}
	}
}
