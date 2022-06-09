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
// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package container defines types for sub-container events and also
// defines an interface for container operation handlers.
package watcher

// SubcontainerEventType indicates an addition or deletion event.
type ContainerEventType int

const (
	ContainerAdd ContainerEventType = iota
	ContainerDelete
)

type ContainerWatchSource int

const (
	Raw ContainerWatchSource = iota
	Rkt
)

// ContainerEvent represents a
type ContainerEvent struct {
	// The type of event that occurred.
	EventType ContainerEventType

	// The full container name of the container where the event occurred.
	Name string

	// The watcher that detected this change event
	WatchSource ContainerWatchSource
}

type ContainerWatcher interface {
	// Registers a channel to listen for events affecting subcontainers (recursively).
	Start(events chan ContainerEvent) error

	// Stops watching for subcontainer changes.
	Stop() error
}
