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
Package groups provides information and interaction with Security Groups
for the OpenStack Networking service.

Example to List Security Groups

	listOpts := groups.ListOpts{
		TenantID: "966b3c7d36a24facaf20b7e458bf2192",
	}

	allPages, err := groups.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		panic(err)
	}

	for _, group := range allGroups {
		fmt.Printf("%+v\n", group)
	}

Example to Create a Security Group

	createOpts := groups.CreateOpts{
		Name:        "group_name",
		Description: "A Security Group",
	}

	group, err := groups.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Security Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"

	updateOpts := groups.UpdateOpts{
		Name: "new_name",
	}

	group, err := groups.Update(networkClient, groupID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Security Group

	groupID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	err := groups.Delete(networkClient, groupID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package groups
