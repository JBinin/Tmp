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

// Implements the none signing method.  This is required by the spec
// but you probably should never use it.
var SigningMethodNone *signingMethodNone

const UnsafeAllowNoneSignatureType unsafeNoneMagicConstant = "none signing method allowed"

var NoneSignatureTypeDisallowedError error

type signingMethodNone struct{}
type unsafeNoneMagicConstant string

func init() {
	SigningMethodNone = &signingMethodNone{}
	NoneSignatureTypeDisallowedError = NewValidationError("'none' signature type is not allowed", ValidationErrorSignatureInvalid)

	RegisterSigningMethod(SigningMethodNone.Alg(), func() SigningMethod {
		return SigningMethodNone
	})
}

func (m *signingMethodNone) Alg() string {
	return "none"
}

// Only allow 'none' alg type if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Verify(signingString, signature string, key interface{}) (err error) {
	// Key must be UnsafeAllowNoneSignatureType to prevent accidentally
	// accepting 'none' signing method
	if _, ok := key.(unsafeNoneMagicConstant); !ok {
		return NoneSignatureTypeDisallowedError
	}
	// If signing method is none, signature must be an empty string
	if signature != "" {
		return NewValidationError(
			"'none' signing method with non-empty signature",
			ValidationErrorSignatureInvalid,
		)
	}

	// Accept 'none' signing method.
	return nil
}

// Only allow 'none' signing if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Sign(signingString string, key interface{}) (string, error) {
	if _, ok := key.(unsafeNoneMagicConstant); ok {
		return "", nil
	}
	return "", NoneSignatureTypeDisallowedError
}
