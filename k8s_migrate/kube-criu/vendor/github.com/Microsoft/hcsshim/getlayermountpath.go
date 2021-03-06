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
package hcsshim

import (
	"syscall"

	"github.com/sirupsen/logrus"
)

// GetLayerMountPath will look for a mounted layer with the given id and return
// the path at which that layer can be accessed.  This path may be a volume path
// if the layer is a mounted read-write layer, otherwise it is expected to be the
// folder path at which the layer is stored.
func GetLayerMountPath(info DriverInfo, id string) (string, error) {
	title := "hcsshim::GetLayerMountPath "
	logrus.Debugf(title+"Flavour %d ID %s", info.Flavour, id)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	var mountPathLength uintptr
	mountPathLength = 0

	// Call the procedure itself.
	logrus.Debugf("Calling proc (1)")
	err = getLayerMountPath(&infop, id, &mountPathLength, nil)
	if err != nil {
		err = makeErrorf(err, title, "(first call) id=%s flavour=%d", id, info.Flavour)
		logrus.Error(err)
		return "", err
	}

	// Allocate a mount path of the returned length.
	if mountPathLength == 0 {
		return "", nil
	}
	mountPathp := make([]uint16, mountPathLength)
	mountPathp[0] = 0

	// Call the procedure again
	logrus.Debugf("Calling proc (2)")
	err = getLayerMountPath(&infop, id, &mountPathLength, &mountPathp[0])
	if err != nil {
		err = makeErrorf(err, title, "(second call) id=%s flavour=%d", id, info.Flavour)
		logrus.Error(err)
		return "", err
	}

	path := syscall.UTF16ToString(mountPathp[0:])
	logrus.Debugf(title+"succeeded flavour=%d id=%s path=%s", info.Flavour, id, path)
	return path, nil
}
