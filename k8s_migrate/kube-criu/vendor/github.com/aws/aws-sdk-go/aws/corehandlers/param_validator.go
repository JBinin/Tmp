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
package corehandlers

import "github.com/aws/aws-sdk-go/aws/request"

// ValidateParametersHandler is a request handler to validate the input parameters.
// Validating parameters only has meaning if done prior to the request being sent.
var ValidateParametersHandler = request.NamedHandler{Name: "core.ValidateParametersHandler", Fn: func(r *request.Request) {
	if !r.ParamsFilled() {
		return
	}

	if v, ok := r.Params.(request.Validator); ok {
		if err := v.Validate(); err != nil {
			r.Error = err
		}
	}
}}
