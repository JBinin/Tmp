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
package oci

import (
	"runtime"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// DefaultSpec returns default oci spec used by docker.
func DefaultSpec() specs.Spec {
	s := specs.Spec{
		Version: "0.6.0",
		Platform: specs.Platform{
			OS:   "SunOS",
			Arch: runtime.GOARCH,
		},
	}
	s.Solaris = &specs.Solaris{}
	return s
}
