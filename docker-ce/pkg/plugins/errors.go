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
package plugins

import (
	"fmt"
	"net/http"
)

type statusError struct {
	status int
	method string
	err    string
}

// Error returns a formatted string for this error type
func (e *statusError) Error() string {
	return fmt.Sprintf("%s: %v", e.method, e.err)
}

// IsNotFound indicates if the passed in error is from an http.StatusNotFound from the plugin
func IsNotFound(err error) bool {
	return isStatusError(err, http.StatusNotFound)
}

func isStatusError(err error, status int) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*statusError)
	if !ok {
		return false
	}
	return e.status == status
}
