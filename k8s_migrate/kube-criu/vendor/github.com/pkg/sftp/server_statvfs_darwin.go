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
package sftp

import (
	"syscall"
)

func statvfsFromStatfst(stat *syscall.Statfs_t) (*StatVFS, error) {
	return &StatVFS{
		Bsize:   uint64(stat.Bsize),
		Frsize:  uint64(stat.Bsize), // fragment size is a linux thing; use block size here
		Blocks:  stat.Blocks,
		Bfree:   stat.Bfree,
		Bavail:  stat.Bavail,
		Files:   stat.Files,
		Ffree:   stat.Ffree,
		Favail:  stat.Ffree,                                                      // not sure how to calculate Favail
		Fsid:    uint64(uint64(stat.Fsid.Val[1])<<32 | uint64(stat.Fsid.Val[0])), // endianness?
		Flag:    uint64(stat.Flags),                                              // assuming POSIX?
		Namemax: 1024,                                                            // man 2 statfs shows: #define MAXPATHLEN      1024
	}, nil
}
