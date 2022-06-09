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

// VolumeUpdateOptions are available parameters for updating existing volumes.
type VolumeUpdateOptions struct {

	// Volume unique ID.
	// Read Only: true
	ID string `json:"id"`

	// Volume name.
	// Read Only: true
	Name string `json:"name"`

	// Description describes the volume.
	Description string `json:"description"`

	// Size in GB.
	// Required: true
	Size int `json:"size"`

	// Namespace is the object scope, such as for teams and projects.
	Namespace string `json:"namespace"`

	// node selector (where volumes should land)
	NodeSelector string `json:"nodeSelector"`

	// Labels are user-defined key/value metadata.
	Labels map[string]string `json:"labels"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}
