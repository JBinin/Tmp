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
	"github.com/golang/glog"

	computealpha "google.golang.org/api/compute/v0.alpha"
	computebeta "google.golang.org/api/compute/v0.beta"
	compute "google.golang.org/api/compute/v1"

	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce/cloud"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce/cloud/filter"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce/cloud/meta"
	"k8s.io/kubernetes/pkg/master/ports"
	utilversion "k8s.io/kubernetes/pkg/util/version"
)

const (
	nodesHealthCheckPath   = "/healthz"
	lbNodesHealthCheckPort = ports.ProxyHealthzPort
)

var (
	minNodesHealthCheckVersion *utilversion.Version
)

func init() {
	if v, err := utilversion.ParseGeneric("1.7.2"); err != nil {
		glog.Fatalf("Failed to parse version for minNodesHealthCheckVersion: %v", err)
	} else {
		minNodesHealthCheckVersion = v
	}
}

func newHealthcheckMetricContext(request string) *metricContext {
	return newHealthcheckMetricContextWithVersion(request, computeV1Version)
}

func newHealthcheckMetricContextWithVersion(request, version string) *metricContext {
	return newGenericMetricContext("healthcheck", request, unusedMetricLabel, unusedMetricLabel, version)
}

// GetHttpHealthCheck returns the given HttpHealthCheck by name.
func (gce *GCECloud) GetHttpHealthCheck(name string) (*compute.HttpHealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("get_legacy")
	v, err := gce.c.HttpHealthChecks().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// UpdateHttpHealthCheck applies the given HttpHealthCheck as an update.
func (gce *GCECloud) UpdateHttpHealthCheck(hc *compute.HttpHealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("update_legacy")
	return mc.Observe(gce.c.HttpHealthChecks().Update(ctx, meta.GlobalKey(hc.Name), hc))
}

// DeleteHttpHealthCheck deletes the given HttpHealthCheck by name.
func (gce *GCECloud) DeleteHttpHealthCheck(name string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("delete_legacy")
	return mc.Observe(gce.c.HttpHealthChecks().Delete(ctx, meta.GlobalKey(name)))
}

// CreateHttpHealthCheck creates the given HttpHealthCheck.
func (gce *GCECloud) CreateHttpHealthCheck(hc *compute.HttpHealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("create_legacy")
	return mc.Observe(gce.c.HttpHealthChecks().Insert(ctx, meta.GlobalKey(hc.Name), hc))
}

// ListHttpHealthChecks lists all HttpHealthChecks in the project.
func (gce *GCECloud) ListHttpHealthChecks() ([]*compute.HttpHealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("list_legacy")
	v, err := gce.c.HttpHealthChecks().List(ctx, filter.None)
	return v, mc.Observe(err)
}

// Legacy HTTPS Health Checks

// GetHttpsHealthCheck returns the given HttpsHealthCheck by name.
func (gce *GCECloud) GetHttpsHealthCheck(name string) (*compute.HttpsHealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("get_legacy")
	v, err := gce.c.HttpsHealthChecks().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// UpdateHttpsHealthCheck applies the given HttpsHealthCheck as an update.
func (gce *GCECloud) UpdateHttpsHealthCheck(hc *compute.HttpsHealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("update_legacy")
	return mc.Observe(gce.c.HttpsHealthChecks().Update(ctx, meta.GlobalKey(hc.Name), hc))
}

// DeleteHttpsHealthCheck deletes the given HttpsHealthCheck by name.
func (gce *GCECloud) DeleteHttpsHealthCheck(name string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("delete_legacy")
	return mc.Observe(gce.c.HttpsHealthChecks().Delete(ctx, meta.GlobalKey(name)))
}

// CreateHttpsHealthCheck creates the given HttpsHealthCheck.
func (gce *GCECloud) CreateHttpsHealthCheck(hc *compute.HttpsHealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("create_legacy")
	return mc.Observe(gce.c.HttpsHealthChecks().Insert(ctx, meta.GlobalKey(hc.Name), hc))
}

// ListHttpsHealthChecks lists all HttpsHealthChecks in the project.
func (gce *GCECloud) ListHttpsHealthChecks() ([]*compute.HttpsHealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("list_legacy")
	v, err := gce.c.HttpsHealthChecks().List(ctx, filter.None)
	return v, mc.Observe(err)
}

// Generic HealthCheck

// GetHealthCheck returns the given HealthCheck by name.
func (gce *GCECloud) GetHealthCheck(name string) (*compute.HealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("get")
	v, err := gce.c.HealthChecks().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// GetAlphaHealthCheck returns the given alpha HealthCheck by name.
func (gce *GCECloud) GetAlphaHealthCheck(name string) (*computealpha.HealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContextWithVersion("get", computeAlphaVersion)
	v, err := gce.c.AlphaHealthChecks().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// GetBetaHealthCheck returns the given beta HealthCheck by name.
func (gce *GCECloud) GetBetaHealthCheck(name string) (*computebeta.HealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContextWithVersion("get", computeBetaVersion)
	v, err := gce.c.BetaHealthChecks().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// UpdateHealthCheck applies the given HealthCheck as an update.
func (gce *GCECloud) UpdateHealthCheck(hc *compute.HealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("update")
	return mc.Observe(gce.c.HealthChecks().Update(ctx, meta.GlobalKey(hc.Name), hc))
}

// UpdateAlphaHealthCheck applies the given alpha HealthCheck as an update.
func (gce *GCECloud) UpdateAlphaHealthCheck(hc *computealpha.HealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContextWithVersion("update", computeAlphaVersion)
	return mc.Observe(gce.c.AlphaHealthChecks().Update(ctx, meta.GlobalKey(hc.Name), hc))
}

// UpdateBetaHealthCheck applies the given beta HealthCheck as an update.
func (gce *GCECloud) UpdateBetaHealthCheck(hc *computebeta.HealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContextWithVersion("update", computeBetaVersion)
	return mc.Observe(gce.c.BetaHealthChecks().Update(ctx, meta.GlobalKey(hc.Name), hc))
}

// DeleteHealthCheck deletes the given HealthCheck by name.
func (gce *GCECloud) DeleteHealthCheck(name string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("delete")
	return mc.Observe(gce.c.HealthChecks().Delete(ctx, meta.GlobalKey(name)))
}

// CreateHealthCheck creates the given HealthCheck.
func (gce *GCECloud) CreateHealthCheck(hc *compute.HealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("create")
	return mc.Observe(gce.c.HealthChecks().Insert(ctx, meta.GlobalKey(hc.Name), hc))
}

// CreateAlphaHealthCheck creates the given alpha HealthCheck.
func (gce *GCECloud) CreateAlphaHealthCheck(hc *computealpha.HealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContextWithVersion("create", computeAlphaVersion)
	return mc.Observe(gce.c.AlphaHealthChecks().Insert(ctx, meta.GlobalKey(hc.Name), hc))
}

// CreateBetaHealthCheck creates the given beta HealthCheck.
func (gce *GCECloud) CreateBetaHealthCheck(hc *computebeta.HealthCheck) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContextWithVersion("create", computeBetaVersion)
	return mc.Observe(gce.c.BetaHealthChecks().Insert(ctx, meta.GlobalKey(hc.Name), hc))
}

// ListHealthChecks lists all HealthCheck in the project.
func (gce *GCECloud) ListHealthChecks() ([]*compute.HealthCheck, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newHealthcheckMetricContext("list")
	v, err := gce.c.HealthChecks().List(ctx, filter.None)
	return v, mc.Observe(err)
}

// GetNodesHealthCheckPort returns the health check port used by the GCE load
// balancers (l4) for performing health checks on nodes.
func GetNodesHealthCheckPort() int32 {
	return lbNodesHealthCheckPort
}

// GetNodesHealthCheckPath returns the health check path used by the GCE load
// balancers (l4) for performing health checks on nodes.
func GetNodesHealthCheckPath() string {
	return nodesHealthCheckPath
}

// isAtLeastMinNodesHealthCheckVersion checks if a version is higher than
// `minNodesHealthCheckVersion`.
func isAtLeastMinNodesHealthCheckVersion(vstring string) bool {
	version, err := utilversion.ParseGeneric(vstring)
	if err != nil {
		glog.Errorf("vstring (%s) is not a valid version string: %v", vstring, err)
		return false
	}
	return version.AtLeast(minNodesHealthCheckVersion)
}

// supportsNodesHealthCheck returns false if anyone of the nodes has version
// lower than `minNodesHealthCheckVersion`.
func supportsNodesHealthCheck(nodes []*v1.Node) bool {
	for _, node := range nodes {
		if !isAtLeastMinNodesHealthCheckVersion(node.Status.NodeInfo.KubeProxyVersion) {
			return false
		}
	}
	return true
}
