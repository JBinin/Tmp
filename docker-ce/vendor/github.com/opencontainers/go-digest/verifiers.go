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
package digest

import (
	"hash"
	"io"
)

// Verifier presents a general verification interface to be used with message
// digests and other byte stream verifications. Users instantiate a Verifier
// from one of the various methods, write the data under test to it then check
// the result with the Verified method.
type Verifier interface {
	io.Writer

	// Verified will return true if the content written to Verifier matches
	// the digest.
	Verified() bool
}

type hashVerifier struct {
	digest Digest
	hash   hash.Hash
}

func (hv hashVerifier) Write(p []byte) (n int, err error) {
	return hv.hash.Write(p)
}

func (hv hashVerifier) Verified() bool {
	return hv.digest == NewDigest(hv.digest.Algorithm(), hv.hash)
}
