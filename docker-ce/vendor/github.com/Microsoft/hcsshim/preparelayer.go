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
	"sync"

	"github.com/Sirupsen/logrus"
)

var prepareLayerLock sync.Mutex

// PrepareLayer finds a mounted read-write layer matching layerId and enables the
// the filesystem filter for use on that layer.  This requires the paths to all
// parent layers, and is necessary in order to view or interact with the layer
// as an actual filesystem (reading and writing files, creating directories, etc).
// Disabling the filter must be done via UnprepareLayer.
func PrepareLayer(info DriverInfo, layerId string, parentLayerPaths []string) error {
	title := "hcsshim::PrepareLayer "
	logrus.Debugf(title+"flavour %d layerId %s", info.Flavour, layerId)

	// Generate layer descriptors
	layers, err := layerPathsToDescriptors(parentLayerPaths)
	if err != nil {
		return err
	}

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// This lock is a temporary workaround for a Windows bug. Only allowing one
	// call to prepareLayer at a time vastly reduces the chance of a timeout.
	prepareLayerLock.Lock()
	defer prepareLayerLock.Unlock()
	err = prepareLayer(&infop, layerId, layers)
	if err != nil {
		err = makeErrorf(err, title, "layerId=%s flavour=%d", layerId, info.Flavour)
		logrus.Error(err)
		return err
	}

	logrus.Debugf(title+"succeeded flavour=%d layerId=%s", info.Flavour, layerId)
	return nil
}
