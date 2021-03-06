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

package errors

import (
	"fmt"
	"net/http"
)

// Validation represents a failure of a precondition
type Validation struct {
	code    int32
	Name    string
	In      string
	Value   interface{}
	message string
	Values  []interface{}
}

func (e *Validation) Error() string {
	return e.message
}

// Code the error code
func (e *Validation) Code() int32 {
	return e.code
}

const (
	contentTypeFail    = `unsupported media type %q, only %v are allowed`
	responseFormatFail = `unsupported media type requested, only %v are available`
)

// InvalidContentType error for an invalid content type
func InvalidContentType(value string, allowed []string) *Validation {
	var values []interface{}
	for _, v := range allowed {
		values = append(values, v)
	}
	return &Validation{
		code:    http.StatusUnsupportedMediaType,
		Name:    "Content-Type",
		In:      "header",
		Value:   value,
		Values:  values,
		message: fmt.Sprintf(contentTypeFail, value, allowed),
	}
}

// InvalidResponseFormat error for an unacceptable response format request
func InvalidResponseFormat(value string, allowed []string) *Validation {
	var values []interface{}
	for _, v := range allowed {
		values = append(values, v)
	}
	return &Validation{
		code:    http.StatusNotAcceptable,
		Name:    "Accept",
		In:      "header",
		Value:   value,
		Values:  values,
		message: fmt.Sprintf(responseFormatFail, allowed),
	}
}
