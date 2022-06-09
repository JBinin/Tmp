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

// ExpandSandboxSize expands the size of a layer to at least size bytes.
func ExpandSandboxSize(info DriverInfo, layerId string, size uint64) error {
	title := "hcsshim::ExpandSandboxSize "
	logrus.Debugf(title+"layerId=%s size=%d", layerId, size)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = expandSandboxSize(&infop, layerId, size)
	if err != nil {
		err = makeErrorf(err, title, "layerId=%s  size=%d", layerId, size)
		logrus.Error(err)
		return err
	}

	logrus.Debugf(title+"- succeeded layerId=%s size=%d", layerId, size)
	return nil
}
