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

type Timer struct {
	C <-chan time.Time

	timer     *time.Timer
	fakeTimer *fakeTimer
}

func (t *Timer) Reset(d time.Duration) bool {
	if t.timer != nil {
		return t.timer.Reset(d)
	}
	return t.fakeTimer.Reset(d)
}

func (t *Timer) Stop() bool {
	if t.timer != nil {
		return t.timer.Stop()
	}
	return t.fakeTimer.Stop()
}

type fakeTimer struct {
	// c is the same chan as C in the Timer that contains this fakeTimer
	c chan<- time.Time
	// clk is kept so we can maintain just one lock and to add and attempt to
	// send the times made by this timer during Resets and Stops
	clk *fake
	// active is true until the fakeTimer's send is attempted or it has been
	// stopped
	active bool
	// sends is where we store all the sends made by this timer so we can
	// deactivate the old ones when Reset or Stop is called.
	sends []*send
}

func (ft *fakeTimer) Reset(d time.Duration) bool {
	ft.clk.Lock()
	defer ft.clk.Unlock()
	target := ft.clk.t.Add(d)
	active := ft.active
	ft.active = true
	for _, s := range ft.sends {
		s.active = false
	}
	s := ft.clk.addSend(target, ft)
	ft.sends = []*send{s}
	ft.clk.sendTimes()
	return active
}

func (ft *fakeTimer) Stop() bool {
	ft.clk.Lock()
	defer ft.clk.Unlock()
	active := ft.active
	ft.active = false
	for _, s := range ft.sends {
		s.active = false
	}
	ft.sends = nil
	ft.clk.sendTimes()
	return active
}
