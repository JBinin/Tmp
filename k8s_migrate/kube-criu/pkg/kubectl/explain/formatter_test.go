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

package explain

import (
	"bytes"
	"testing"
)

func TestFormatterWrite(t *testing.T) {
	buf := bytes.Buffer{}
	f := Formatter{
		Writer: &buf,
	}

	f.Write("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	// Indent creates a new Formatter
	f.Indent(5).Write("Morbi at turpis faucibus, gravida dolor ut, fringilla velit.")
	// So Indent(2) doesn't indent to 7 here.
	f.Indent(2).Write("Etiam maximus urna at tellus faucibus mattis.")

	want := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
     Morbi at turpis faucibus, gravida dolor ut, fringilla velit.
  Etiam maximus urna at tellus faucibus mattis.
`

	if buf.String() != want {
		t.Errorf("Got:\n%v\nWant:\n%v\n", buf.String(), want)
	}
}

func TestFormatterWrappedWrite(t *testing.T) {
	buf := bytes.Buffer{}
	f := Formatter{
		Writer: &buf,
		Wrap:   50,
	}

	f.WriteWrapped("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi at turpis faucibus, gravida dolor ut, fringilla velit.")
	f.Indent(10).WriteWrapped("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi at turpis faucibus, gravida dolor ut, fringilla velit.")
	// Test long words (especially urls) on their own line.
	f.Indent(20).WriteWrapped("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ThisIsAVeryLongWordThatDoesn'tFitOnALineOnItsOwn. Morbi at turpis faucibus, gravida dolor ut, fringilla velit.")

	want := `Lorem ipsum dolor sit amet, consectetur adipiscing
elit. Morbi at turpis faucibus, gravida dolor ut,
fringilla velit.
          Lorem ipsum dolor sit amet, consectetur
          adipiscing elit. Morbi at turpis
          faucibus, gravida dolor ut, fringilla
          velit.
                    Lorem ipsum dolor sit amet,
                    consectetur adipiscing elit.
                    ThisIsAVeryLongWordThatDoesn'tFitOnALineOnItsOwn.
                    Morbi at turpis faucibus,
                    gravida dolor ut, fringilla
                    velit.
`

	if buf.String() != want {
		t.Errorf("Got:\n%v\nWant:\n%v\n", buf.String(), want)
	}
}

func TestDefaultWrap(t *testing.T) {
	buf := bytes.Buffer{}
	f := Formatter{
		Writer: &buf,
		// Wrap is not set
	}

	f.WriteWrapped("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi at turpis faucibus, gravida dolor ut, fringilla velit. Etiam maximus urna at tellus faucibus mattis.")
	want := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi at turpis faucibus, gravida dolor ut, fringilla velit. Etiam maximus urna at tellus faucibus mattis.
`
	if buf.String() != want {
		t.Errorf("Got:\n%v\nWant:\n%v\n", buf.String(), want)
	}
}
