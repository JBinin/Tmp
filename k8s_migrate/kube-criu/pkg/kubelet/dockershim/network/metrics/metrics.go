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

package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// NetworkPluginOperationsKey is the key for operation count metrics.
	NetworkPluginOperationsKey = "network_plugin_operations"
	// NetworkPluginOperationsLatencyKey is the key for the operation latency metrics.
	NetworkPluginOperationsLatencyKey = "network_plugin_operations_latency_microseconds"

	// Keep the "kubelet" subsystem for backward compatibility.
	kubeletSubsystem = "kubelet"
)

var (
	// NetworkPluginOperationsLatency collects operation latency numbers by operation
	// type.
	NetworkPluginOperationsLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: kubeletSubsystem,
			Name:      NetworkPluginOperationsLatencyKey,
			Help:      "Latency in microseconds of network plugin operations. Broken down by operation type.",
		},
		[]string{"operation_type"},
	)
)

var registerMetrics sync.Once

// Register all metrics.
func Register() {
	registerMetrics.Do(func() {
		prometheus.MustRegister(NetworkPluginOperationsLatency)
	})
}

// SinceInMicroseconds gets the time since the specified start in microseconds.
func SinceInMicroseconds(start time.Time) float64 {
	return float64(time.Since(start).Nanoseconds() / time.Microsecond.Nanoseconds())
}
