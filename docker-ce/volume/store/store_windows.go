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
package store

import "strings"

// normaliseVolumeName is a platform specific function to normalise the name
// of a volume. On Windows, as NTFS is case insensitive, under
// c:\ProgramData\Docker\Volumes\, the folders John and john would be synonymous.
// Hence we can't allow the volume "John" and "john" to be created as separate
// volumes.
func normaliseVolumeName(name string) string {
	return strings.ToLower(name)
}
