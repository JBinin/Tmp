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
package main

import (
	"strings"
	"text/template"
)

func printArgs(args []arg) string {
	var argStr []string
	for _, arg := range args {
		argStr = append(argStr, arg.String())
	}
	return strings.Join(argStr, ", ")
}

func buildImports(specs []importSpec) string {
	if len(specs) == 0 {
		return `import "errors"`
	}
	imports := "import(\n"
	imports += "\t\"errors\"\n"
	for _, i := range specs {
		imports += "\t" + i.String() + "\n"
	}
	imports += ")"
	return imports
}

func marshalType(t string) string {
	switch t {
	case "error":
		// convert error types to plain strings to ensure the values are encoded/decoded properly
		return "string"
	default:
		return t
	}
}

func isErr(t string) bool {
	switch t {
	case "error":
		return true
	default:
		return false
	}
}

// Need to use this helper due to issues with go-vet
func buildTag(s string) string {
	return "+build " + s
}

var templFuncs = template.FuncMap{
	"printArgs":   printArgs,
	"marshalType": marshalType,
	"isErr":       isErr,
	"lower":       strings.ToLower,
	"title":       title,
	"tag":         buildTag,
	"imports":     buildImports,
}

func title(s string) string {
	if strings.ToLower(s) == "id" {
		return "ID"
	}
	return strings.Title(s)
}

var generatedTempl = template.Must(template.New("rpc_cient").Funcs(templFuncs).Parse(`
// generated code - DO NOT EDIT
{{ range $k, $v := .BuildTags }}
	// {{ tag $k }} {{ end }}

package {{ .Name }}

{{ imports .Imports }}

type client interface{
	Call(string, interface{}, interface{}) error
}

type {{ .InterfaceType }}Proxy struct {
	client
}

{{ range .Functions }}
	type {{ $.InterfaceType }}Proxy{{ .Name }}Request struct{
		{{ range .Args }}
			{{ title .Name }} {{ .ArgType }} {{ end }}
	}

	type {{ $.InterfaceType }}Proxy{{ .Name }}Response struct{
		{{ range .Returns }}
			{{ title .Name }} {{ marshalType .ArgType }} {{ end }}
	}

	func (pp *{{ $.InterfaceType }}Proxy) {{ .Name }}({{ printArgs .Args }}) ({{ printArgs .Returns }}) {
		var(
			req {{ $.InterfaceType }}Proxy{{ .Name }}Request
			ret {{ $.InterfaceType }}Proxy{{ .Name }}Response
		)
		{{ range .Args }}
			req.{{ title .Name }} = {{ lower .Name }} {{ end }}
		if err = pp.Call("{{ $.RPCName }}.{{ .Name }}", req, &ret); err != nil {
			return
		}
		{{ range $r := .Returns }}
			{{ if isErr .ArgType }}
				if ret.{{ title .Name }} != "" {
					{{ lower .Name }} = errors.New(ret.{{ title .Name }})
				} {{ end }}
			{{ if isErr .ArgType | not }} {{ lower .Name }} = ret.{{ title .Name }} {{ end }} {{ end }}

		return
	}
{{ end }}
`))
