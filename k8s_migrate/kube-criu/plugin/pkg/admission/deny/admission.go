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

package deny

import (
	"errors"
	"io"

	"github.com/golang/glog"

	"k8s.io/apiserver/pkg/admission"
)

// PluginName indicates name of admission plugin.
const PluginName = "AlwaysDeny"

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return NewAlwaysDeny(), nil
	})
}

// alwaysDeny is an implementation of admission.Interface which always says no to an admission request.
type alwaysDeny struct{}

var _ admission.MutationInterface = alwaysDeny{}
var _ admission.ValidationInterface = alwaysDeny{}

// Admit makes an admission decision based on the request attributes.
func (alwaysDeny) Admit(a admission.Attributes) (err error) {
	return admission.NewForbidden(a, errors.New("admission control is denying all modifications"))
}

// Validate makes an admission decision based on the request attributes.  It is NOT allowed to mutate.
func (alwaysDeny) Validate(a admission.Attributes) (err error) {
	return admission.NewForbidden(a, errors.New("admission control is denying all modifications"))
}

// Handles returns true if this admission controller can handle the given operation
// where operation can be one of CREATE, UPDATE, DELETE, or CONNECT
func (alwaysDeny) Handles(operation admission.Operation) bool {
	return true
}

// NewAlwaysDeny creates an always deny admission handler
func NewAlwaysDeny() admission.Interface {
	// DEPRECATED: AlwaysDeny denys all admission request, it is no use.
	glog.Warningf("%s admission controller is deprecated. "+
		"Please remove this controller from your configuration files and scripts.", PluginName)
	return new(alwaysDeny)
}
