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
package volumeattach

import "github.com/gophercloud/gophercloud"

const resourcePath = "os-volume_attachments"

func resourceURL(c *gophercloud.ServiceClient, serverID string) string {
	return c.ServiceURL("servers", serverID, resourcePath)
}

func listURL(c *gophercloud.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func createURL(c *gophercloud.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func getURL(c *gophercloud.ServiceClient, serverID, aID string) string {
	return c.ServiceURL("servers", serverID, resourcePath, aID)
}

func deleteURL(c *gophercloud.ServiceClient, serverID, aID string) string {
	return getURL(c, serverID, aID)
}
