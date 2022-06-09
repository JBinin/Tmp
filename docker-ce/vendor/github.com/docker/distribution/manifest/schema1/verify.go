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
package schema1

import (
	"crypto/x509"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libtrust"
)

// Verify verifies the signature of the signed manifest returning the public
// keys used during signing.
func Verify(sm *SignedManifest) ([]libtrust.PublicKey, error) {
	js, err := libtrust.ParsePrettySignature(sm.all, "signatures")
	if err != nil {
		logrus.WithField("err", err).Debugf("(*SignedManifest).Verify")
		return nil, err
	}

	return js.Verify()
}

// VerifyChains verifies the signature of the signed manifest against the
// certificate pool returning the list of verified chains. Signatures without
// an x509 chain are not checked.
func VerifyChains(sm *SignedManifest, ca *x509.CertPool) ([][]*x509.Certificate, error) {
	js, err := libtrust.ParsePrettySignature(sm.all, "signatures")
	if err != nil {
		return nil, err
	}

	return js.VerifyChains(ca)
}
