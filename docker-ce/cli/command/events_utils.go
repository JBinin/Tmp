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
package command

import (
	"sync"

	"github.com/Sirupsen/logrus"
	eventtypes "github.com/docker/docker/api/types/events"
)

type eventProcessor func(eventtypes.Message, error) error

// EventHandler is abstract interface for user to customize
// own handle functions of each type of events
type EventHandler interface {
	Handle(action string, h func(eventtypes.Message))
	Watch(c <-chan eventtypes.Message)
}

// InitEventHandler initializes and returns an EventHandler
func InitEventHandler() EventHandler {
	return &eventHandler{handlers: make(map[string]func(eventtypes.Message))}
}

type eventHandler struct {
	handlers map[string]func(eventtypes.Message)
	mu       sync.Mutex
}

func (w *eventHandler) Handle(action string, h func(eventtypes.Message)) {
	w.mu.Lock()
	w.handlers[action] = h
	w.mu.Unlock()
}

// Watch ranges over the passed in event chan and processes the events based on the
// handlers created for a given action.
// To stop watching, close the event chan.
func (w *eventHandler) Watch(c <-chan eventtypes.Message) {
	for e := range c {
		w.mu.Lock()
		h, exists := w.handlers[e.Action]
		w.mu.Unlock()
		if !exists {
			continue
		}
		logrus.Debugf("event handler: received event: %v", e)
		go h(e)
	}
}
