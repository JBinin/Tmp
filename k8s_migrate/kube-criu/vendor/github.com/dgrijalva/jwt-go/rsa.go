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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// Implements the RSA family of signing methods signing methods
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for RS256 and company
var (
	SigningMethodRS256 *SigningMethodRSA
	SigningMethodRS384 *SigningMethodRSA
	SigningMethodRS512 *SigningMethodRSA
)

func init() {
	// RS256
	SigningMethodRS256 = &SigningMethodRSA{"RS256", crypto.SHA256}
	RegisterSigningMethod(SigningMethodRS256.Alg(), func() SigningMethod {
		return SigningMethodRS256
	})

	// RS384
	SigningMethodRS384 = &SigningMethodRSA{"RS384", crypto.SHA384}
	RegisterSigningMethod(SigningMethodRS384.Alg(), func() SigningMethod {
		return SigningMethodRS384
	})

	// RS512
	SigningMethodRS512 = &SigningMethodRSA{"RS512", crypto.SHA512}
	RegisterSigningMethod(SigningMethodRS512.Alg(), func() SigningMethod {
		return SigningMethodRS512
	})
}

func (m *SigningMethodRSA) Alg() string {
	return m.Name
}

// Implements the Verify method from SigningMethod
// For this signing method, must be an rsa.PublicKey structure.
func (m *SigningMethodRSA) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = DecodeSegment(signature); err != nil {
		return err
	}

	var rsaKey *rsa.PublicKey
	var ok bool

	if rsaKey, ok = key.(*rsa.PublicKey); !ok {
		return ErrInvalidKeyType
	}

	// Create hasher
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}
	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Verify the signature
	return rsa.VerifyPKCS1v15(rsaKey, m.Hash, hasher.Sum(nil), sig)
}

// Implements the Sign method from SigningMethod
// For this signing method, must be an rsa.PrivateKey structure.
func (m *SigningMethodRSA) Sign(signingString string, key interface{}) (string, error) {
	var rsaKey *rsa.PrivateKey
	var ok bool

	// Validate type of key
	if rsaKey, ok = key.(*rsa.PrivateKey); !ok {
		return "", ErrInvalidKey
	}

	// Create the hasher
	if !m.Hash.Available() {
		return "", ErrHashUnavailable
	}

	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Sign the string and return the encoded bytes
	if sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, hasher.Sum(nil)); err == nil {
		return EncodeSegment(sigBytes), nil
	} else {
		return "", err
	}
}
