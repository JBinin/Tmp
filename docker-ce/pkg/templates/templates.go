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
	"encoding/json"
	"strings"
	"text/template"
)

// basicFunctions are the set of initial
// functions provided to every template.
var basicFunctions = template.FuncMap{
	"json": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
	"split":    strings.Split,
	"join":     strings.Join,
	"title":    strings.Title,
	"lower":    strings.ToLower,
	"upper":    strings.ToUpper,
	"pad":      padWithSpace,
	"truncate": truncateWithLength,
}

// Parse creates a new anonymous template with the basic functions
// and parses the given format.
func Parse(format string) (*template.Template, error) {
	return NewParse("", format)
}

// NewParse creates a new tagged template with the basic functions
// and parses the given format.
func NewParse(tag, format string) (*template.Template, error) {
	return template.New(tag).Funcs(basicFunctions).Parse(format)
}

// padWithSpace adds whitespace to the input if the input is non-empty
func padWithSpace(source string, prefix, suffix int) string {
	if source == "" {
		return source
	}
	return strings.Repeat(" ", prefix) + source + strings.Repeat(" ", suffix)
}

// truncateWithLength truncates the source string up to the length provided by the input
func truncateWithLength(source string, length int) string {
	if len(source) < length {
		return source
	}
	return source[:length]
}
