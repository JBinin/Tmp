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
package main

import (
	"strings"
)

// LicenseFilePrefix is a list of filename prefixes that indicate it
//  might contain a software license
var LicenseFilePrefix = []string{
	"licence", // UK spelling
	"license", // US spelling
	"copying",
	"unlicense",
	"copyright",
	"copyleft",
	"authors",
	"contributors",
}

// LegalFileSubstring are substrings that indicate the file is likely
// to contain some type of legal declaration.  "legal" is often used
// that it might moved to LicenseFilePrefix
var LegalFileSubstring = []string{
	"legal",
	"notice",
	"disclaimer",
	"patent",
	"third-party",
	"thirdparty",
}

// IsLicenseFile returns true if the filename might be contain a
// software license
func IsLicenseFile(filename string) bool {
	lowerfile := strings.ToLower(filename)
	for _, prefix := range LicenseFilePrefix {
		if strings.HasPrefix(lowerfile, prefix) {
			return true
		}
	}
	return false
}

// IsLegalFile returns true if the file is likely to contain some type
// of of legal declaration or licensing information
func IsLegalFile(filename string) bool {
	lowerfile := strings.ToLower(filename)
	for _, prefix := range LicenseFilePrefix {
		if strings.HasPrefix(lowerfile, prefix) {
			return true
		}
	}
	for _, substring := range LegalFileSubstring {
		if strings.Contains(lowerfile, substring) {
			return true
		}
	}
	return false
}
