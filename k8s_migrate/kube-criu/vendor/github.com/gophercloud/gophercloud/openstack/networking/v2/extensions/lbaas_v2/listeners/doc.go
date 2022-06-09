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
Package listeners provides information and interaction with Listeners of the
LBaaS v2 extension for the OpenStack Networking service.

Example to List Listeners

	listOpts := listeners.ListOpts{
		LoadbalancerID : "ca430f80-1737-4712-8dc6-3f640d55594b",
	}

	allPages, err := listeners.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allListeners, err := listeners.ExtractListeners(allPages)
	if err != nil {
		panic(err)
	}

	for _, listener := range allListeners {
		fmt.Printf("%+v\n", listener)
	}

Example to Create a Listener

	createOpts := listeners.CreateOpts{
		Protocol:               "TCP",
		Name:                   "db",
		LoadbalancerID:         "79e05663-7f03-45d2-a092-8b94062f22ab",
		AdminStateUp:           gophercloud.Enabled,
		DefaultPoolID:          "41efe233-7591-43c5-9cf7-923964759f9e",
		ProtocolPort:           3306,
	}

	listener, err := listeners.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Listener

	listenerID := "d67d56a6-4a86-4688-a282-f46444705c64"

	i1001 := 1001
	updateOpts := listeners.UpdateOpts{
		ConnLimit: &i1001,
	}

	listener, err := listeners.Update(networkClient, listenerID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Listener

	listenerID := "d67d56a6-4a86-4688-a282-f46444705c64"
	err := listeners.Delete(networkClient, listenerID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package listeners
