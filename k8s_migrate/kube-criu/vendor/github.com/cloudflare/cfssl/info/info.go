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
// Package info contains the definitions for the info endpoint
package info

// Req is the request struct for an info API request.
type Req struct {
	Label   string `json:"label"`
	Profile string `json:"profile"`
}

// Resp is the response for an Info API request.
type Resp struct {
	Certificate  string   `json:"certificate"`
	Usage        []string `json:"usages"`
	ExpiryString string   `json:"expiry"`
}
