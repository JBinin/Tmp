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

package gce

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce/cloud"

	"github.com/stretchr/testify/assert"
)

func TestServiceNetworkTierAnnotationKey(t *testing.T) {
	createTestService := func() *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				UID:       "randome-uid",
				Name:      "test-svc",
				Namespace: "test-ns",
			},
		}
	}

	for testName, testCase := range map[string]struct {
		annotations  map[string]string
		expectedTier cloud.NetworkTier
		expectErr    bool
	}{
		"Use the default when the annotation does not exist": {
			annotations:  nil,
			expectedTier: cloud.NetworkTierDefault,
		},
		"Standard tier": {
			annotations:  map[string]string{NetworkTierAnnotationKey: "Standard"},
			expectedTier: cloud.NetworkTierStandard,
		},
		"Premium tier": {
			annotations:  map[string]string{NetworkTierAnnotationKey: "Premium"},
			expectedTier: cloud.NetworkTierPremium,
		},
		"Report an error on invalid network tier value": {
			annotations:  map[string]string{NetworkTierAnnotationKey: "Unknown-tier"},
			expectedTier: cloud.NetworkTierPremium,
			expectErr:    true,
		},
	} {
		t.Run(testName, func(t *testing.T) {
			svc := createTestService()
			svc.Annotations = testCase.annotations
			actualTier, err := GetServiceNetworkTier(svc)
			assert.Equal(t, testCase.expectedTier, actualTier)
			assert.Equal(t, testCase.expectErr, err != nil)
		})
	}
}
