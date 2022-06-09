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
package pools

import "github.com/gophercloud/gophercloud"

const (
	rootPath     = "lbaas"
	resourcePath = "pools"
	memberPath   = "members"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func memberRootURL(c *gophercloud.ServiceClient, poolId string) string {
	return c.ServiceURL(rootPath, resourcePath, poolId, memberPath)
}

func memberResourceURL(c *gophercloud.ServiceClient, poolID string, memeberID string) string {
	return c.ServiceURL(rootPath, resourcePath, poolID, memberPath, memeberID)
}
