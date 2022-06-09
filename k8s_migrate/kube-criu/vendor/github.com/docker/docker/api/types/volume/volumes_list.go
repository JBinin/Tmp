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
package volume

// ----------------------------------------------------------------------------
// DO NOT EDIT THIS FILE
// This file was generated by `swagger generate operation`
//
// See hack/generate-swagger-api.sh
// ----------------------------------------------------------------------------

import "github.com/docker/docker/api/types"

// VolumesListOKBody volumes list o k body
// swagger:model VolumesListOKBody
type VolumesListOKBody struct {

	// List of volumes
	// Required: true
	Volumes []*types.Volume `json:"Volumes"`

	// Warnings that occurred when fetching the list of volumes
	// Required: true
	Warnings []string `json:"Warnings"`
}
