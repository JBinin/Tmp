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
Package loadbalancers provides information and interaction with Load Balancers
of the LBaaS v2 extension for the OpenStack Networking service.

Example to List Load Balancers

	listOpts := loadbalancers.ListOpts{
		Provider: "haproxy",
	}

	allPages, err := loadbalancers.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		panic(err)
	}

	for _, lb := range allLoadbalancers {
		fmt.Printf("%+v\n", lb)
	}

Example to Create a Load Balancer

	createOpts := loadbalancers.CreateOpts{
		Name:         "db_lb",
		AdminStateUp: gophercloud.Enabled,
		VipSubnetID:  "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		VipAddress:   "10.30.176.48",
		Flavor:       "medium",
		Provider:     "haproxy",
	}

	lb, err := loadbalancers.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Load Balancer

	lbID := "d67d56a6-4a86-4688-a282-f46444705c64"

	i1001 := 1001
	updateOpts := loadbalancers.UpdateOpts{
		Name: "new-name",
	}

	lb, err := loadbalancers.Update(networkClient, lbID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Load Balancers

	lbID := "d67d56a6-4a86-4688-a282-f46444705c64"
	err := loadbalancers.Delete(networkClient, lbID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Get the Status of a Load Balancer

	lbID := "d67d56a6-4a86-4688-a282-f46444705c64"
	status, err := loadbalancers.GetStatuses(networkClient, LBID).Extract()
	if err != nil {
		panic(err)
	}
*/
package loadbalancers
