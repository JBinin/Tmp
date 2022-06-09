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
package jwt

import (
	"errors"
)

// Error constants
var (
	ErrInvalidKey      = errors.New("key is invalid")
	ErrInvalidKeyType  = errors.New("key is of invalid type")
	ErrHashUnavailable = errors.New("the requested hash function is unavailable")
)

// The errors that might occur when parsing and validating a token
const (
	ValidationErrorMalformed        uint32 = 1 << iota // Token is malformed
	ValidationErrorUnverifiable                        // Token could not be verified because of signing problems
	ValidationErrorSignatureInvalid                    // Signature validation failed

	// Standard Claim validation errors
	ValidationErrorAudience      // AUD validation failed
	ValidationErrorExpired       // EXP validation failed
	ValidationErrorIssuedAt      // IAT validation failed
	ValidationErrorIssuer        // ISS validation failed
	ValidationErrorNotValidYet   // NBF validation failed
	ValidationErrorId            // JTI validation failed
	ValidationErrorClaimsInvalid // Generic claims validation error
)

// Helper for constructing a ValidationError with a string error message
func NewValidationError(errorText string, errorFlags uint32) *ValidationError {
	return &ValidationError{
		text:   errorText,
		Errors: errorFlags,
	}
}

// The error from Parse if token is not valid
type ValidationError struct {
	Inner  error  // stores the error returned by external dependencies, i.e.: KeyFunc
	Errors uint32 // bitfield.  see ValidationError... constants
	text   string // errors that do not have a valid error just have text
}

// Validation error is an error type
func (e ValidationError) Error() string {
	if e.Inner != nil {
		return e.Inner.Error()
	} else if e.text != "" {
		return e.text
	} else {
		return "token is invalid"
	}
	return e.Inner.Error()
}

// No errors
func (e *ValidationError) valid() bool {
	if e.Errors > 0 {
		return false
	}
	return true
}
