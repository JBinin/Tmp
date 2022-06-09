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
package jsoniter

import (
	"io"
)

// IteratorPool a thread safe pool of iterators with same configuration
type IteratorPool interface {
	BorrowIterator(data []byte) *Iterator
	ReturnIterator(iter *Iterator)
}

// StreamPool a thread safe pool of streams with same configuration
type StreamPool interface {
	BorrowStream(writer io.Writer) *Stream
	ReturnStream(stream *Stream)
}

func (cfg *frozenConfig) BorrowStream(writer io.Writer) *Stream {
	stream := cfg.streamPool.Get().(*Stream)
	stream.Reset(writer)
	return stream
}

func (cfg *frozenConfig) ReturnStream(stream *Stream) {
	stream.out = nil
	stream.Error = nil
	stream.Attachment = nil
	cfg.streamPool.Put(stream)
}

func (cfg *frozenConfig) BorrowIterator(data []byte) *Iterator {
	iter := cfg.iteratorPool.Get().(*Iterator)
	iter.ResetBytes(data)
	return iter
}

func (cfg *frozenConfig) ReturnIterator(iter *Iterator) {
	iter.Error = nil
	iter.Attachment = nil
	cfg.iteratorPool.Put(iter)
}
