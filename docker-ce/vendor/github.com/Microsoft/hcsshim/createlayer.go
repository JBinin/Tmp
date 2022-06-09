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

// CreateLayer creates a new, empty, read-only layer on the filesystem based on
// the parent layer provided.
func CreateLayer(info DriverInfo, id, parent string) error {
	title := "hcsshim::CreateLayer "
	logrus.Debugf(title+"Flavour %d ID %s parent %s", info.Flavour, id, parent)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = createLayer(&infop, id, parent)
	if err != nil {
		err = makeErrorf(err, title, "id=%s parent=%s flavour=%d", id, parent, info.Flavour)
		logrus.Error(err)
		return err
	}

	logrus.Debugf(title+" - succeeded id=%s parent=%s flavour=%d", id, parent, info.Flavour)
	return nil
}
