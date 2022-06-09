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
// Package tailfile provides helper functions to read the nth lines of any
// ReadSeeker.
package tailfile

import (
	"bytes"
	"errors"
	"io"
	"os"
)

const blockSize = 1024

var eol = []byte("\n")

// ErrNonPositiveLinesNumber is an error returned if the lines number was negative.
var ErrNonPositiveLinesNumber = errors.New("The number of lines to extract from the file must be positive")

//TailFile returns last n lines of reader f (could be a fil).
func TailFile(f io.ReadSeeker, n int) ([][]byte, error) {
	if n <= 0 {
		return nil, ErrNonPositiveLinesNumber
	}
	size, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		return nil, err
	}
	block := -1
	var data []byte
	var cnt int
	for {
		var b []byte
		step := int64(block * blockSize)
		left := size + step // how many bytes to beginning
		if left < 0 {
			if _, err := f.Seek(0, os.SEEK_SET); err != nil {
				return nil, err
			}
			b = make([]byte, blockSize+left)
			if _, err := f.Read(b); err != nil {
				return nil, err
			}
			data = append(b, data...)
			break
		} else {
			b = make([]byte, blockSize)
			if _, err := f.Seek(left, os.SEEK_SET); err != nil {
				return nil, err
			}
			if _, err := f.Read(b); err != nil {
				return nil, err
			}
			data = append(b, data...)
		}
		cnt += bytes.Count(b, eol)
		if cnt > n {
			break
		}
		block--
	}
	lines := bytes.Split(data, eol)
	if n < len(lines) {
		return lines[len(lines)-n-1 : len(lines)-1], nil
	}
	return lines[:len(lines)-1], nil
}
