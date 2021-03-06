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
// The `fwd` package provides a buffered reader
// and writer. Each has methods that help improve
// the encoding/decoding performance of some binary
// protocols.
//
// The `fwd.Writer` and `fwd.Reader` type provide similar
// functionality to their counterparts in `bufio`, plus
// a few extra utility methods that simplify read-ahead
// and write-ahead. I wrote this package to improve serialization
// performance for http://github.com/tinylib/msgp,
// where it provided about a 2x speedup over `bufio` for certain
// workloads. However, care must be taken to understand the semantics of the
// extra methods provided by this package, as they allow
// the user to access and manipulate the buffer memory
// directly.
//
// The extra methods for `fwd.Reader` are `Peek`, `Skip`
// and `Next`. `(*fwd.Reader).Peek`, unlike `(*bufio.Reader).Peek`,
// will re-allocate the read buffer in order to accommodate arbitrarily
// large read-ahead. `(*fwd.Reader).Skip` skips the next `n` bytes
// in the stream, and uses the `io.Seeker` interface if the underlying
// stream implements it. `(*fwd.Reader).Next` returns a slice pointing
// to the next `n` bytes in the read buffer (like `Peek`), but also
// increments the read position. This allows users to process streams
// in aribtrary block sizes without having to manage appropriately-sized
// slices. Additionally, obviating the need to copy the data from the
// buffer to another location in memory can improve performance dramatically
// in CPU-bound applications.
//
// `fwd.Writer` only has one extra method, which is `(*fwd.Writer).Next`, which
// returns a slice pointing to the next `n` bytes of the writer, and increments
// the write position by the length of the returned slice. This allows users
// to write directly to the end of the buffer.
//
package fwd

import "io"

const (
	// DefaultReaderSize is the default size of the read buffer
	DefaultReaderSize = 2048

	// minimum read buffer; straight from bufio
	minReaderSize = 16
)

// NewReader returns a new *Reader that reads from 'r'
func NewReader(r io.Reader) *Reader {
	return NewReaderSize(r, DefaultReaderSize)
}

// NewReaderSize returns a new *Reader that
// reads from 'r' and has a buffer size 'n'
func NewReaderSize(r io.Reader, n int) *Reader {
	rd := &Reader{
		r:    r,
		data: make([]byte, 0, max(minReaderSize, n)),
	}
	if s, ok := r.(io.Seeker); ok {
		rd.rs = s
	}
	return rd
}

// Reader is a buffered look-ahead reader
type Reader struct {
	r io.Reader // underlying reader

	// data[n:len(data)] is buffered data; data[len(data):cap(data)] is free buffer space
	data  []byte // data
	n     int    // read offset
	state error  // last read error

	// if the reader past to NewReader was
	// also an io.Seeker, this is non-nil
	rs io.Seeker
}

// Reset resets the underlying reader
// and the read buffer.
func (r *Reader) Reset(rd io.Reader) {
	r.r = rd
	r.data = r.data[0:0]
	r.n = 0
	r.state = nil
	if s, ok := rd.(io.Seeker); ok {
		r.rs = s
	} else {
		r.rs = nil
	}
}

// more() does one read on the underlying reader
func (r *Reader) more() {
	// move data backwards so that
	// the read offset is 0; this way
	// we can supply the maximum number of
	// bytes to the reader
	if r.n != 0 {
		r.data = r.data[:copy(r.data[0:], r.data[r.n:])]
		r.n = 0
	}
	var a int
	a, r.state = r.r.Read(r.data[len(r.data):cap(r.data)])
	if a == 0 && r.state == nil {
		r.state = io.ErrNoProgress
		return
	}
	r.data = r.data[:len(r.data)+a]
}

// pop error
func (r *Reader) err() (e error) {
	e, r.state = r.state, nil
	return
}

// pop error; EOF -> io.ErrUnexpectedEOF
func (r *Reader) noEOF() (e error) {
	e, r.state = r.state, nil
	if e == io.EOF {
		e = io.ErrUnexpectedEOF
	}
	return
}

// buffered bytes
func (r *Reader) buffered() int { return len(r.data) - r.n }

// Buffered returns the number of bytes currently in the buffer
func (r *Reader) Buffered() int { return len(r.data) - r.n }

// BufferSize returns the total size of the buffer
func (r *Reader) BufferSize() int { return cap(r.data) }

// Peek returns the next 'n' buffered bytes,
// reading from the underlying reader if necessary.
// It will only return a slice shorter than 'n' bytes
// if it also returns an error. Peek does not advance
// the reader. EOF errors are *not* returned as
// io.ErrUnexpectedEOF.
func (r *Reader) Peek(n int) ([]byte, error) {
	// in the degenerate case,
	// we may need to realloc
	// (the caller asked for more
	// bytes than the size of the buffer)
	if cap(r.data) < n {
		old := r.data[r.n:]
		r.data = make([]byte, n+r.buffered())
		r.data = r.data[:copy(r.data, old)]
	}

	// keep filling until
	// we hit an error or
	// read enough bytes
	for r.buffered() < n && r.state == nil {
		r.more()
	}

	// we must have hit an error
	if r.buffered() < n {
		return r.data[r.n:], r.err()
	}

	return r.data[r.n : r.n+n], nil
}

