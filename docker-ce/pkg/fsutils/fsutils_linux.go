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
// +build linux

package fsutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
)

func locateDummyIfEmpty(path string) (string, error) {
	children, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	if len(children) != 0 {
		return "", nil
	}
	dummyFile, err := ioutil.TempFile(path, "fsutils-dummy")
	if err != nil {
		return "", err
	}
	name := dummyFile.Name()
	err = dummyFile.Close()
	return name, err
}

// SupportsDType returns whether the filesystem mounted on path supports d_type
func SupportsDType(path string) (bool, error) {
	// locate dummy so that we have at least one dirent
	dummy, err := locateDummyIfEmpty(path)
	if err != nil {
		return false, err
	}
	if dummy != "" {
		defer os.Remove(dummy)
	}

	visited := 0
	supportsDType := true
	fn := func(ent *syscall.Dirent) bool {
		visited++
		if ent.Type == syscall.DT_UNKNOWN {
			supportsDType = false
			// stop iteration
			return true
		}
		// continue iteration
		return false
	}
	if err = iterateReadDir(path, fn); err != nil {
		return false, err
	}
	if visited == 0 {
		return false, fmt.Errorf("did not hit any dirent during iteration %s", path)
	}
	return supportsDType, nil
}

func iterateReadDir(path string, fn func(*syscall.Dirent) bool) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()
	fd := int(d.Fd())
	buf := make([]byte, 4096)
	for {
		nbytes, err := syscall.ReadDirent(fd, buf)
		if err != nil {
			return err
		}
		if nbytes == 0 {
			break
		}
		for off := 0; off < nbytes; {
			ent := (*syscall.Dirent)(unsafe.Pointer(&buf[off]))
			if stop := fn(ent); stop {
				return nil
			}
			off += int(ent.Reclen)
		}
	}
	return nil
}
