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
// Package logger defines interfaces that logger drivers implement to
// log messages.
//
// The other half of a logger driver is the implementation of the
// factory, which holds the contextual instance information that
// allows multiple loggers of the same type to perform different
// actions, such as logging to different locations.
package logger

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/pkg/jsonlog"
)

// ErrReadLogsNotSupported is returned when the logger does not support reading logs.
var ErrReadLogsNotSupported = errors.New("configured logging driver does not support reading")

const (
	// TimeFormat is the time format used for timestamps sent to log readers.
	TimeFormat           = jsonlog.RFC3339NanoFixed
	logWatcherBufferSize = 4096
)

var messagePool = &sync.Pool{New: func() interface{} { return &Message{Line: make([]byte, 0, 256)} }}

// NewMessage returns a new message from the message sync.Pool
func NewMessage() *Message {
	return messagePool.Get().(*Message)
}

// PutMessage puts the specified message back n the message pool.
// The message fields are reset before putting into the pool.
func PutMessage(msg *Message) {
	msg.reset()
	messagePool.Put(msg)
}

// Message is datastructure that represents piece of output produced by some
// container.  The Line member is a slice of an array whose contents can be
// changed after a log driver's Log() method returns.
// Any changes made to this struct must also be updated in the `reset` function
type Message struct {
	Line      []byte
	Source    string
	Timestamp time.Time
	Attrs     LogAttributes
	Partial   bool
}

// reset sets the message back to default values
// This is used when putting a message back into the message pool.
// Any changes to the `Message` struct should be reflected here.
func (m *Message) reset() {
	m.Line = m.Line[:0]
	m.Source = ""
	m.Attrs = nil
	m.Partial = false
}

// LogAttributes is used to hold the extra attributes available in the log message
// Primarily used for converting the map type to string and sorting.
type LogAttributes map[string]string
type byKey []string

func (s byKey) Len() int { return len(s) }
func (s byKey) Less(i, j int) bool {
	keyI := strings.Split(s[i], "=")
	keyJ := strings.Split(s[j], "=")
	return keyI[0] < keyJ[0]
}
func (s byKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (a LogAttributes) String() string {
	var ss byKey
	for k, v := range a {
		ss = append(ss, k+"="+v)
	}
	sort.Sort(ss)
	return strings.Join(ss, ",")
}

// Logger is the interface for docker logging drivers.
type Logger interface {
	Log(*Message) error
	Name() string
	Close() error
}

// ReadConfig is the configuration passed into ReadLogs.
type ReadConfig struct {
	Since  time.Time
	Tail   int
	Follow bool
}

// LogReader is the interface for reading log messages for loggers that support reading.
type LogReader interface {
	// Read logs from underlying logging backend
	ReadLogs(ReadConfig) *LogWatcher
}

// LogWatcher is used when consuming logs read from the LogReader interface.
type LogWatcher struct {
	// For sending log messages to a reader.
	Msg chan *Message
	// For sending error messages that occur while while reading logs.
	Err           chan error
	closeOnce     sync.Once
	closeNotifier chan struct{}
}

// NewLogWatcher returns a new LogWatcher.
func NewLogWatcher() *LogWatcher {
	return &LogWatcher{
		Msg:           make(chan *Message, logWatcherBufferSize),
		Err:           make(chan error, 1),
		closeNotifier: make(chan struct{}),
	}
}

// Close notifies the underlying log reader to stop.
func (w *LogWatcher) Close() {
	// only close if not already closed
	w.closeOnce.Do(func() {
		close(w.closeNotifier)
	})
}

// WatchClose returns a channel receiver that receives notification
// when the watcher has been closed. This should only be called from
// one goroutine.
func (w *LogWatcher) WatchClose() <-chan struct{} {
	return w.closeNotifier
}
