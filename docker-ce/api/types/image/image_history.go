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
package image

// ----------------------------------------------------------------------------
// DO NOT EDIT THIS FILE
// This file was generated by `swagger generate operation`
//
// See hack/generate-swagger-api.sh
// ----------------------------------------------------------------------------

// HistoryResponseItem history response item
// swagger:model HistoryResponseItem
type HistoryResponseItem struct {

	// comment
	// Required: true
	Comment string `json:"Comment"`

	// created
	// Required: true
	Created int64 `json:"Created"`

	// created by
	// Required: true
	CreatedBy string `json:"CreatedBy"`

	// Id
	// Required: true
	ID string `json:"Id"`

	// size
	// Required: true
	Size int64 `json:"Size"`

	// tags
	// Required: true
	Tags []string `json:"Tags"`
}
