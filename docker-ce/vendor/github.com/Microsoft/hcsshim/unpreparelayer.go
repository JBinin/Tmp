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

import "github.com/Sirupsen/logrus"

// UnprepareLayer disables the filesystem filter for the read-write layer with
// the given id.
func UnprepareLayer(info DriverInfo, layerId string) error {
	title := "hcsshim::UnprepareLayer "
	logrus.Debugf(title+"flavour %d layerId %s", info.Flavour, layerId)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = unprepareLayer(&infop, layerId)
	if err != nil {
		err = makeErrorf(err, title, "layerId=%s flavour=%d", layerId, info.Flavour)
		logrus.Error(err)
		return err
	}

	logrus.Debugf(title+"succeeded flavour %d layerId=%s", info.Flavour, layerId)
	return nil
}
