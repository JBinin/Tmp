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
/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flushwriter

import (
	"io"
	"net/http"
)

// Wrap wraps an io.Writer into a writer that flushes after every write if
// the writer implements the Flusher interface.
func Wrap(w io.Writer) io.Writer {
	fw := &flushWriter{
		writer: w,
	}
	if flusher, ok := w.(http.Flusher); ok {
		fw.flusher = flusher
	}
	return fw
}

// flushWriter provides wrapper for responseWriter with HTTP streaming capabilities
type flushWriter struct {
	flusher http.Flusher
	writer  io.Writer
}

// Write is a FlushWriter implementation of the io.Writer that sends any buffered
// data to the client.
func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.writer.Write(p)
	if err != nil {
		return
	}
	if fw.flusher != nil {
		fw.flusher.Flush()
	}
	return
}
