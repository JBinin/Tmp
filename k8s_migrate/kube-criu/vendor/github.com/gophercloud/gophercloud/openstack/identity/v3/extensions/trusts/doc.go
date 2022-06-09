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
Package trusts enables management of OpenStack Identity Trusts.

Example to Create a Token with Username, Password, and Trust ID

	var trustToken struct {
		tokens.Token
		trusts.TokenExt
	}

	authOptions := tokens.AuthOptions{
		UserID:   "username",
		Password: "password",
	}

	createOpts := trusts.AuthOptsExt{
		AuthOptionsBuilder: authOptions,
		TrustID:            "de0945a",
	}

	err := tokens.Create(identityClient, createOpts).ExtractInto(&trustToken)
	if err != nil {
		panic(err)
	}
*/
package trusts
