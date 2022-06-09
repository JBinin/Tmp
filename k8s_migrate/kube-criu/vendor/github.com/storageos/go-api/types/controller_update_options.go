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

import "context"

// ControllerUpdateOptions are available parameters for updating existing controllers.
type ControllerUpdateOptions struct {

	// Controller unique ID.
	// Read Only: true
	ID string `json:"id"`

	// Controller name.
	// Read Only: true
	Name string `json:"name"`

	// Description of the controller.
	Description string `json:"description"`

	// Labels are user-defined key/value metadata.
	Labels map[string]string `json:"labels"`

	// Cordon sets the controler into an unschedulable state if true
	Cordon bool `json:"unschedulable"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}
