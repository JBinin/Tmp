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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
)

// DeploymentControllerOptions holds the DeploymentController options.
type DeploymentControllerOptions struct {
	ConcurrentDeploymentSyncs      int32
	DeploymentControllerSyncPeriod metav1.Duration
}

// AddFlags adds flags related to DeploymentController for controller manager to the specified FlagSet.
func (o *DeploymentControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentDeploymentSyncs, "concurrent-deployment-syncs", o.ConcurrentDeploymentSyncs, "The number of deployment objects that are allowed to sync concurrently. Larger number = more responsive deployments, but more CPU (and network) load")
	fs.DurationVar(&o.DeploymentControllerSyncPeriod.Duration, "deployment-controller-sync-period", o.DeploymentControllerSyncPeriod.Duration, "Period for syncing the deployments.")
}

// ApplyTo fills up DeploymentController config with options.
func (o *DeploymentControllerOptions) ApplyTo(cfg *componentconfig.DeploymentControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentDeploymentSyncs = o.ConcurrentDeploymentSyncs
	cfg.DeploymentControllerSyncPeriod = o.DeploymentControllerSyncPeriod

	return nil
}

// Validate checks validation of DeploymentControllerOptions.
func (o *DeploymentControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
