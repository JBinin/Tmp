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

package explain

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
)

func TestSplitAndParseResourceRequest(t *testing.T) {
	tests := []struct {
		name       string
		inresource string

		expectedInResource string
		expectedFieldsPath []string
		expectedErr        bool
	}{
		{
			name:       "no trailing period",
			inresource: "field1.field2.field3",

			expectedInResource: "field1",
			expectedFieldsPath: []string{"field2", "field3"},
		},
		{
			name:       "trailing period with correct fieldsPath",
			inresource: "field1.field2.field3.",

			expectedInResource: "field1",
			expectedFieldsPath: []string{"field2", "field3"},
		},
		{
			name:       "trailing period with incorrect fieldsPath",
			inresource: "field1.field2.field3.",

			expectedInResource: "field1",
			expectedFieldsPath: []string{"field2", "field3", ""},
			expectedErr:        true,
		},
	}

	mapper := testrestmapper.TestOnlyStaticRESTMapper(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInResource, gotFieldsPath, err := SplitAndParseResourceRequest(tt.inresource, mapper)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.expectedInResource, gotInResource) && !tt.expectedErr {
				t.Errorf("%s: expected inresource: %s, got: %s", tt.name, tt.expectedInResource, gotInResource)
			}

			if !reflect.DeepEqual(tt.expectedFieldsPath, gotFieldsPath) && !tt.expectedErr {
				t.Errorf("%s: expected fieldsPath: %s, got: %s", tt.name, tt.expectedFieldsPath, gotFieldsPath)
			}
		})
	}
}