// Skip moves the reader forward 'n' bytes.
// Returns the number of bytes skipped and any
// errors encountered. It is analagous to Seek(n, 1).
// If the underlying reader implements io.Seeker, then
// that method will be used to skip forward.
//
// If the reader encounters
// an EOF before skipping 'n' bytes, it
// returns io.ErrUnexpectedEOF. If the
// underlying reader implements io.Seeker, then
// those rules apply instead. (Many implementations
// will not return `io.EOF` until the next call
// to Read.)
func (r *Reader) Skip(n int) (int, error) {

	// fast path
	if r.buffered() >= n {
		r.n += n
		return n, nil
	}

	// use seeker implementation
	// if we can
	if r.rs != nil {
		return r.skipSeek(n)
	}

	// loop on filling
	// and then erasing
	o := n
	for r.buffered() < n && r.state == nil {
		r.more()
		// we can skip forward
		// up to r.buffered() bytes
		step := min(r.buffered(), n)
		r.n += step
		n -= step
	}
	// at this point, n should be
	// 0 if everything went smoothly
	return o - n, r.noEOF()
}

// Next returns the next 'n' bytes in the stream.
// Unlike Peek, Next advances the reader position.
// The returned bytes point to the same
// data as the buffer, so the slice is
// only valid until the next reader method call.
// An EOF is considered an unexpected error.
// If an the returned slice is less than the
// length asked for, an error will be returned,
// and the reader position will not be incremented.
func (r *Reader) Next(n int) ([]byte, error) {

	// in case the buffer is too small
	if cap(r.data) < n {
		old := r.data[r.n:]
		r.data = make([]byte, n+r.buffered())
		r.data = r.data[:copy(r.data, old)]
	}

	// fill at least 'n' bytes
	for r.buffered() < n && r.state == nil {
		r.more()
	}

	if r.buffered() < n {
		return r.data[r.n:], r.noEOF()
	}
	out := r.data[r.n : r.n+n]
	r.n += n
	return out, nil
}

// skipSeek uses the io.Seeker to seek forward.
// only call this function when n > r.buffered()
func (r *Reader) skipSeek(n int) (int, error) {
	o := r.buffered()
	// first, clear buffer
	n -= o
	r.n = 0
	r.data = r.data[:0]

	// then seek forward remaning bytes
	i, err := r.rs.Seek(int64(n), 1)
	return int(i) + o, err
}

// Read implements `io.Reader`
func (r *Reader) Read(b []byte) (int, error) {
	if len(b) <= r.buffered() {
		x := copy(b, r.data[r.n:])
		r.n += x
		return x, nil
	}
	r.more()
	if r.buffered() > 0 {
		x := copy(b, r.data[r.n:])
		r.n += x
		return x, nil
	}

	// io.Reader is supposed to return
	// 0 read bytes on error
	return 0, r.err()
}

// ReadFull attempts to read len(b) bytes into
// 'b'. It returns the number of bytes read into
// 'b', and an error if it does not return len(b).
// EOF is considered an unexpected error.
func (r *Reader) ReadFull(b []byte) (int, error) {
	var x int
	l := len(b)
	for x < l {
		if r.buffered() == 0 {
			r.more()
		}
		c := copy(b[x:], r.data[r.n:])
		x += c
		r.n += c
		if r.state != nil {
			return x, r.noEOF()
		}
	}
	return x, nil
}

// ReadByte implements `io.ByteReader`
func (r *Reader) ReadByte() (byte, error) {
	for r.buffered() < 1 && r.state == nil {
		r.more()
	}
	if r.buffered() < 1 {
		return 0, r.err()
	}
	b := r.data[r.n]
	r.n++
	return b, nil
}

// WriteTo implements `io.WriterTo`
func (r *Reader) WriteTo(w io.Writer) (int64, error) {
	var (
		i   int64
		ii  int
		err error
	)
	// first, clear buffer
	if r.buffered() > 0 {
		ii, err = w.Write(r.data[r.n:])
		i += int64(ii)
		if err != nil {
			return i, err
		}
		r.data = r.data[0:0]
		r.n = 0
	}
	for r.state == nil {
		// here we just do
		// 1:1 reads and writes
		r.more()
		if r.buffered() > 0 {
			ii, err = w.Write(r.data)
			i += int64(ii)
			if err != nil {
				return i, err
			}
			r.data = r.data[0:0]
			r.n = 0
		}
	}
	if r.state != io.EOF {
		return i, r.err()
	}
	return i, nil
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}
