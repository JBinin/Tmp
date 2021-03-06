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
// +build solaris,cgo

package zfs

/*
#include <sys/statvfs.h>
#include <stdlib.h>

static inline struct statvfs *getstatfs(char *s) {
        struct statvfs *buf;
        int err;
        buf = (struct statvfs *)malloc(sizeof(struct statvfs));
        err = statvfs(s, buf);
        return buf;
}
*/
import "C"
import (
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/graphdriver"
)

func checkRootdirFs(rootdir string) error {

	cs := C.CString(filepath.Dir(rootdir))
	buf := C.getstatfs(cs)

	// on Solaris buf.f_basetype contains ['z', 'f', 's', 0 ... ]
	if (buf.f_basetype[0] != 122) || (buf.f_basetype[1] != 102) || (buf.f_basetype[2] != 115) ||
		(buf.f_basetype[3] != 0) {
		logrus.Debugf("[zfs] no zfs dataset found for rootdir '%s'", rootdir)
		C.free(unsafe.Pointer(buf))
		return graphdriver.ErrPrerequisites
	}

	C.free(unsafe.Pointer(buf))
	C.free(unsafe.Pointer(cs))
	return nil
}

/* rootfs is introduced to comply with the OCI spec
which states that root filesystem must be mounted at <CID>/rootfs/ instead of <CID>/
*/
func getMountpoint(id string) string {
	maxlen := 12

	// we need to preserve filesystem suffix
	suffix := strings.SplitN(id, "-", 2)

	if len(suffix) > 1 {
		return filepath.Join(id[:maxlen]+"-"+suffix[1], "rootfs", "root")
	}

	return filepath.Join(id[:maxlen], "rootfs", "root")
}
