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
// +build !go1.7

package tlsconfig

import (
	"crypto/x509"

	"github.com/Sirupsen/logrus"
)

// SystemCertPool returns an new empty cert pool,
// accessing system cert pool is supported in go 1.7
func SystemCertPool() (*x509.CertPool, error) {
	logrus.Warn("Unable to use system certificate pool: requires building with go 1.7 or later")
	return x509.NewCertPool(), nil
}
