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
package templates

import (
	"bytes"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

func TestParseStringFunctions(t *testing.T) {
	tm, err := Parse(`{{join (split . ":") "/"}}`)
	assert.NilError(t, err)

	var b bytes.Buffer
	assert.NilError(t, tm.Execute(&b, "text:with:colon"))
	want := "text/with/colon"
	assert.Equal(t, b.String(), want)
}

func TestNewParse(t *testing.T) {
	tm, err := NewParse("foo", "this is a {{ . }}")
	assert.NilError(t, err)

	var b bytes.Buffer
	assert.NilError(t, tm.Execute(&b, "string"))
	want := "this is a string"
	assert.Equal(t, b.String(), want)
}

func TestParseTruncateFunction(t *testing.T) {
	source := "tupx5xzf6hvsrhnruz5cr8gwp"

	testCases := []struct {
		template string
		expected string
	}{
		{
			template: `{{truncate . 5}}`,
			expected: "tupx5",
		},
		{
			template: `{{truncate . 25}}`,
			expected: "tupx5xzf6hvsrhnruz5cr8gwp",
		},
		{
			template: `{{truncate . 30}}`,
			expected: "tupx5xzf6hvsrhnruz5cr8gwp",
		},
	}

	for _, testCase := range testCases {
		tm, err := Parse(testCase.template)
		assert.NilError(t, err)

		var b bytes.Buffer
		assert.NilError(t, tm.Execute(&b, source))
		assert.Equal(t, b.String(), testCase.expected)
	}
}
