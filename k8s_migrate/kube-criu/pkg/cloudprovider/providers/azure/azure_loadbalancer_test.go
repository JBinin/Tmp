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

package azure

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
)

func TestFindProbe(t *testing.T) {
	tests := []struct {
		msg           string
		existingProbe []network.Probe
		curProbe      network.Probe
		expected      bool
	}{
		{
			msg:      "empty existing probes should return false",
			expected: false,
		},
		{
			msg: "probe names match while ports unmatch should return false",
			existingProbe: []network.Probe{
				{
					Name: to.StringPtr("httpProbe"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port: to.Int32Ptr(1),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.StringPtr("httpProbe"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port: to.Int32Ptr(2),
				},
			},
			expected: false,
		},
		{
			msg: "probe ports match while names unmatch should return false",
			existingProbe: []network.Probe{
				{
					Name: to.StringPtr("probe1"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port: to.Int32Ptr(1),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.StringPtr("probe2"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port: to.Int32Ptr(1),
				},
			},
			expected: false,
		},
		{
			msg: "both probe ports and names match should return true",
			existingProbe: []network.Probe{
				{
					Name: to.StringPtr("matchName"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port: to.Int32Ptr(1),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.StringPtr("matchName"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port: to.Int32Ptr(1),
				},
			},
			expected: true,
		},
	}

	for i, test := range tests {
		findResult := findProbe(test.existingProbe, test.curProbe)
		assert.Equal(t, test.expected, findResult, fmt.Sprintf("TestCase[%d]: %s", i, test.msg))
	}
}

func TestFindRule(t *testing.T) {
	tests := []struct {
		msg          string
		existingRule []network.LoadBalancingRule
		curRule      network.LoadBalancingRule
		expected     bool
	}{
		{
			msg:      "empty existing rules should return false",
			expected: false,
		},
		{
			msg: "rule names unmatch should return false",
			existingRule: []network.LoadBalancingRule{
				{
					Name: to.StringPtr("httpProbe1"),
					LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
						FrontendPort: to.Int32Ptr(1),
					},
				},
			},
			curRule: network.LoadBalancingRule{
				Name: to.StringPtr("httpProbe2"),
				LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
					FrontendPort: to.Int32Ptr(1),
				},
			},
			expected: false,
		},
		{
			msg: "rule names match while frontend ports unmatch should return false",
			existingRule: []network.LoadBalancingRule{
				{
					Name: to.StringPtr("httpProbe"),
					LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
						FrontendPort: to.Int32Ptr(1),
					},
				},
			},
			curRule: network.LoadBalancingRule{
				Name: to.StringPtr("httpProbe"),
				LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
					FrontendPort: to.Int32Ptr(2),
				},
			},
			expected: false,
		},
		{
			msg: "rule names match while backend ports unmatch should return false",
			existingRule: []network.LoadBalancingRule{
				{
					Name: to.StringPtr("httpProbe"),
					LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
						BackendPort: to.Int32Ptr(1),
					},
				},
			},
			curRule: network.LoadBalancingRule{
				Name: to.StringPtr("httpProbe"),
				LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
					BackendPort: to.Int32Ptr(2),
				},
			},
			expected: false,
		},
		{
			msg: "rule names match while LoadDistribution unmatch should return false",
			existingRule: []network.LoadBalancingRule{
				{
					Name: to.StringPtr("probe1"),
					LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
						LoadDistribution: network.Default,
					},
				},
			},
			curRule: network.LoadBalancingRule{
				Name: to.StringPtr("probe2"),
				LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
					LoadDistribution: network.SourceIP,
				},
			},
			expected: false,
		},
		{
			msg: "both rule names and LoadBalancingRulePropertiesFormats match should return true",
			existingRule: []network.LoadBalancingRule{
				{
					Name: to.StringPtr("matchName"),
					LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
						BackendPort:      to.Int32Ptr(2),
						FrontendPort:     to.Int32Ptr(2),
						LoadDistribution: network.SourceIP,
					},
				},
			},
			curRule: network.LoadBalancingRule{
				Name: to.StringPtr("matchName"),
				LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
					BackendPort:      to.Int32Ptr(2),
					FrontendPort:     to.Int32Ptr(2),
					LoadDistribution: network.SourceIP,
				},
			},
			expected: true,
		},
	}

	for i, test := range tests {
		findResult := findRule(test.existingRule, test.curRule)
		assert.Equal(t, test.expected, findResult, fmt.Sprintf("TestCase[%d]: %s", i, test.msg))
	}
}

func TestGetIdleTimeout(t *testing.T) {
	for _, c := range []struct {
		desc        string
		annotations map[string]string
		i           *int32
		err         bool
	}{
		{desc: "no annotation"},
		{desc: "annotation empty value", annotations: map[string]string{ServiceAnnotationLoadBalancerIdleTimeout: ""}, err: true},
		{desc: "annotation not a number", annotations: map[string]string{ServiceAnnotationLoadBalancerIdleTimeout: "cookies"}, err: true},
		{desc: "annotation negative value", annotations: map[string]string{ServiceAnnotationLoadBalancerIdleTimeout: "-6"}, err: true},
		{desc: "annotation zero value", annotations: map[string]string{ServiceAnnotationLoadBalancerIdleTimeout: "0"}, err: true},
		{desc: "annotation too low value", annotations: map[string]string{ServiceAnnotationLoadBalancerIdleTimeout: "3"}, err: true},
		{desc: "annotation too high value", annotations: map[string]string{ServiceAnnotationLoadBalancerIdleTimeout: "31"}, err: true},
		{desc: "annotation good value", annotations: map[string]string{ServiceAnnotationLoadBalancerIdleTimeout: "24"}, i: to.Int32Ptr(24)},
	} {
		t.Run(c.desc, func(t *testing.T) {
			s := &v1.Service{}
			s.Annotations = c.annotations
			i, err := getIdleTimeout(s)

			if !reflect.DeepEqual(c.i, i) {
				t.Fatalf("got unexpected value: %d", to.Int32(i))
			}
			if (err != nil) != c.err {
				t.Fatalf("expected error=%v, got %v", c.err, err)
			}
		})
	}
}
