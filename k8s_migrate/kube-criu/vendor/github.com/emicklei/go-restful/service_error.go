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
package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "fmt"

// ServiceError is a transport object to pass information about a non-Http error occurred in a WebService while processing a request.
type ServiceError struct {
	Code    int
	Message string
}

// NewError returns a ServiceError using the code and reason
func NewError(code int, message string) ServiceError {
	return ServiceError{Code: code, Message: message}
}

// Error returns a text representation of the service error
func (s ServiceError) Error() string {
	return fmt.Sprintf("[ServiceError:%v] %v", s.Code, s.Message)
}
