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
package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaults = map[string]string{
	"FOO": "first",
	"BAR": "",
}

func defaultMapping(name string) (string, bool) {
	val, ok := defaults[name]
	return val, ok
}

func TestEscaped(t *testing.T) {
	result, err := Substitute("$${foo}", defaultMapping)
	assert.NoError(t, err)
	assert.Equal(t, "${foo}", result)
}

func TestInvalid(t *testing.T) {
	invalidTemplates := []string{
		"${",
		"$}",
		"${}",
		"${ }",
		"${ foo}",
		"${foo }",
		"${foo!}",
	}

	for _, template := range invalidTemplates {
		_, err := Substitute(template, defaultMapping)
		assert.Error(t, err)
		assert.IsType(t, &InvalidTemplateError{}, err)
	}
}

func TestNoValueNoDefault(t *testing.T) {
	for _, template := range []string{"This ${missing} var", "This ${BAR} var"} {
		result, err := Substitute(template, defaultMapping)
		assert.NoError(t, err)
		assert.Equal(t, "This  var", result)
	}
}

func TestValueNoDefault(t *testing.T) {
	for _, template := range []string{"This $FOO var", "This ${FOO} var"} {
		result, err := Substitute(template, defaultMapping)
		assert.NoError(t, err)
		assert.Equal(t, "This first var", result)
	}
}

func TestNoValueWithDefault(t *testing.T) {
	for _, template := range []string{"ok ${missing:-def}", "ok ${missing-def}"} {
		result, err := Substitute(template, defaultMapping)
		assert.NoError(t, err)
		assert.Equal(t, "ok def", result)
	}
}

func TestEmptyValueWithSoftDefault(t *testing.T) {
	result, err := Substitute("ok ${BAR:-def}", defaultMapping)
	assert.NoError(t, err)
	assert.Equal(t, "ok def", result)
}

func TestEmptyValueWithHardDefault(t *testing.T) {
	result, err := Substitute("ok ${BAR-def}", defaultMapping)
	assert.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

func TestNonAlphanumericDefault(t *testing.T) {
	result, err := Substitute("ok ${BAR:-/non:-alphanumeric}", defaultMapping)
	assert.NoError(t, err)
	assert.Equal(t, "ok /non:-alphanumeric", result)
}
