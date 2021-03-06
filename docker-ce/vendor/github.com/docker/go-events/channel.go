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
package events

import (
	"fmt"
	"sync"
)

// Channel provides a sink that can be listened on. The writer and channel
// listener must operate in separate goroutines.
//
// Consumers should listen on Channel.C until Closed is closed.
type Channel struct {
	C chan Event

	closed chan struct{}
	once   sync.Once
}

// NewChannel returns a channel. If buffer is zero, the channel is
// unbuffered.
func NewChannel(buffer int) *Channel {
	return &Channel{
		C:      make(chan Event, buffer),
		closed: make(chan struct{}),
	}
}

// Done returns a channel that will always proceed once the sink is closed.
func (ch *Channel) Done() chan struct{} {
	return ch.closed
}

// Write the event to the channel. Must be called in a separate goroutine from
// the listener.
func (ch *Channel) Write(event Event) error {
	select {
	case ch.C <- event:
		return nil
	case <-ch.closed:
		return ErrSinkClosed
	}
}

// Close the channel sink.
func (ch *Channel) Close() error {
	ch.once.Do(func() {
		close(ch.closed)
	})

	return nil
}

func (ch Channel) String() string {
	// Serialize a copy of the Channel that doesn't contain the sync.Once,
	// to avoid a data race.
	ch2 := map[string]interface{}{
		"C":      ch.C,
		"closed": ch.closed,
	}
	return fmt.Sprint(ch2)
}
