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
// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import "time"

// Timer is a helper type to time functions. Use NewTimer to create new
// instances.
type Timer struct {
	begin    time.Time
	observer Observer
}

// NewTimer creates a new Timer. The provided Observer is used to observe a
// duration in seconds. Timer is usually used to time a function call in the
// following way:
//    func TimeMe() {
//        timer := NewTimer(myHistogram)
//        defer timer.ObserveDuration()
//        // Do actual work.
//    }
func NewTimer(o Observer) *Timer {
	return &Timer{
		begin:    time.Now(),
		observer: o,
	}
}

// ObserveDuration records the duration passed since the Timer was created with
// NewTimer. It calls the Observe method of the Observer provided during
// construction with the duration in seconds as an argument. ObserveDuration is
// usually called with a defer statement.
func (t *Timer) ObserveDuration() {
	if t.observer != nil {
		t.observer.Observe(time.Since(t.begin).Seconds())
	}
}
