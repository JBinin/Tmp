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
package heartbeat

import (
	"sync/atomic"
	"time"
)

// Heartbeat is simple way to track heartbeats.
type Heartbeat struct {
	timeout int64
	timer   *time.Timer
}

// New creates new Heartbeat with specified duration. timeoutFunc will be called
// if timeout for heartbeat is expired. Note that in case of timeout you need to
// call Beat() to reactivate Heartbeat.
func New(timeout time.Duration, timeoutFunc func()) *Heartbeat {
	hb := &Heartbeat{
		timeout: int64(timeout),
		timer:   time.AfterFunc(timeout, timeoutFunc),
	}
	return hb
}

// Beat resets internal timer to zero. It also can be used to reactivate
// Heartbeat after timeout.
func (hb *Heartbeat) Beat() {
	hb.timer.Reset(time.Duration(atomic.LoadInt64(&hb.timeout)))
}

// Update updates internal timeout to d. It does not do Beat.
func (hb *Heartbeat) Update(d time.Duration) {
	atomic.StoreInt64(&hb.timeout, int64(d))
}

// Stop stops Heartbeat timer.
func (hb *Heartbeat) Stop() {
	hb.timer.Stop()
}
