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
package swarm

import "time"

// Version represents the internal object version.
type Version struct {
	Index uint64 `json:",omitempty"`
}

// Meta is a base object inherited by most of the other once.
type Meta struct {
	Version   Version   `json:",omitempty"`
	CreatedAt time.Time `json:",omitempty"`
	UpdatedAt time.Time `json:",omitempty"`
}

// Annotations represents how to describe an object.
type Annotations struct {
	Name   string            `json:",omitempty"`
	Labels map[string]string `json:"Labels"`
}

// Driver represents a driver (network, logging).
type Driver struct {
	Name    string            `json:",omitempty"`
	Options map[string]string `json:",omitempty"`
}
