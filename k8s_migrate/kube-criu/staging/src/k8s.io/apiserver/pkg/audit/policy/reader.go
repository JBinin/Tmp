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

package policy

import (
	"fmt"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/runtime/schema"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
	auditv1alpha1 "k8s.io/apiserver/pkg/apis/audit/v1alpha1"
	auditv1beta1 "k8s.io/apiserver/pkg/apis/audit/v1beta1"
	"k8s.io/apiserver/pkg/apis/audit/validation"
	"k8s.io/apiserver/pkg/audit"

	"github.com/golang/glog"
)

var (
	apiGroupVersions = []schema.GroupVersion{
		auditv1beta1.SchemeGroupVersion,
		auditv1alpha1.SchemeGroupVersion,
		auditv1.SchemeGroupVersion,
	}
	apiGroupVersionSet = map[schema.GroupVersion]bool{}
)

func init() {
	for _, gv := range apiGroupVersions {
		apiGroupVersionSet[gv] = true
	}
}

func LoadPolicyFromFile(filePath string) (*auditinternal.Policy, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path not specified")
	}
	policyDef, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file path %q: %+v", filePath, err)
	}

	policy := &auditinternal.Policy{}
	decoder := audit.Codecs.UniversalDecoder(apiGroupVersions...)

	_, gvk, err := decoder.Decode(policyDef, nil, policy)
	if err != nil {
		return nil, fmt.Errorf("failed decoding file %q: %v", filePath, err)
	}

	// Ensure the policy file contained an apiVersion and kind.
	if !apiGroupVersionSet[schema.GroupVersion{Group: gvk.Group, Version: gvk.Version}] {
		return nil, fmt.Errorf("unknown group version field %v in policy file %s", gvk, filePath)
	}

	if err := validation.ValidatePolicy(policy); err != nil {
		return nil, err.ToAggregate()
	}

	policyCnt := len(policy.Rules)
	if policyCnt == 0 {
		return nil, fmt.Errorf("loaded illegal policy with 0 rules from file %s", filePath)
	}
	glog.V(4).Infof("Loaded %d audit policy rules from file %s", policyCnt, filePath)
	return policy, nil
}
