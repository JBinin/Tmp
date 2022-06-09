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
package dbus

import (
	"errors"
)

// Call represents a pending or completed method call.
type Call struct {
	Destination string
	Path        ObjectPath
	Method      string
	Args        []interface{}

	// Strobes when the call is complete.
	Done chan *Call

	// After completion, the error status. If this is non-nil, it may be an
	// error message from the peer (with Error as its type) or some other error.
	Err error

	// Holds the response once the call is done.
	Body []interface{}
}

var errSignature = errors.New("dbus: mismatched signature")

// Store stores the body of the reply into the provided pointers. It returns
// an error if the signatures of the body and retvalues don't match, or if
// the error status is not nil.
func (c *Call) Store(retvalues ...interface{}) error {
	if c.Err != nil {
		return c.Err
	}

	return Store(c.Body, retvalues...)
}
