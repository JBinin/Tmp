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
package writer

import (
	"bytes"
	"io"
	"sync"
)

type WriterInterface interface {
	io.Writer

	Truncate()
	DumpOut()
	DumpOutWithHeader(header string)
	Bytes() []byte
}

type Writer struct {
	buffer    *bytes.Buffer
	outWriter io.Writer
	lock      *sync.Mutex
	stream    bool
}

func New(outWriter io.Writer) *Writer {
	return &Writer{
		buffer:    &bytes.Buffer{},
		lock:      &sync.Mutex{},
		outWriter: outWriter,
		stream:    true,
	}
}

func (w *Writer) SetStream(stream bool) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.stream = stream
}

func (w *Writer) Write(b []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	n, err = w.buffer.Write(b)
	if w.stream {
		return w.outWriter.Write(b)
	}
	return n, err
}

func (w *Writer) Truncate() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.buffer.Reset()
}

func (w *Writer) DumpOut() {
	w.lock.Lock()
	defer w.lock.Unlock()
	if !w.stream {
		w.buffer.WriteTo(w.outWriter)
	}
}

func (w *Writer) Bytes() []byte {
	w.lock.Lock()
	defer w.lock.Unlock()
	b := w.buffer.Bytes()
	copied := make([]byte, len(b))
	copy(copied, b)
	return copied
}

func (w *Writer) DumpOutWithHeader(header string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if !w.stream && w.buffer.Len() > 0 {
		w.outWriter.Write([]byte(header))
		w.buffer.WriteTo(w.outWriter)
	}
}
