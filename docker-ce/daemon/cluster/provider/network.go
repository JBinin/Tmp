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
package provider

import "github.com/docker/docker/api/types"

// NetworkCreateRequest is a request when creating a network.
type NetworkCreateRequest struct {
	ID string
	types.NetworkCreateRequest
}

// NetworkCreateResponse is a response when creating a network.
type NetworkCreateResponse struct {
	ID string `json:"Id"`
}

// VirtualAddress represents a virtual address.
type VirtualAddress struct {
	IPv4 string
	IPv6 string
}

// PortConfig represents a port configuration.
type PortConfig struct {
	Name          string
	Protocol      int32
	TargetPort    uint32
	PublishedPort uint32
}

// ServiceConfig represents a service configuration.
type ServiceConfig struct {
	ID               string
	Name             string
	Aliases          map[string][]string
	VirtualAddresses map[string]*VirtualAddress
	ExposedPorts     []*PortConfig
}
