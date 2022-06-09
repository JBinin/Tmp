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
	"bufio"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/pubsub"
)

type supervisor interface {
	// GetContainerStats collects all the stats related to a container
	GetContainerStats(container *container.Container) (*types.StatsJSON, error)
}

// NewCollector creates a stats collector that will poll the supervisor with the specified interval
func NewCollector(supervisor supervisor, interval time.Duration) *Collector {
	s := &Collector{
		interval:   interval,
		supervisor: supervisor,
		publishers: make(map[*container.Container]*pubsub.Publisher),
		bufReader:  bufio.NewReaderSize(nil, 128),
	}

	platformNewStatsCollector(s)

	return s
}

// Collector manages and provides container resource stats
type Collector struct {
	m          sync.Mutex
	supervisor supervisor
	interval   time.Duration
	publishers map[*container.Container]*pubsub.Publisher
	bufReader  *bufio.Reader

	// The following fields are not set on Windows currently.
	clockTicksPerSecond uint64
}
