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

package algorithmprovider

import (
	"fmt"
	"testing"

	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/scheduler/factory"
)

var (
	algorithmProviderNames = []string{
		factory.DefaultProvider,
	}
)

func TestDefaultConfigExists(t *testing.T) {
	p, err := factory.GetAlgorithmProvider(factory.DefaultProvider)
	if err != nil {
		t.Errorf("error retrieving default provider: %v", err)
	}
	if p == nil {
		t.Error("algorithm provider config should not be nil")
	}
	if len(p.FitPredicateKeys) == 0 {
		t.Error("default algorithm provider shouldn't have 0 fit predicates")
	}
}

func TestAlgorithmProviders(t *testing.T) {
	for _, pn := range algorithmProviderNames {
		t.Run(pn, func(t *testing.T) {
			p, err := factory.GetAlgorithmProvider(pn)
			if err != nil {
				t.Fatalf("error retrieving provider: %v", err)
			}
			if len(p.PriorityFunctionKeys) == 0 {
				t.Errorf("algorithm provider shouldn't have 0 priority functions")
			}
			for _, pf := range p.PriorityFunctionKeys.List() {
				t.Run(fmt.Sprintf("priorityfunction/%s", pf), func(t *testing.T) {
					if !factory.IsPriorityFunctionRegistered(pf) {
						t.Errorf("priority function is not registered but is used in the algorithm provider")
					}
				})
			}
			for _, fp := range p.FitPredicateKeys.List() {
				t.Run(fmt.Sprintf("fitpredicate/%s", fp), func(t *testing.T) {
					if !factory.IsFitPredicateRegistered(fp) {
						t.Errorf("fit predicate is not registered but is used in the algorithm provider")
					}
				})
			}
		})
	}
}

func TestApplyFeatureGates(t *testing.T) {
	for _, pn := range algorithmProviderNames {
		t.Run(pn, func(t *testing.T) {
			p, err := factory.GetAlgorithmProvider(pn)
			if err != nil {
				t.Fatalf("Error retrieving provider: %v", err)
			}

			if !p.FitPredicateKeys.Has("CheckNodeCondition") {
				t.Fatalf("Failed to find predicate: 'CheckNodeCondition'")
			}

			if !p.FitPredicateKeys.Has("PodToleratesNodeTaints") {
				t.Fatalf("Failed to find predicate: 'PodToleratesNodeTaints'")
			}
		})
	}

	// Apply features for algorithm providers.
	utilfeature.DefaultFeatureGate.Set("TaintNodesByCondition=True")

	ApplyFeatureGates()

	for _, pn := range algorithmProviderNames {
		t.Run(pn, func(t *testing.T) {
			p, err := factory.GetAlgorithmProvider(pn)
			if err != nil {
				t.Fatalf("Error retrieving '%s' provider: %v", pn, err)
			}

			if !p.FitPredicateKeys.Has("PodToleratesNodeTaints") {
				t.Fatalf("Failed to find predicate: 'PodToleratesNodeTaints'")
			}

			if p.FitPredicateKeys.Has("CheckNodeCondition") {
				t.Fatalf("Unexpected predicate: 'CheckNodeCondition'")
			}
		})
	}
}
