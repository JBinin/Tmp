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
package request

import (
	"io"
	"sync"
)

// offsetReader is a thread-safe io.ReadCloser to prevent racing
// with retrying requests
type offsetReader struct {
	buf    io.ReadSeeker
	lock   sync.Mutex
	closed bool
}

func newOffsetReader(buf io.ReadSeeker, offset int64) *offsetReader {
	reader := &offsetReader{}
	buf.Seek(offset, 0)

	reader.buf = buf
	return reader
}

// Close will close the instance of the offset reader's access to
// the underlying io.ReadSeeker.
func (o *offsetReader) Close() error {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.closed = true
	return nil
}

// Read is a thread-safe read of the underlying io.ReadSeeker
func (o *offsetReader) Read(p []byte) (int, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.closed {
		return 0, io.EOF
	}

	return o.buf.Read(p)
}

// Seek is a thread-safe seeking operation.
func (o *offsetReader) Seek(offset int64, whence int) (int64, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	return o.buf.Seek(offset, whence)
}

// CloseAndCopy will return a new offsetReader with a copy of the old buffer
// and close the old buffer.
func (o *offsetReader) CloseAndCopy(offset int64) *offsetReader {
	o.Close()
	return newOffsetReader(o.buf, offset)
}
