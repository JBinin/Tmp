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
package cryptoservice

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/utils"
)

// GenerateCertificate generates an X509 Certificate from a template, given a GUN and validity interval
func GenerateCertificate(rootKey data.PrivateKey, gun string, startTime, endTime time.Time) (*x509.Certificate, error) {
	signer := rootKey.CryptoSigner()
	if signer == nil {
		return nil, fmt.Errorf("key type not supported for Certificate generation: %s\n", rootKey.Algorithm())
	}

	return generateCertificate(signer, gun, startTime, endTime)
}

func generateCertificate(signer crypto.Signer, gun string, startTime, endTime time.Time) (*x509.Certificate, error) {
	template, err := utils.NewCertificate(gun, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create the certificate template for: %s (%v)", gun, err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, signer.Public(), signer)
	if err != nil {
		return nil, fmt.Errorf("failed to create the certificate for: %s (%v)", gun, err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the certificate for key: %s (%v)", gun, err)
	}

	return cert, nil
}
