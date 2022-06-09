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
Copyright 2015 The Kubernetes Authors.

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

package fake

import (
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	core "k8s.io/client-go/testing"
)

func (c *FakeScales) Get(kind string, name string) (result *v1beta1.Scale, err error) {
	action := core.GetActionImpl{}
	action.Verb = "get"
	action.Namespace = c.ns
	action.Resource = schema.GroupVersionResource{Resource: kind}
	action.Subresource = "scale"
	action.Name = name
	obj, err := c.Fake.Invokes(action, &v1beta1.Scale{})
	result = obj.(*v1beta1.Scale)
	return
}

func (c *FakeScales) Update(kind string, scale *v1beta1.Scale) (result *v1beta1.Scale, err error) {
	action := core.UpdateActionImpl{}
	action.Verb = "update"
	action.Namespace = c.ns
	action.Resource = schema.GroupVersionResource{Resource: kind}
	action.Subresource = "scale"
	action.Object = scale
	obj, err := c.Fake.Invokes(action, scale)
	result = obj.(*v1beta1.Scale)
	return
}
