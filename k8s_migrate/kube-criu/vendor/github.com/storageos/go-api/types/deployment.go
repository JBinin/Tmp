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

import "time"

// Deployment Volume master or replica deployment details.
// swagger:model Deployment
type Deployment struct {

	// Deployment unique ID
	// Read Only: true
	ID string `json:"id"`

	// Inode number
	// Read Only: true
	Inode uint32 `json:"inode"`

	// Controller ID
	// Read Only: true
	Controller string `json:"controller"`

	// Controller name
	// Read Only: true
	ControllerName string `json:"controllerName"`

	// Health
	// Read Only: true
	Health string `json:"health"`

	// Status
	// Read Only: true
	Status string `json:"status"`

	// Created at
	// Read Only: true
	CreatedAt time.Time `json:"createdAt"`
}
