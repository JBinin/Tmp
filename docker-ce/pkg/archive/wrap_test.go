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
	"testing"
)

func TestGenerateEmptyFile(t *testing.T) {
	archive, err := Generate("emptyFile")
	if err != nil {
		t.Fatal(err)
	}
	if archive == nil {
		t.Fatal("The generated archive should not be nil.")
	}

	expectedFiles := [][]string{
		{"emptyFile", ""},
	}

	tr := tar.NewReader(archive)
	actualFiles := make([][]string, 0, 10)
	i := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(tr)
		content := buf.String()
		actualFiles = append(actualFiles, []string{hdr.Name, content})
		i++
	}
	if len(actualFiles) != len(expectedFiles) {
		t.Fatalf("Number of expected file %d, got %d.", len(expectedFiles), len(actualFiles))
	}
	for i := 0; i < len(expectedFiles); i++ {
		actual := actualFiles[i]
		expected := expectedFiles[i]
		if actual[0] != expected[0] {
			t.Fatalf("Expected name '%s', Actual name '%s'", expected[0], actual[0])
		}
		if actual[1] != expected[1] {
			t.Fatalf("Expected content '%s', Actual content '%s'", expected[1], actual[1])
		}
	}
}

func TestGenerateWithContent(t *testing.T) {
	archive, err := Generate("file", "content")
	if err != nil {
		t.Fatal(err)
	}
	if archive == nil {
		t.Fatal("The generated archive should not be nil.")
	}

	expectedFiles := [][]string{
		{"file", "content"},
	}

	tr := tar.NewReader(archive)
	actualFiles := make([][]string, 0, 10)
	i := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(tr)
		content := buf.String()
		actualFiles = append(actualFiles, []string{hdr.Name, content})
		i++
	}
	if len(actualFiles) != len(expectedFiles) {
		t.Fatalf("Number of expected file %d, got %d.", len(expectedFiles), len(actualFiles))
	}
	for i := 0; i < len(expectedFiles); i++ {
		actual := actualFiles[i]
		expected := expectedFiles[i]
		if actual[0] != expected[0] {
			t.Fatalf("Expected name '%s', Actual name '%s'", expected[0], actual[0])
		}
		if actual[1] != expected[1] {
			t.Fatalf("Expected content '%s', Actual content '%s'", expected[1], actual[1])
		}
	}
}
