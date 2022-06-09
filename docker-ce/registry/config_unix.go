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
// +build !windows

package registry

import (
	"github.com/spf13/pflag"
)

var (
	// CertsDir is the directory where certificates are stored
	CertsDir = "/etc/docker/certs.d"
)

// cleanPath is used to ensure that a directory name is valid on the target
// platform. It will be passed in something *similar* to a URL such as
// https:/index.docker.io/v1. Not all platforms support directory names
// which contain those characters (such as : on Windows)
func cleanPath(s string) string {
	return s
}

// installCliPlatformFlags handles any platform specific flags for the service.
func (options *ServiceOptions) installCliPlatformFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&options.V2Only, "disable-legacy-registry", false, "Disable contacting legacy registries")
}
