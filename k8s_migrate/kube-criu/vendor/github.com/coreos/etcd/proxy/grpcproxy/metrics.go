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
// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpcproxy

import "github.com/prometheus/client_golang/prometheus"

var (
	watchersCoalescing = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "grpc_proxy",
		Name:      "watchers_coalescing_total",
		Help:      "Total number of current watchers coalescing",
	})
	eventsCoalescing = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "grpc_proxy",
		Name:      "events_coalescing_total",
		Help:      "Total number of events coalescing",
	})
	cacheKeys = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "grpc_proxy",
		Name:      "cache_keys_total",
		Help:      "Total number of keys/ranges cached",
	})
	cacheHits = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "grpc_proxy",
		Name:      "cache_hits_total",
		Help:      "Total number of cache hits",
	})
	cachedMisses = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "grpc_proxy",
		Name:      "cache_misses_total",
		Help:      "Total number of cache misses",
	})
)

func init() {
	prometheus.MustRegister(watchersCoalescing)
	prometheus.MustRegister(eventsCoalescing)
	prometheus.MustRegister(cacheKeys)
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cachedMisses)
}
