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

package rest

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

// RESTExportStrategy is the interface that defines how to export a Kubernetes
// object.  An exported object is stripped of non-user-settable fields and
// optionally, the identifying information related to the object's identity in
// the cluster so that it can be loaded into a different namespace or entirely
// different cluster without conflict.
type RESTExportStrategy interface {
	// Export strips fields that can not be set by the user.  If 'exact' is false
	// fields specific to the cluster are also stripped
	Export(ctx context.Context, obj runtime.Object, exact bool) error
}
