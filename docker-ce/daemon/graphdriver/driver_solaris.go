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

package graphdriver

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
	"unsafe"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/mount"
)

const (
	// FsMagicZfs filesystem id for Zfs
	FsMagicZfs = FsMagic(0x2fc12fc1)
)

var (
	// Slice of drivers that should be used in an order
	priority = []string{
		"zfs",
	}

	// FsNames maps filesystem id to name of the filesystem.
	FsNames = map[FsMagic]string{
		FsMagicZfs: "zfs",
	}
)

// GetFSMagic returns the filesystem id given the path.
func GetFSMagic(rootpath string) (FsMagic, error) {
	return 0, nil
}

type fsChecker struct {
	t FsMagic
}

func (c *fsChecker) IsMounted(path string) bool {
	m, _ := Mounted(c.t, path)
	return m
}

// NewFsChecker returns a checker configured for the provied FsMagic
func NewFsChecker(t FsMagic) Checker {
	return &fsChecker{
		t: t,
	}
}

// NewDefaultChecker returns a check that parses /proc/mountinfo to check
// if the specified path is mounted.
// No-op on Solaris.
func NewDefaultChecker() Checker {
	return &defaultChecker{}
}

type defaultChecker struct {
}

func (c *defaultChecker) IsMounted(path string) bool {
	m, _ := mount.Mounted(path)
	return m
}

// Mounted checks if the given path is mounted as the fs type
//Solaris supports only ZFS for now
func Mounted(fsType FsMagic, mountPath string) (bool, error) {

	cs := C.CString(filepath.Dir(mountPath))
	buf := C.getstatfs(cs)

	// on Solaris buf.f_basetype contains ['z', 'f', 's', 0 ... ]
	if (buf.f_basetype[0] != 122) || (buf.f_basetype[1] != 102) || (buf.f_basetype[2] != 115) ||
		(buf.f_basetype[3] != 0) {
		logrus.Debugf("[zfs] no zfs dataset found for rootdir '%s'", mountPath)
		C.free(unsafe.Pointer(buf))
		return false, ErrPrerequisites
	}

	C.free(unsafe.Pointer(buf))
	C.free(unsafe.Pointer(cs))
	return true, nil
}
