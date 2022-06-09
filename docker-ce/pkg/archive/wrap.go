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
package archive

import (
	"archive/tar"
	"bytes"
	"io"
)

// Generate generates a new archive from the content provided
// as input.
//
// `files` is a sequence of path/content pairs. A new file is
// added to the archive for each pair.
// If the last pair is incomplete, the file is created with an
// empty content. For example:
//
// Generate("foo.txt", "hello world", "emptyfile")
//
// The above call will return an archive with 2 files:
//  * ./foo.txt with content "hello world"
//  * ./empty with empty content
//
// FIXME: stream content instead of buffering
// FIXME: specify permissions and other archive metadata
func Generate(input ...string) (io.Reader, error) {
	files := parseStringPairs(input...)
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for _, file := range files {
		name, content := file[0], file[1]
		hdr := &tar.Header{
			Name: name,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

func parseStringPairs(input ...string) (output [][2]string) {
	output = make([][2]string, 0, len(input)/2+1)
	for i := 0; i < len(input); i += 2 {
		var pair [2]string
		pair[0] = input[i]
		if i+1 < len(input) {
			pair[1] = input[i+1]
		}
		output = append(output, pair)
	}
	return
}
