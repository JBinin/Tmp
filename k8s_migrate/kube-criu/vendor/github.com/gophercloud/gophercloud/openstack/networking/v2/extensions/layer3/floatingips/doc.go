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
package floatingips enables management and retrieval of Floating IPs from the
OpenStack Networking service.

Example to List Floating IPs

	listOpts := floatingips.ListOpts{
		FloatingNetworkID: "a6917946-38ab-4ffd-a55a-26c0980ce5ee",
	}

	allPages, err := floatingips.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allFIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		panic(err)
	}

	for _, fip := range allFIPs {
		fmt.Printf("%+v\n", fip)
	}

Example to Create a Floating IP

	createOpts := floatingips.CreateOpts{
		FloatingNetworkID: "a6917946-38ab-4ffd-a55a-26c0980ce5ee",
	}

	fip, err := floatingips.Create(networkingClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Floating IP

	fipID := "2f245a7b-796b-4f26-9cf9-9e82d248fda7"
	portID := "76d0a61b-b8e5-490c-9892-4cf674f2bec8"

	updateOpts := floatingips.UpdateOpts{
		PortID: &portID,
	}

	fip, err := floatingips.Update(networkingClient, fipID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Disassociate a Floating IP with a Port

	fipID := "2f245a7b-796b-4f26-9cf9-9e82d248fda7"

	updateOpts := floatingips.UpdateOpts{
		PortID: nil,
	}

	fip, err := floatingips.Update(networkingClient, fipID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Floating IP

	fipID := "2f245a7b-796b-4f26-9cf9-9e82d248fda7"
	err := floatingips.Delete(networkClient, fipID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package floatingips
