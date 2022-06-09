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

// +build dragonfly freebsd linux nacl netbsd openbsd solaris

package x509

import (
	"io/ioutil"
	"os"
)

// Possible directories with certificate files; stop after successfully
// reading at least one file from a directory.
var certDirectories = []string{
	"/etc/ssl/certs",               // SLES10/SLES11, https://golang.org/issue/12139
	"/system/etc/security/cacerts", // Android
	"/usr/local/share/certs",       // FreeBSD
	"/etc/pki/tls/certs",           // Fedora/RHEL
	"/etc/openssl/certs",           // NetBSD
}

const (
	// certFileEnv is the environment variable which identifies where to locate
	// the SSL certificate file. If set this overrides the system default.
	certFileEnv = "SSL_CERT_FILE"

	// certDirEnv is the environment variable which identifies which directory
	// to check for SSL certificate files. If set this overrides the system default.
	certDirEnv = "SSL_CERT_DIR"
)

func (c *Certificate) systemVerify(opts *VerifyOptions) (chains [][]*Certificate, err error) {
	return nil, nil
}

func loadSystemRoots() (*CertPool, error) {
	roots := NewCertPool()

	files := certFiles
	if f := os.Getenv(certFileEnv); f != "" {
		files = []string{f}
	}

	var firstErr error
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err == nil {
			roots.AppendCertsFromPEM(data)
			break
		}
		if firstErr == nil && !os.IsNotExist(err) {
			firstErr = err
		}
	}

	dirs := certDirectories
	if d := os.Getenv(certDirEnv); d != "" {
		dirs = []string{d}
	}

	for _, directory := range dirs {
		fis, err := ioutil.ReadDir(directory)
		if err != nil {
			if firstErr == nil && !os.IsNotExist(err) {
				firstErr = err
			}
			continue
		}
		rootsAdded := false
		for _, fi := range fis {
			data, err := ioutil.ReadFile(directory + "/" + fi.Name())
			if err == nil && roots.AppendCertsFromPEM(data) {
				rootsAdded = true
			}
		}
		if rootsAdded {
			return roots, nil
		}
	}

	if len(roots.certs) > 0 {
		return roots, nil
	}

	return nil, firstErr
}
