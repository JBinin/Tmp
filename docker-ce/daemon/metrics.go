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
package daemon

import "github.com/docker/go-metrics"

var (
	containerActions          metrics.LabeledTimer
	imageActions              metrics.LabeledTimer
	networkActions            metrics.LabeledTimer
	engineVersion             metrics.LabeledGauge
	engineCpus                metrics.Gauge
	engineMemory              metrics.Gauge
	healthChecksCounter       metrics.Counter
	healthChecksFailedCounter metrics.Counter
)

func init() {
	ns := metrics.NewNamespace("engine", "daemon", nil)
	containerActions = ns.NewLabeledTimer("container_actions", "The number of seconds it takes to process each container action", "action")
	for _, a := range []string{
		"start",
		"changes",
		"commit",
		"create",
		"delete",
	} {
		containerActions.WithValues(a).Update(0)
	}
	networkActions = ns.NewLabeledTimer("network_actions", "The number of seconds it takes to process each network action", "action")
	engineVersion = ns.NewLabeledGauge("engine", "The version and commit information for the engine process", metrics.Unit("info"),
		"version",
		"commit",
		"architecture",
		"graph_driver", "kernel",
		"os",
	)
	engineCpus = ns.NewGauge("engine_cpus", "The number of cpus that the host system of the engine has", metrics.Unit("cpus"))
	engineMemory = ns.NewGauge("engine_memory", "The number of bytes of memory that the host system of the engine has", metrics.Bytes)
	healthChecksCounter = ns.NewCounter("health_checks", "The total number of health checks")
	healthChecksFailedCounter = ns.NewCounter("health_checks_failed", "The total number of failed health checks")
	imageActions = ns.NewLabeledTimer("image_actions", "The number of seconds it takes to process each image action", "action")
	metrics.Register(ns)
}
