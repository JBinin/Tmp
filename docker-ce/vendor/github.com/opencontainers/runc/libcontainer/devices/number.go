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
// +build linux freebsd

package devices

/*

This code provides support for manipulating linux device numbers.  It should be replaced by normal syscall functions once http://code.google.com/p/go/issues/detail?id=8106 is solved.

You can read what they are here:

 - http://www.makelinux.net/ldd3/chp-3-sect-2
 - http://www.linux-tutorial.info/modules.php?name=MContent&pageid=94

Note! These are NOT the same as the MAJOR(dev_t device);, MINOR(dev_t device); and MKDEV(int major, int minor); functions as defined in <linux/kdev_t.h> as the representation of device numbers used by go is different than the one used internally to the kernel! - https://github.com/torvalds/linux/blob/master/include/linux/kdev_t.h#L9

*/

func Major(devNumber int) int64 {
	return int64((devNumber >> 8) & 0xfff)
}

func Minor(devNumber int) int64 {
	return int64((devNumber & 0xff) | ((devNumber >> 12) & 0xfff00))
}
