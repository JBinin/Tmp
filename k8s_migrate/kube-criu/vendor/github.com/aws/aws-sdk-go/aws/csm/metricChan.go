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
package csm

import (
	"sync/atomic"
)

const (
	runningEnum = iota
	pausedEnum
)

var (
	// MetricsChannelSize of metrics to hold in the channel
	MetricsChannelSize = 100
)

type metricChan struct {
	ch     chan metric
	paused int64
}

func newMetricChan(size int) metricChan {
	return metricChan{
		ch: make(chan metric, size),
	}
}

func (ch *metricChan) Pause() {
	atomic.StoreInt64(&ch.paused, pausedEnum)
}

func (ch *metricChan) Continue() {
	atomic.StoreInt64(&ch.paused, runningEnum)
}

func (ch *metricChan) IsPaused() bool {
	v := atomic.LoadInt64(&ch.paused)
	return v == pausedEnum
}

// Push will push metrics to the metric channel if the channel
// is not paused
func (ch *metricChan) Push(m metric) bool {
	if ch.IsPaused() {
		return false
	}

	select {
	case ch.ch <- m:
		return true
	default:
		return false
	}
}
