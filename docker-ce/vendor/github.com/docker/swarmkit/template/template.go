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
	"strings"
	"text/template"
)

// funcMap defines functions for our template system.
var funcMap = template.FuncMap{
	"join": func(s ...string) string {
		// first arg is sep, remaining args are strings to join
		return strings.Join(s[1:], s[0])
	},
}

func newTemplate(s string) (*template.Template, error) {
	return template.New("expansion").Option("missingkey=error").Funcs(funcMap).Parse(s)
}
