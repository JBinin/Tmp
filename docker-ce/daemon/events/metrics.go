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
package events

import "github.com/docker/go-metrics"

var (
	eventsCounter    metrics.Counter
	eventSubscribers metrics.Gauge
)

func init() {
	ns := metrics.NewNamespace("engine", "daemon", nil)
	eventsCounter = ns.NewCounter("events", "The number of events logged")
	eventSubscribers = ns.NewGauge("events_subscribers", "The number of current subscribers to events", metrics.Total)
	metrics.Register(ns)
}
