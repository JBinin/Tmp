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
package credentials

// ErrCredentialsNotFound standarizes the not found error, so every helper returns
// the same message and docker can handle it properly.
const errCredentialsNotFoundMessage = "credentials not found in native keychain"

// errCredentialsNotFound represents an error
// raised when credentials are not in the store.
type errCredentialsNotFound struct{}

// Error returns the standard error message
// for when the credentials are not in the store.
func (errCredentialsNotFound) Error() string {
	return errCredentialsNotFoundMessage
}

// NewErrCredentialsNotFound creates a new error
// for when the credentials are not in the store.
func NewErrCredentialsNotFound() error {
	return errCredentialsNotFound{}
}

// IsErrCredentialsNotFound returns true if the error
// was caused by not having a set of credentials in a store.
func IsErrCredentialsNotFound(err error) bool {
	_, ok := err.(errCredentialsNotFound)
	return ok
}

// IsErrCredentialsNotFoundMessage returns true if the error
// was caused by not having a set of credentials in a store.
//
// This function helps to check messages returned by an
// external program via its standard output.
func IsErrCredentialsNotFoundMessage(err string) bool {
	return err == errCredentialsNotFoundMessage
}
