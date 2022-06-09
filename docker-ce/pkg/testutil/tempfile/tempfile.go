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
package tempfile

import (
	"github.com/docker/docker/pkg/testutil/assert"
	"io/ioutil"
	"os"
)

// TempFile is a temporary file that can be used with unit tests. TempFile
// reduces the boilerplate setup required in each test case by handling
// setup errors.
type TempFile struct {
	File *os.File
}

// NewTempFile returns a new temp file with contents
func NewTempFile(t assert.TestingT, prefix string, content string) *TempFile {
	file, err := ioutil.TempFile("", prefix+"-")
	assert.NilError(t, err)

	_, err = file.Write([]byte(content))
	assert.NilError(t, err)
	file.Close()
	return &TempFile{File: file}
}

// Name returns the filename
func (f *TempFile) Name() string {
	return f.File.Name()
}

// Remove removes the file
func (f *TempFile) Remove() {
	os.Remove(f.Name())
}
