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
// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mergo merges same-type structs and maps by setting default values in zero-value fields.

Mergo won't merge unexported (private) fields but will do recursively any exported one. It also won't merge structs inside maps (because they are not addressable using Go reflection).

Usage

From my own work-in-progress project:

	type networkConfig struct {
		Protocol string
		Address string
		ServerType string `json: "server_type"`
		Port uint16
	}

	type FssnConfig struct {
		Network networkConfig
	}

	var fssnDefault = FssnConfig {
		networkConfig {
			"tcp",
			"127.0.0.1",
			"http",
			31560,
		},
	}

	// Inside a function [...]

	if err := mergo.Merge(&config, fssnDefault); err != nil {
		log.Fatal(err)
	}

	// More code [...]

*/
package mergo
