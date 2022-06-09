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
// +build !go1.6

package testsuite

import (
	"os"
	"path"
)

// "This change will only be enabled if the go command is run with
// GO15VENDOREXPERIMENT=1 in its environment."
// c.f. the vendor-experiment proposal https://goo.gl/2ucMeC
func vendorExperimentCheck(dir string) bool {
	vendorExperiment := os.Getenv("GO15VENDOREXPERIMENT")
	return vendorExperiment == "1" && path.Base(dir) == "vendor"
}
