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
// +build linux,!btrfs_noversion

package btrfs

/*
#include <btrfs/version.h>

// around version 3.16, they did not define lib version yet
#ifndef BTRFS_LIB_VERSION
#define BTRFS_LIB_VERSION -1
#endif

// upstream had removed it, but now it will be coming back
#ifndef BTRFS_BUILD_VERSION
#define BTRFS_BUILD_VERSION "-"
#endif
*/
import "C"

func btrfsBuildVersion() string {
	return string(C.BTRFS_BUILD_VERSION)
}

func btrfsLibVersion() int {
	return int(C.BTRFS_LIB_VERSION)
}
