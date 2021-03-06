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
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// StartTimer begins a timer observation at the callsite. When the target
// operation is completed, the caller should call the return done func().
func StartTimer(timer Timer) (done func()) {
	start := time.Now()
	return func() {
		timer.Update(time.Since(start))
	}
}

// Timer is a metric that allows collecting the duration of an action in seconds
type Timer interface {
	// Update records an observation, duration, and converts to the target
	// units.
	Update(duration time.Duration)

	// UpdateSince will add the duration from the provided starting time to the
	// timer's summary with the precisions that was used in creation of the timer
	UpdateSince(time.Time)
}

// LabeledTimer is a timer that must have label values populated before use.
type LabeledTimer interface {
	WithValues(labels ...string) Timer
}

type labeledTimer struct {
	m *prometheus.HistogramVec
}

func (lt *labeledTimer) WithValues(labels ...string) Timer {
	return &timer{m: lt.m.WithLabelValues(labels...)}
}

func (lt *labeledTimer) Describe(c chan<- *prometheus.Desc) {
	lt.m.Describe(c)
}

func (lt *labeledTimer) Collect(c chan<- prometheus.Metric) {
	lt.m.Collect(c)
}

type timer struct {
	m prometheus.Histogram
}

func (t *timer) Update(duration time.Duration) {
	t.m.Observe(duration.Seconds())
}

func (t *timer) UpdateSince(since time.Time) {
	t.m.Observe(time.Since(since).Seconds())
}

func (t *timer) Describe(c chan<- *prometheus.Desc) {
	t.m.Describe(c)
}

func (t *timer) Collect(c chan<- prometheus.Metric) {
	t.m.Collect(c)
}
