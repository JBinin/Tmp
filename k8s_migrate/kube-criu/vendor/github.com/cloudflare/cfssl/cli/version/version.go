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
// Package version implements the version command.
package version

import (
	"fmt"
	"runtime"

	"github.com/cloudflare/cfssl/cli"
)

// Version stores the semantic versioning information for CFSSL.
var version = struct {
	Major    int
	Minor    int
	Patch    int
	Revision string
}{1, 3, 2, "release"}

func versionString() string {
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
}

// Usage text for 'cfssl version'
var versionUsageText = `cfssl version -- print out the version of CF SSL

Usage of version:
	cfssl version
`

// FormatVersion returns the formatted version string.
func FormatVersion() string {
	return fmt.Sprintf("Version: %s\nRevision: %s\nRuntime: %s\n", versionString(), version.Revision, runtime.Version())
}

// The main functionality of 'cfssl version' is to print out the version info.
func versionMain(args []string, c cli.Config) (err error) {
	fmt.Printf("%s", FormatVersion())
	return nil
}

// Command assembles the definition of Command 'version'
var Command = &cli.Command{UsageText: versionUsageText, Flags: nil, Main: versionMain}
