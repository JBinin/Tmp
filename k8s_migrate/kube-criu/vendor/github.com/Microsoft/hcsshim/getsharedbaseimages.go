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

// GetSharedBaseImages will enumerate the images stored in the common central
// image store and return descriptive info about those images for the purpose
// of registering them with the graphdriver, graph, and tagstore.
func GetSharedBaseImages() (imageData string, err error) {
	title := "hcsshim::GetSharedBaseImages "

	logrus.Debugf("Calling proc")
	var buffer *uint16
	err = getBaseImages(&buffer)
	if err != nil {
		err = makeError(err, title, "")
		logrus.Error(err)
		return
	}
	imageData = convertAndFreeCoTaskMemString(buffer)
	logrus.Debugf(title+" - succeeded output=%s", imageData)
	return
}
