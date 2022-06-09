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
// +build go1.7

package pq

import "crypto/tls"

// Accept renegotiation requests initiated by the backend.
//
// Renegotiation was deprecated then removed from PostgreSQL 9.5, but
// the default configuration of older versions has it enabled. Redshift
// also initiates renegotiations and cannot be reconfigured.
func sslRenegotiation(conf *tls.Config) {
	conf.Renegotiation = tls.RenegotiateFreelyAsClient
}
