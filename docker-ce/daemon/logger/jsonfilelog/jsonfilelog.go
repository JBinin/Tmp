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
// Package jsonfilelog provides the default Logger implementation for
// Docker logging. This logger logs to files on the host server in the
// JSON format.
package jsonfilelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/docker/docker/pkg/jsonlog"
	"github.com/docker/go-units"
)

// Name is the name of the file that the jsonlogger logs to.
const Name = "json-file"

// JSONFileLogger is Logger implementation for default Docker logging.
type JSONFileLogger struct {
	buf     *bytes.Buffer
	writer  *loggerutils.RotateFileWriter
	mu      sync.Mutex
	readers map[*logger.LogWatcher]struct{} // stores the active log followers
	extra   []byte                          // json-encoded extra attributes
}

func init() {
	if err := logger.RegisterLogDriver(Name, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(Name, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// New creates new JSONFileLogger which writes to filename passed in
// on given context.
func New(info logger.Info) (logger.Logger, error) {
	var capval int64 = -1
	if capacity, ok := info.Config["max-size"]; ok {
		var err error
		capval, err = units.FromHumanSize(capacity)
		if err != nil {
			return nil, err
		}
	}
	var maxFiles = 1
	if maxFileString, ok := info.Config["max-file"]; ok {
		var err error
		maxFiles, err = strconv.Atoi(maxFileString)
		if err != nil {
			return nil, err
		}
		if maxFiles < 1 {
			return nil, fmt.Errorf("max-file cannot be less than 1")
		}
	}

	writer, err := loggerutils.NewRotateFileWriter(info.LogPath, capval, maxFiles)
	if err != nil {
		return nil, err
	}

	var extra []byte
	if attrs := info.ExtraAttributes(nil); len(attrs) > 0 {
		var err error
		extra, err = json.Marshal(attrs)
		if err != nil {
			return nil, err
		}
	}

	return &JSONFileLogger{
		buf:     bytes.NewBuffer(nil),
		writer:  writer,
		readers: make(map[*logger.LogWatcher]struct{}),
		extra:   extra,
	}, nil
}

// Log converts logger.Message to jsonlog.JSONLog and serializes it to file.
func (l *JSONFileLogger) Log(msg *logger.Message) error {
	timestamp, err := jsonlog.FastTimeMarshalJSON(msg.Timestamp)
	if err != nil {
		return err
	}
	l.mu.Lock()
	logline := msg.Line
	if !msg.Partial {
		logline = append(msg.Line, '\n')
	}
	err = (&jsonlog.JSONLogs{
		Log:      logline,
		Stream:   msg.Source,
		Created:  timestamp,
		RawAttrs: l.extra,
	}).MarshalJSONBuf(l.buf)
	logger.PutMessage(msg)
	if err != nil {
		l.mu.Unlock()
		return err
	}

	l.buf.WriteByte('\n')
	_, err = l.writer.Write(l.buf.Bytes())
	l.buf.Reset()
	l.mu.Unlock()

	return err
}

// ValidateLogOpt looks for json specific log options max-file & max-size.
func ValidateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case "max-file":
		case "max-size":
		case "labels":
		case "env":
		default:
			return fmt.Errorf("unknown log opt '%s' for json-file log driver", key)
		}
	}
	return nil
}

// LogPath returns the location the given json logger logs to.
func (l *JSONFileLogger) LogPath() string {
	return l.writer.LogPath()
}

// Close closes underlying file and signals all readers to stop.
func (l *JSONFileLogger) Close() error {
	l.mu.Lock()
	err := l.writer.Close()
	for r := range l.readers {
		r.Close()
		delete(l.readers, r)
	}
	l.mu.Unlock()
	return err
}

// Name returns name of this logger.
func (l *JSONFileLogger) Name() string {
	return Name
}
