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
package stats

import (
	"github.com/docker/docker/container"
)

// platformNewStatsCollector performs platform specific initialisation of the
// Collector structure. This is a no-op on Windows.
func platformNewStatsCollector(s *Collector) {
}

// Collect registers the container with the collector and adds it to
// the event loop for collection on the specified interval returning
// a channel for the subscriber to receive on.
// Currently not supported on Solaris
func (s *Collector) Collect(c *container.Container) chan interface{} {
	return nil
}

// StopCollection closes the channels for all subscribers and removes
// the container from metrics collection.
// Currently not supported on Solaris
func (s *Collector) StopCollection(c *container.Container) {
}

// Unsubscribe removes a specific subscriber from receiving updates for a container's stats.
// Currently not supported on Solaris
func (s *Collector) Unsubscribe(c *container.Container, ch chan interface{}) {
}
