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
Package networks contains functionality for working with Neutron network
resources. A network is an isolated virtual layer-2 broadcast domain that is
typically reserved for the tenant who created it (unless you configure the
network to be shared). Tenants can create multiple networks until the
thresholds per-tenant quota is reached.

In the v2.0 Networking API, the network is the main entity. Ports and subnets
are always associated with a network.

Example to List Networks

	listOpts := networks.ListOpts{
		TenantID: "a99e9b4e620e4db09a2dfb6e42a01e66",
	}

	allPages, err := networks.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		panic(err)
	}

	for _, network := range allNetworks {
		fmt.Printf("%+v", network)
	}

Example to Create a Network

	iTrue := true
	createOpts := networks.CreateOpts{
		Name:         "network_1",
		AdminStateUp: &iTrue,
	}

	network, err := networks.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Network

	networkID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"

	updateOpts := networks.UpdateOpts{
		Name: "new_name",
	}

	network, err := networks.Update(networkClient, networkID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Network

	networkID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := networks.Delete(networkClient, networkID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package networks
