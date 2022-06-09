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
/*
Package openstack contains resources for the individual OpenStack projects
supported in Gophercloud. It also includes functions to authenticate to an
OpenStack cloud and for provisioning various service-level clients.

Example of Creating a Service Client

	ao, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.AuthenticatedClient(ao)
	client, err := openstack.NewNetworkV2(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
*/
package openstack
