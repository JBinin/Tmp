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
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build dragonfly freebsd linux openbsd netbsd solaris

package x509

import "io/ioutil"

// Possible certificate files; stop after finding one.
var certFiles = []string{
	"/etc/ssl/certs/ca-certificates.crt",     // Debian/Ubuntu/Gentoo etc.
	"/etc/pki/tls/certs/ca-bundle.crt",       // Fedora/RHEL
	"/etc/ssl/ca-bundle.pem",                 // OpenSUSE
	"/etc/ssl/cert.pem",                      // OpenBSD
	"/usr/local/share/certs/ca-root-nss.crt", // FreeBSD/DragonFly
}

func (c *Certificate) systemVerify(opts *VerifyOptions) (chains [][]*Certificate, err error) {
	return nil, nil
}

func initSystemRoots() {
	roots := NewCertPool()
	for _, file := range certFiles {
		data, err := ioutil.ReadFile(file)
		if err == nil {
			roots.AppendCertsFromPEM(data)
			systemRoots = roots
			return
		}
	}

	// All of the files failed to load. systemRoots will be nil which will
	// trigger a specific error at verification time.
}
