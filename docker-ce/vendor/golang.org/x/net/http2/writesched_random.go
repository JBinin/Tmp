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
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import "math"

// NewRandomWriteScheduler constructs a WriteScheduler that ignores HTTP/2
// priorities. Control frames like SETTINGS and PING are written before DATA
// frames, but if no control frames are queued and multiple streams have queued
// HEADERS or DATA frames, Pop selects a ready stream arbitrarily.
func NewRandomWriteScheduler() WriteScheduler {
	return &randomWriteScheduler{sq: make(map[uint32]*writeQueue)}
}

type randomWriteScheduler struct {
	// zero are frames not associated with a specific stream.
	zero writeQueue

	// sq contains the stream-specific queues, keyed by stream ID.
	// When a stream is idle or closed, it's deleted from the map.
	sq map[uint32]*writeQueue

	// pool of empty queues for reuse.
	queuePool writeQueuePool
}

func (ws *randomWriteScheduler) OpenStream(streamID uint32, options OpenStreamOptions) {
	// no-op: idle streams are not tracked
}

func (ws *randomWriteScheduler) CloseStream(streamID uint32) {
	q, ok := ws.sq[streamID]
	if !ok {
		return
	}
	delete(ws.sq, streamID)
	ws.queuePool.put(q)
}

func (ws *randomWriteScheduler) AdjustStream(streamID uint32, priority PriorityParam) {
	// no-op: priorities are ignored
}

func (ws *randomWriteScheduler) Push(wr FrameWriteRequest) {
	id := wr.StreamID()
	if id == 0 {
		ws.zero.push(wr)
		return
	}
	q, ok := ws.sq[id]
	if !ok {
		q = ws.queuePool.get()
		ws.sq[id] = q
	}
	q.push(wr)
}

func (ws *randomWriteScheduler) Pop() (FrameWriteRequest, bool) {
	// Control frames first.
	if !ws.zero.empty() {
		return ws.zero.shift(), true
	}
	// Iterate over all non-idle streams until finding one that can be consumed.
	for _, q := range ws.sq {
		if wr, ok := q.consume(math.MaxInt32); ok {
			return wr, true
		}
	}
	return FrameWriteRequest{}, false
}
