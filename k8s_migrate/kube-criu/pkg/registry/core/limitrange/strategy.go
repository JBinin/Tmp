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

package limitrange

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/validation"
)

type limitrangeStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// LimitRange objects via the REST API.
var Strategy = limitrangeStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

func (limitrangeStrategy) NamespaceScoped() bool {
	return true
}

func (limitrangeStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	limitRange := obj.(*api.LimitRange)
	if len(limitRange.Name) == 0 {
		limitRange.Name = string(uuid.NewUUID())
	}
}

func (limitrangeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (limitrangeStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	limitRange := obj.(*api.LimitRange)
	return validation.ValidateLimitRange(limitRange)
}

// Canonicalize normalizes the object after validation.
func (limitrangeStrategy) Canonicalize(obj runtime.Object) {
}

func (limitrangeStrategy) AllowCreateOnUpdate() bool {
	return true
}

func (limitrangeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	limitRange := obj.(*api.LimitRange)
	return validation.ValidateLimitRange(limitRange)
}

func (limitrangeStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (limitrangeStrategy) Export(context.Context, runtime.Object, bool) error {
	// Copied from OpenShift exporter
	// TODO: this needs to be fixed
	//  limitrange.Strategy.PrepareForCreate(ctx, obj)
	return nil
}
