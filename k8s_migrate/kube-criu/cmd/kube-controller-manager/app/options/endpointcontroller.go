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
Copyright 2018 The Kubernetes Authors.

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

package options

import (
	"github.com/spf13/pflag"

	"k8s.io/kubernetes/pkg/apis/componentconfig"
)

// EndPointControllerOptions holds the EndPointController options.
type EndPointControllerOptions struct {
	ConcurrentEndpointSyncs int32
}

// AddFlags adds flags related to EndPointController for controller manager to the specified FlagSet.
func (o *EndPointControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentEndpointSyncs, "concurrent-endpoint-syncs", o.ConcurrentEndpointSyncs, "The number of endpoint syncing operations that will be done concurrently. Larger number = faster endpoint updating, but more CPU (and network) load")
}

// ApplyTo fills up EndPointController config with options.
func (o *EndPointControllerOptions) ApplyTo(cfg *componentconfig.EndPointControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentEndpointSyncs = o.ConcurrentEndpointSyncs

	return nil
}

// Validate checks validation of EndPointControllerOptions.
func (o *EndPointControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
