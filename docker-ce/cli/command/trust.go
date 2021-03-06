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
package command

import (
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

var (
	// TODO: make this not global
	untrusted bool
)

// AddTrustVerificationFlags adds content trust flags to the provided flagset
func AddTrustVerificationFlags(fs *pflag.FlagSet) {
	trusted := getDefaultTrustState()
	fs.BoolVar(&untrusted, "disable-content-trust", !trusted, "Skip image verification")
}

// AddTrustSigningFlags adds "signing" flags to the provided flagset
func AddTrustSigningFlags(fs *pflag.FlagSet) {
	trusted := getDefaultTrustState()
	fs.BoolVar(&untrusted, "disable-content-trust", !trusted, "Skip image signing")
}

// getDefaultTrustState returns true if content trust is enabled through the $DOCKER_CONTENT_TRUST environment variable.
func getDefaultTrustState() bool {
	var trusted bool
	if e := os.Getenv("DOCKER_CONTENT_TRUST"); e != "" {
		if t, err := strconv.ParseBool(e); t || err != nil {
			// treat any other value as true
			trusted = true
		}
	}
	return trusted
}

// IsTrusted returns true if content trust is enabled, either through the $DOCKER_CONTENT_TRUST environment variable,
// or through `--disabled-content-trust=false` on a command.
func IsTrusted() bool {
	return !untrusted
}
