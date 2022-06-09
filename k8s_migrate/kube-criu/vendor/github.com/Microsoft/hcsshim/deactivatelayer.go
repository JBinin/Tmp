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

import "github.com/sirupsen/logrus"

// DeactivateLayer will dismount a layer that was mounted via ActivateLayer.
func DeactivateLayer(info DriverInfo, id string) error {
	title := "hcsshim::DeactivateLayer "
	logrus.Debugf(title+"Flavour %d ID %s", info.Flavour, id)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = deactivateLayer(&infop, id)
	if err != nil {
		err = makeErrorf(err, title, "id=%s flavour=%d", id, info.Flavour)
		logrus.Error(err)
		return err
	}

	logrus.Debugf(title+"succeeded flavour=%d id=%s", info.Flavour, id)
	return nil
}
