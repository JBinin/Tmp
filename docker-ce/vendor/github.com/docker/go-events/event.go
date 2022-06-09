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

// Event marks items that can be sent as events.
type Event interface{}

// Sink accepts and sends events.
type Sink interface {
	// Write an event to the Sink. If no error is returned, the caller will
	// assume that all events have been committed to the sink. If an error is
	// received, the caller may retry sending the event.
	Write(event Event) error

	// Close the sink, possibly waiting for pending events to flush.
	Close() error
}
