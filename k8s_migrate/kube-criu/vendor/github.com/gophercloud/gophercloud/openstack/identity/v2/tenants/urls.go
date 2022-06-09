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
package tenants

import "github.com/gophercloud/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("tenants")
}

func getURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("tenants")
}

func deleteURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func updateURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}
