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
package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// ImageSummary image summary
// swagger:model ImageSummary
type ImageSummary struct {

	// containers
	// Required: true
	Containers int64 `json:"Containers"`

	// created
	// Required: true
	Created int64 `json:"Created"`

	// Id
	// Required: true
	ID string `json:"Id"`

	// labels
	// Required: true
	Labels map[string]string `json:"Labels"`

	// parent Id
	// Required: true
	ParentID string `json:"ParentId"`

	// repo digests
	// Required: true
	RepoDigests []string `json:"RepoDigests"`

	// repo tags
	// Required: true
	RepoTags []string `json:"RepoTags"`

	// shared size
	// Required: true
	SharedSize int64 `json:"SharedSize"`

	// size
	// Required: true
	Size int64 `json:"Size"`

	// virtual size
	// Required: true
	VirtualSize int64 `json:"VirtualSize"`
}
