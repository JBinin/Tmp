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
Package external provides information and interaction with the external
extension for the OpenStack Networking service.

Example to List Networks with External Information

	type NetworkWithExternalExt struct {
		networks.Network
		external.NetworkExternalExt
	}

	var allNetworks []NetworkWithExternalExt

	allPages, err := networks.List(networkClient, nil).AllPages()
	if err != nil {
		panic(err)
	}

	err = networks.ExtractNetworksInto(allPages, &allNetworks)
	if err != nil {
		panic(err)
	}

	for _, network := range allNetworks {
		fmt.Println("%+v\n", network)
	}

Example to Create a Network with External Information

	iTrue := true
	networkCreateOpts := networks.CreateOpts{
		Name:         "private",
		AdminStateUp: &iTrue,
	}

	createOpts := external.CreateOptsExt{
		networkCreateOpts,
		&iTrue,
	}

	network, err := networks.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package external
