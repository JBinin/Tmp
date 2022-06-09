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
Package tokens provides information and interaction with the token API
resource for the OpenStack Identity service.

For more information, see:
http://developer.openstack.org/api-ref-identity-v2.html#identity-auth-v2

Example to Create an Unscoped Token from a Password

	authOpts := gophercloud.AuthOptions{
		Username: "user",
		Password: "pass"
	}

	token, err := tokens.Create(identityClient, authOpts).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Create a Token from a Tenant ID and Password

	authOpts := gophercloud.AuthOptions{
		Username: "user",
		Password: "password",
		TenantID: "fc394f2ab2df4114bde39905f800dc57"
	}

	token, err := tokens.Create(identityClient, authOpts).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Create a Token from a Tenant Name and Password

	authOpts := gophercloud.AuthOptions{
		Username:   "user",
		Password:   "password",
		TenantName: "tenantname"
	}

	token, err := tokens.Create(identityClient, authOpts).ExtractToken()
	if err != nil {
		panic(err)
	}
*/
package tokens
