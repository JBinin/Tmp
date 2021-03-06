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
package clock

import "time"

type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	Since(t time.Time) time.Duration

	NewTimer(d time.Duration) Timer
	NewTicker(d time.Duration) Ticker
}

type realClock struct{}

func NewClock() Clock {
	return &realClock{}
}

func (clock *realClock) Now() time.Time {
	return time.Now()
}

func (clock *realClock) Since(t time.Time) time.Duration {
	return time.Now().Sub(t)
}

func (clock *realClock) Sleep(d time.Duration) {
	<-clock.NewTimer(d).C()
}

func (clock *realClock) NewTimer(d time.Duration) Timer {
	return &realTimer{
		t: time.NewTimer(d),
	}
}

func (clock *realClock) NewTicker(d time.Duration) Ticker {
	return &realTicker{
		t: time.NewTicker(d),
	}
}
