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
package errors

import "net/http"

// apiError is an error wrapper that also
// holds information about response status codes.
type apiError struct {
	error
	statusCode int
}

// HTTPErrorStatusCode returns a status code.
func (e apiError) HTTPErrorStatusCode() int {
	return e.statusCode
}

// NewErrorWithStatusCode allows you to associate
// a specific HTTP Status Code to an error.
// The server will take that code and set
// it as the response status.
func NewErrorWithStatusCode(err error, code int) error {
	return apiError{err, code}
}

// NewBadRequestError creates a new API error
// that has the 400 HTTP status code associated to it.
func NewBadRequestError(err error) error {
	return NewErrorWithStatusCode(err, http.StatusBadRequest)
}

// NewRequestForbiddenError creates a new API error
// that has the 403 HTTP status code associated to it.
func NewRequestForbiddenError(err error) error {
	return NewErrorWithStatusCode(err, http.StatusForbidden)
}

// NewRequestNotFoundError creates a new API error
// that has the 404 HTTP status code associated to it.
func NewRequestNotFoundError(err error) error {
	return NewErrorWithStatusCode(err, http.StatusNotFound)
}

// NewRequestConflictError creates a new API error
// that has the 409 HTTP status code associated to it.
func NewRequestConflictError(err error) error {
	return NewErrorWithStatusCode(err, http.StatusConflict)
}
