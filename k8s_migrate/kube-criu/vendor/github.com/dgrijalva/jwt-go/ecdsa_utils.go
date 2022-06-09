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
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrNotECPublicKey  = errors.New("Key is not a valid ECDSA public key")
	ErrNotECPrivateKey = errors.New("Key is not a valid ECDSA private key")
)

// Parse PEM encoded Elliptic Curve Private Key Structure
func ParseECPrivateKeyFromPEM(key []byte) (*ecdsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		return nil, err
	}

	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, ErrNotECPrivateKey
	}

	return pkey, nil
}

// Parse PEM encoded PKCS1 or PKCS8 public key
func ParseECPublicKeyFromPEM(key []byte) (*ecdsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *ecdsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PublicKey); !ok {
		return nil, ErrNotECPublicKey
	}

	return pkey, nil
}
