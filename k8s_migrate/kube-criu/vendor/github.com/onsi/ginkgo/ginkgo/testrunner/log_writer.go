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
package testrunner

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

type logWriter struct {
	buffer *bytes.Buffer
	lock   *sync.Mutex
	log    *log.Logger
}

func newLogWriter(target io.Writer, node int) *logWriter {
	return &logWriter{
		buffer: &bytes.Buffer{},
		lock:   &sync.Mutex{},
		log:    log.New(target, fmt.Sprintf("[%d] ", node), 0),
	}
}

func (w *logWriter) Write(data []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.buffer.Write(data)
	contents := w.buffer.String()

	lines := strings.Split(contents, "\n")
	for _, line := range lines[0 : len(lines)-1] {
		w.log.Println(line)
	}

	w.buffer.Reset()
	w.buffer.Write([]byte(lines[len(lines)-1]))
	return len(data), nil
}

func (w *logWriter) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.buffer.Len() > 0 {
		w.log.Println(w.buffer.String())
	}

	return nil
}
