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

// ListOptions are optional parameters for finding and listing most objects.
type ListOptions struct {

	// FieldSelector restricts the list of returned objects by their fields. Defaults to everything.
	FieldSelector string

	// LabelSelector restricts the list of returned objects by their labels. Defaults to everything.
	LabelSelector string

	// Namespace is the object scope, such as for teams and projects.
	Namespace string

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context
}
