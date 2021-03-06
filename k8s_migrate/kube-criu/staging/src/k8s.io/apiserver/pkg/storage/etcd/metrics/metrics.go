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

package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	cacheHitCounterOpts = prometheus.CounterOpts{
		Name: "etcd_helper_cache_hit_count",
		Help: "Counter of etcd helper cache hits.",
	}
	cacheHitCounter      = prometheus.NewCounter(cacheHitCounterOpts)
	cacheMissCounterOpts = prometheus.CounterOpts{
		Name: "etcd_helper_cache_miss_count",
		Help: "Counter of etcd helper cache miss.",
	}
	cacheMissCounter      = prometheus.NewCounter(cacheMissCounterOpts)
	cacheEntryCounterOpts = prometheus.CounterOpts{
		Name: "etcd_helper_cache_entry_count",
		Help: "Counter of etcd helper cache entries. This can be different from etcd_helper_cache_miss_count " +
			"because two concurrent threads can miss the cache and generate the same entry twice.",
	}
	cacheEntryCounter = prometheus.NewCounter(cacheEntryCounterOpts)
	cacheGetLatency   = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "etcd_request_cache_get_latencies_summary",
			Help: "Latency in microseconds of getting an object from etcd cache",
		},
	)
	cacheAddLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "etcd_request_cache_add_latencies_summary",
			Help: "Latency in microseconds of adding an object to etcd cache",
		},
	)
	etcdRequestLatenciesSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "etcd_request_latencies_summary",
			Help: "Etcd request latency summary in microseconds for each operation and object type.",
		},
		[]string{"operation", "type"},
	)
	objectCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcd_object_counts",
			Help: "Number of stored objects at the time of last check split by kind.",
		},
		[]string{"resource"},
	)
)

var registerMetrics sync.Once

// Register all metrics.
func Register() {
	// Register the metrics.
	registerMetrics.Do(func() {
		prometheus.MustRegister(cacheHitCounter)
		prometheus.MustRegister(cacheMissCounter)
		prometheus.MustRegister(cacheEntryCounter)
		prometheus.MustRegister(cacheAddLatency)
		prometheus.MustRegister(cacheGetLatency)
		prometheus.MustRegister(etcdRequestLatenciesSummary)
		prometheus.MustRegister(objectCounts)
	})
}

func UpdateObjectCount(resourcePrefix string, count int64) {
	objectCounts.WithLabelValues(resourcePrefix).Set(float64(count))
}

func RecordEtcdRequestLatency(verb, resource string, startTime time.Time) {
	etcdRequestLatenciesSummary.WithLabelValues(verb, resource).Observe(float64(time.Since(startTime) / time.Microsecond))
}

func ObserveGetCache(startTime time.Time) {
	cacheGetLatency.Observe(float64(time.Since(startTime) / time.Microsecond))
}

func ObserveAddCache(startTime time.Time) {
	cacheAddLatency.Observe(float64(time.Since(startTime) / time.Microsecond))
}

func ObserveCacheHit() {
	cacheHitCounter.Inc()
}

func ObserveCacheMiss() {
	cacheMissCounter.Inc()
}

func ObserveNewEntry() {
	cacheEntryCounter.Inc()
}

func Reset() {
	cacheHitCounter = prometheus.NewCounter(cacheHitCounterOpts)
	cacheMissCounter = prometheus.NewCounter(cacheMissCounterOpts)
	cacheEntryCounter = prometheus.NewCounter(cacheEntryCounterOpts)
	// TODO: Reset cacheAddLatency.
	// TODO: Reset cacheGetLatency.
	etcdRequestLatenciesSummary.Reset()
}
