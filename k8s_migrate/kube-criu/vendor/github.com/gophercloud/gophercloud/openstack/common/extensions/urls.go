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
package extensions

import "github.com/gophercloud/gophercloud"

// ExtensionURL generates the URL for an extension resource by name.
func ExtensionURL(c *gophercloud.ServiceClient, name string) string {
	return c.ServiceURL("extensions", name)
}

// ListExtensionURL generates the URL for the extensions resource collection.
func ListExtensionURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("extensions")
}
