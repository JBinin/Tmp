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
package flavors

import (
	"github.com/gophercloud/gophercloud"
)

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors", "detail")
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func accessURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "os-flavor-access")
}

func accessActionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "action")
}

func extraSpecsListURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs")
}

func extraSpecsGetURL(client *gophercloud.ServiceClient, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}

func extraSpecsCreateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs")
}

func extraSpecUpdateURL(client *gophercloud.ServiceClient, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}

func extraSpecDeleteURL(client *gophercloud.ServiceClient, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}
