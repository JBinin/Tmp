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

// NameToGuid converts the given string into a GUID using the algorithm in the
// Host Compute Service, ensuring GUIDs generated with the same string are common
// across all clients.
func NameToGuid(name string) (id GUID, err error) {
	title := "hcsshim::NameToGuid "
	logrus.Debugf(title+"Name %s", name)

	err = nameToGuid(name, &id)
	if err != nil {
		err = makeErrorf(err, title, "name=%s", name)
		logrus.Error(err)
		return
	}

	return
}
