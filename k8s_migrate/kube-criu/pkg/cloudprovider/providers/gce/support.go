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

package gce

import (
	"context"

	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce/cloud"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce/cloud/meta"
)

// gceProjectRouter sends requests to the appropriate project ID.
type gceProjectRouter struct {
	gce *GCECloud
}

// ProjectID returns the project ID to be used for the given operation.
func (r *gceProjectRouter) ProjectID(ctx context.Context, version meta.Version, service string) string {
	switch service {
	case "Firewalls", "Routes":
		return r.gce.NetworkProjectID()
	default:
		return r.gce.projectID
	}
}

// gceRateLimiter implements cloud.RateLimiter.
type gceRateLimiter struct {
	gce *GCECloud
}

// Accept blocks until the operation can be performed.
//
// TODO: the current cloud provider policy doesn't seem to be correct as it
// only rate limits the polling operations, but not the /submission/ of
// operations.
func (l *gceRateLimiter) Accept(ctx context.Context, key *cloud.RateLimitKey) error {
	if key.Operation == "Get" && key.Service == "Operations" {
		// Wait a minimum amount of time regardless of rate limiter.
		rl := &cloud.MinimumRateLimiter{
			// Convert flowcontrol.RateLimiter into cloud.RateLimiter
			RateLimiter: &cloud.AcceptRateLimiter{
				Acceptor: l.gce.operationPollRateLimiter,
			},
			Minimum: operationPollInterval,
		}
		return rl.Accept(ctx, key)
	}
	return nil
}

// CreateGCECloudWithCloud is a helper function to create an instance of GCECloud with the
// given Cloud interface implementation. Typical usage is to use cloud.NewMockGCE to get a
// handle to a mock Cloud instance and then use that for testing.
func CreateGCECloudWithCloud(config *CloudConfig, c cloud.Cloud) (*GCECloud, error) {
	gceCloud, err := CreateGCECloud(config)
	if err == nil {
		gceCloud.c = c
	}
	return gceCloud, err
}
