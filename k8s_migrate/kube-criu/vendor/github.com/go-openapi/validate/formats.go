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
// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validate

import (
	"log"
	"reflect"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
)

type formatValidator struct {
	Format       string
	Path         string
	In           string
	KnownFormats strfmt.Registry
}

func (f *formatValidator) SetPath(path string) {
	f.Path = path
}

func (f *formatValidator) Applies(source interface{}, kind reflect.Kind) bool {
	doit := func() bool {
		if source == nil {
			return false
		}
		switch source.(type) {
		case *spec.Items:
			it := source.(*spec.Items)
			return kind == reflect.String && f.KnownFormats.ContainsName(it.Format)
		case *spec.Parameter:
			par := source.(*spec.Parameter)
			return kind == reflect.String && f.KnownFormats.ContainsName(par.Format)
		case *spec.Schema:
			sch := source.(*spec.Schema)
			return kind == reflect.String && f.KnownFormats.ContainsName(sch.Format)
		}
		return false
	}
	r := doit()
	if Debug {
		log.Printf("format validator for %q applies %t for %T (kind: %v)\n", f.Path, r, source, kind)
	}
	return r
}

func (f *formatValidator) Validate(val interface{}) *Result {
	result := new(Result)

	if err := FormatOf(f.Path, f.In, f.Format, val.(string), f.KnownFormats); err != nil {
		result.AddErrors(err)
	}

	if result.HasErrors() {
		return result
	}
	return nil
}
