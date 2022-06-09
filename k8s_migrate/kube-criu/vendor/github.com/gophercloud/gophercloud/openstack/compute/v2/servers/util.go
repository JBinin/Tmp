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
package servers

import "github.com/gophercloud/gophercloud"

// WaitForStatus will continually poll a server until it successfully
// transitions to a specified status. It will do this for at most the number
// of seconds specified.
func WaitForStatus(c *gophercloud.ServiceClient, id, status string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := Get(c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}
