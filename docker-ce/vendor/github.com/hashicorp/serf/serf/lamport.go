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
package serf

import (
	"sync/atomic"
)

// LamportClock is a thread safe implementation of a lamport clock. It
// uses efficient atomic operations for all of its functions, falling back
// to a heavy lock only if there are enough CAS failures.
type LamportClock struct {
	counter uint64
}

// LamportTime is the value of a LamportClock.
type LamportTime uint64

// Time is used to return the current value of the lamport clock
func (l *LamportClock) Time() LamportTime {
	return LamportTime(atomic.LoadUint64(&l.counter))
}

// Increment is used to increment and return the value of the lamport clock
func (l *LamportClock) Increment() LamportTime {
	return LamportTime(atomic.AddUint64(&l.counter, 1))
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *LamportClock) Witness(v LamportTime) {
WITNESS:
	// If the other value is old, we do not need to do anything
	cur := atomic.LoadUint64(&l.counter)
	other := uint64(v)
	if other < cur {
		return
	}

	// Ensure that our local clock is at least one ahead.
	if !atomic.CompareAndSwapUint64(&l.counter, cur, other+1) {
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future witness will pass us by and our witness
		// will end.
		goto WITNESS
	}
}
