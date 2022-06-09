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
package dhcp4client

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
)

func CryptoGenerateXID(b []byte) {
	if _, err := cryptorand.Read(b); err != nil {
		panic(err)
	}
}

func MathGenerateXID(b []byte) {
	if _, err := mathrand.Read(b); err != nil {
		panic(err)
	}
}
