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
Copyright 2016 The Kubernetes Authors.

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

package generators

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

// genExpansion produces a file for a group client, e.g. ExtensionsClient for the extension group.
type genExpansion struct {
	generator.DefaultGen
	groupPackagePath string
	// types in a group
	types []*types.Type
}

// We only want to call GenerateType() once per group.
func (g *genExpansion) Filter(c *generator.Context, t *types.Type) bool {
	return len(g.types) == 0 || t == g.types[0]
}

func (g *genExpansion) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	for _, t := range g.types {
		if _, err := os.Stat(filepath.Join(g.groupPackagePath, strings.ToLower(t.Name.Name+"_expansion.go"))); os.IsNotExist(err) {
			sw.Do(expansionInterfaceTemplate, t)
		}
	}
	return sw.Error()
}

var expansionInterfaceTemplate = `
type $.|public$Expansion interface {}
`
