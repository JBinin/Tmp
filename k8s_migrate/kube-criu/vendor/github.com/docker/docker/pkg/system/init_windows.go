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
package system

import "os"

// LCOWSupported determines if Linux Containers on Windows are supported.
// Note: This feature is in development (06/17) and enabled through an
// environment variable. At a future time, it will be enabled based
// on build number. @jhowardmsft
var lcowSupported = false

func init() {
	// LCOW initialization
	if os.Getenv("LCOW_SUPPORTED") != "" {
		lcowSupported = true
	}

}
