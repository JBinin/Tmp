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
Package rules provides information and interaction with Security Group Rules
for the OpenStack Networking service.

Example to List Security Groups Rules

	listOpts := rules.ListOpts{
		Protocol: "tcp",
	}

	allPages, err := rules.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allRules, err := rules.ExtractRules(allPages)
	if err != nil {
		panic(err)
	}

	for _, rule := range allRules {
		fmt.Printf("%+v\n", rule)
	}

Example to Create a Security Group Rule

	createOpts := rules.CreateOpts{
		Direction:     "ingress",
		PortRangeMin:  80,
		EtherType:     rules.EtherType4,
		PortRangeMax:  80,
		Protocol:      "tcp",
		RemoteGroupID: "85cc3048-abc3-43cc-89b3-377341426ac5",
		SecGroupID:    "a7734e61-b545-452d-a3cd-0189cbd9747a",
	}

	rule, err := rules.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Security Group Rule

	ruleID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	err := rules.Delete(networkClient, ruleID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package rules
