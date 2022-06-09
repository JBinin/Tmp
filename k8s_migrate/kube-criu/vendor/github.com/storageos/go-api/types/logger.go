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

// Logger is the runtime configuration of the node's logging services.
// swagger:model Logger
type Logger struct {

	// Node name
	Node string `json:"node"`

	// Log level
	Level string `json:"level"`

	// Log filter
	Filter string `json:"filter"`

	// Log filters by category
	// Read Only: true
	Categories map[string]string `json:"categories"`
}

// LoggerUpdateOptions are the available parameters for updating loggers.
type LoggerUpdateOptions struct {

	// Log level
	Level string `json:"level"`

	// Log filter
	Filter string `json:"filter"`

	// List of nodes to update.  All if not set.
	Nodes []string `json:"nodes"`

	// List of fields to update.  Must be set.
	Fields []string `json:"fields"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}
