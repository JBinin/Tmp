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

// DeleteOptions are available parameters for deleting existing volumes.
type DeleteOptions struct {

	// Volume unique ID.
	// Read Only: true
	ID string `json:"id"`

	// Volume name.
	// Read Only: true
	Name string `json:"name"`

	// Namespace is the object scope, such as for teams and projects.
	Namespace string `json:"namespace"`

	// Force will cause the volume to be deleted even if it's in use.
	Force bool `json:"force"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}
