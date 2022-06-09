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
package discovery

import (
	"errors"
	"time"
)

var (
	// ErrNotSupported is returned when a discovery service is not supported.
	ErrNotSupported = errors.New("discovery service not supported")

	// ErrNotImplemented is returned when discovery feature is not implemented
	// by discovery backend.
	ErrNotImplemented = errors.New("not implemented in this discovery service")
)

// Watcher provides watching over a cluster for nodes joining and leaving.
type Watcher interface {
	// Watch the discovery for entry changes.
	// Returns a channel that will receive changes or an error.
	// Providing a non-nil stopCh can be used to stop watching.
	Watch(stopCh <-chan struct{}) (<-chan Entries, <-chan error)
}

// Backend is implemented by discovery backends which manage cluster entries.
type Backend interface {
	// Watcher must be provided by every backend.
	Watcher

	// Initialize the discovery with URIs, a heartbeat, a ttl and optional settings.
	Initialize(string, time.Duration, time.Duration, map[string]string) error

	// Register to the discovery.
	Register(string) error
}
