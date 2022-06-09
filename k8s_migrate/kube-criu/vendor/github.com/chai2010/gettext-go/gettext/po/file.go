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
// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
)

// File represents an PO File.
//
// See http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html
type File struct {
	MimeHeader Header
	Messages   []Message
}

// Load loads a named po file.
func Load(name string) (*File, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return LoadData(data)
}

// LoadData loads po file format data.
func LoadData(data []byte) (*File, error) {
	r := newLineReader(string(data))
	var file File
	for {
		var msg Message
		if err := msg.readPoEntry(r); err != nil {
			if err == io.EOF {
				return &file, nil
			}
			return nil, err
		}
		if msg.MsgId == "" {
			file.MimeHeader.parseHeader(&msg)
			continue
		}
		file.Messages = append(file.Messages, msg)
	}
}

// Save saves a po file.
func (f *File) Save(name string) error {
	return ioutil.WriteFile(name, []byte(f.String()), 0666)
}

// Save returns a po file format data.
func (f *File) Data() []byte {
	// sort the massge as ReferenceFile/ReferenceLine field
	var messages []Message
	messages = append(messages, f.Messages...)
	sort.Sort(byMessages(messages))

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", f.MimeHeader.String())
	for i := 0; i < len(messages); i++ {
		fmt.Fprintf(&buf, "%s\n", messages[i].String())
	}
	return buf.Bytes()
}

// String returns the po format file string.
func (f *File) String() string {
	return string(f.Data())
}
