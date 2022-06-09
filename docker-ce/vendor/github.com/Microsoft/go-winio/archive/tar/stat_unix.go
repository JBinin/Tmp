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
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin dragonfly freebsd openbsd netbsd solaris

package tar

import (
	"os"
	"syscall"
)

func init() {
	sysStat = statUnix
}

func statUnix(fi os.FileInfo, h *Header) error {
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return nil
	}
	h.Uid = int(sys.Uid)
	h.Gid = int(sys.Gid)
	// TODO(bradfitz): populate username & group.  os/user
	// doesn't cache LookupId lookups, and lacks group
	// lookup functions.
	h.AccessTime = statAtime(sys)
	h.ChangeTime = statCtime(sys)
	// TODO(bradfitz): major/minor device numbers?
	return nil
}
