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
	"time"

	autoscaling "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// PodMetricsInfo contains pod metric values as a map from pod names to
// metric values (the metric values are expected to be the metric as a milli-value)
type PodMetricsInfo map[string]int64

// MetricsClient knows how to query a remote interface to retrieve container-level
// resource metrics as well as pod-level arbitrary metrics
type MetricsClient interface {
	// GetResourceMetric gets the given resource metric (and an associated oldest timestamp)
	// for all pods matching the specified selector in the given namespace
	GetResourceMetric(resource v1.ResourceName, namespace string, selector labels.Selector) (PodMetricsInfo, time.Time, error)

	// GetRawMetric gets the given metric (and an associated oldest timestamp)
	// for all pods matching the specified selector in the given namespace
	GetRawMetric(metricName string, namespace string, selector labels.Selector) (PodMetricsInfo, time.Time, error)

	// GetObjectMetric gets the given metric (and an associated timestamp) for the given
	// object in the given namespace
	GetObjectMetric(metricName string, namespace string, objectRef *autoscaling.CrossVersionObjectReference) (int64, time.Time, error)

	// GetExternalMetric gets all the values of a given external metric
	// that match the specified selector.
	GetExternalMetric(metricName string, namespace string, selector labels.Selector) ([]int64, time.Time, error)
}
