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
package restful

// Copyright 2015 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"compress/gzip"
	"compress/zlib"
)

// CompressorProvider describes a component that can provider compressors for the std methods.
type CompressorProvider interface {
	// Returns a *gzip.Writer which needs to be released later.
	// Before using it, call Reset().
	AcquireGzipWriter() *gzip.Writer

	// Releases an aqcuired *gzip.Writer.
	ReleaseGzipWriter(w *gzip.Writer)

	// Returns a *gzip.Reader which needs to be released later.
	AcquireGzipReader() *gzip.Reader

	// Releases an aqcuired *gzip.Reader.
	ReleaseGzipReader(w *gzip.Reader)

	// Returns a *zlib.Writer which needs to be released later.
	// Before using it, call Reset().
	AcquireZlibWriter() *zlib.Writer

	// Releases an aqcuired *zlib.Writer.
	ReleaseZlibWriter(w *zlib.Writer)
}

// DefaultCompressorProvider is the actual provider of compressors (zlib or gzip).
var currentCompressorProvider CompressorProvider

func init() {
	currentCompressorProvider = NewSyncPoolCompessors()
}

// CurrentCompressorProvider returns the current CompressorProvider.
// It is initialized using a SyncPoolCompessors.
func CurrentCompressorProvider() CompressorProvider {
	return currentCompressorProvider
}

// CompressorProvider sets the actual provider of compressors (zlib or gzip).
func SetCompressorProvider(p CompressorProvider) {
	if p == nil {
		panic("cannot set compressor provider to nil")
	}
	currentCompressorProvider = p
}
