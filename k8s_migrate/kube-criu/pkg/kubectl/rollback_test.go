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

package kubectl

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
)

var rollbacktests = map[schema.GroupKind]reflect.Type{
	{Group: "apps", Kind: "DaemonSet"}:   reflect.TypeOf(&DaemonSetRollbacker{}),
	{Group: "apps", Kind: "StatefulSet"}: reflect.TypeOf(&StatefulSetRollbacker{}),
	{Group: "apps", Kind: "Deployment"}:  reflect.TypeOf(&DeploymentRollbacker{}),
}

func TestRollbackerFor(t *testing.T) {
	fakeClientset := &fake.Clientset{}

	for kind, expectedType := range rollbacktests {
		result, err := RollbackerFor(kind, fakeClientset)
		if err != nil {
			t.Fatalf("error getting Rollbacker for a %v: %v", kind.String(), err)
		}

		if reflect.TypeOf(result) != expectedType {
			t.Fatalf("unexpected output type (%v was expected but got %v)", expectedType, reflect.TypeOf(result))
		}
	}
}
