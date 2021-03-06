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
Copyright 2017 The Kubernetes Authors.

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

package util

import (
	"testing"
)

const (
	validTmpl    = "image: {{ .ImageRepository }}/pause-{{ .Arch }}:3.1"
	validTmplOut = "image: k8s.gcr.io/pause-amd64:3.1"
	doNothing    = "image: k8s.gcr.io/pause-amd64:3.1"
	invalidTmpl1 = "{{ .baz }/d}"
	invalidTmpl2 = "{{ !foobar }}"
)

func TestParseTemplate(t *testing.T) {
	var tmplTests = []struct {
		template    string
		data        interface{}
		output      string
		errExpected bool
	}{
		// should parse a valid template and set the right values
		{
			template: validTmpl,
			data: struct{ ImageRepository, Arch string }{
				ImageRepository: "k8s.gcr.io",
				Arch:            "amd64",
			},
			output:      validTmplOut,
			errExpected: false,
		},
		// should noop if there aren't any {{ .foo }} present
		{
			template: doNothing,
			data: struct{ ImageRepository, Arch string }{
				ImageRepository: "k8s.gcr.io",
				Arch:            "amd64",
			},
			output:      doNothing,
			errExpected: false,
		},
		// invalid syntax, passing nil
		{
			template:    invalidTmpl1,
			data:        nil,
			output:      "",
			errExpected: true,
		},
		// invalid syntax
		{
			template:    invalidTmpl2,
			data:        struct{}{},
			output:      "",
			errExpected: true,
		},
	}
	for _, tt := range tmplTests {
		outbytes, err := ParseTemplate(tt.template, tt.data)
		if tt.errExpected != (err != nil) {
			t.Errorf(
				"failed TestParseTemplate:\n\texpected err: %t\n\t  actual: %s",
				tt.errExpected,
				err,
			)
		}
		if tt.output != string(outbytes) {
			t.Errorf(
				"failed TestParseTemplate:\n\texpected bytes: %s\n\t  actual: %s",
				tt.output,
				outbytes,
			)
		}
	}
}
