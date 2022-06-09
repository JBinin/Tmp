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
package serror

import (
	"encoding/json"
)

type StorageOSError interface {
	// embedding error provides compatibility with standard error handling code
	error

	// Encoding/decoding methods to help errors traverse API boundaries
	json.Marshaler
	json.Unmarshaler

	Err() error               // Returns the underlying error that caused this event
	String() string           // A short string representing the error (for logging etc)
	Help() string             // A larger string that should provide informative debug instruction to users
	Kind() StorageOSErrorKind // A type representing a set of known error conditions, helpful to switch on
	Extra() map[string]string // A container for error specific information

	// TODO: should we include callstack traces here? We could have a debug mode for it.
}

func ErrorKind(err error) StorageOSErrorKind {
	if serr, ok := err.(StorageOSError); ok {
		return serr.Kind()
	}
	return UnknownError
}

func IsStorageOSError(err error) bool {
	_, ok := err.(StorageOSError)
	return ok
}
