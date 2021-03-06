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
package network

// Address represents an IP address
type Address struct {
	Addr      string
	PrefixLen int
}

// IPAM represents IP Address Management
type IPAM struct {
	Driver  string
	Options map[string]string //Per network IPAM driver options
	Config  []IPAMConfig
}

// IPAMConfig represents IPAM configurations
type IPAMConfig struct {
	Subnet     string            `json:",omitempty"`
	IPRange    string            `json:",omitempty"`
	Gateway    string            `json:",omitempty"`
	AuxAddress map[string]string `json:"AuxiliaryAddresses,omitempty"`
}

// EndpointIPAMConfig represents IPAM configurations for the endpoint
type EndpointIPAMConfig struct {
	IPv4Address  string   `json:",omitempty"`
	IPv6Address  string   `json:",omitempty"`
	LinkLocalIPs []string `json:",omitempty"`
}

// Copy makes a copy of the endpoint ipam config
func (cfg *EndpointIPAMConfig) Copy() *EndpointIPAMConfig {
	cfgCopy := *cfg
	cfgCopy.LinkLocalIPs = make([]string, 0, len(cfg.LinkLocalIPs))
	cfgCopy.LinkLocalIPs = append(cfgCopy.LinkLocalIPs, cfg.LinkLocalIPs...)
	return &cfgCopy
}

// PeerInfo represents one peer of an overlay network
type PeerInfo struct {
	Name string
	IP   string
}

// EndpointSettings stores the network endpoint details
type EndpointSettings struct {
	// Configurations
	IPAMConfig *EndpointIPAMConfig
	Links      []string
	Aliases    []string
	// Operational data
	NetworkID           string
	EndpointID          string
	Gateway             string
	IPAddress           string
	IPPrefixLen         int
	IPv6Gateway         string
	GlobalIPv6Address   string
	GlobalIPv6PrefixLen int
	MacAddress          string
}

// Copy makes a deep copy of `EndpointSettings`
func (es *EndpointSettings) Copy() *EndpointSettings {
	epCopy := *es
	if es.IPAMConfig != nil {
		epCopy.IPAMConfig = es.IPAMConfig.Copy()
	}

	if es.Links != nil {
		links := make([]string, 0, len(es.Links))
		epCopy.Links = append(links, es.Links...)
	}

	if es.Aliases != nil {
		aliases := make([]string, 0, len(es.Aliases))
		epCopy.Aliases = append(aliases, es.Aliases...)
	}
	return &epCopy
}

// NetworkingConfig represents the container's networking configuration for each of its interfaces
// Carries the networking configs specified in the `docker run` and `docker network connect` commands
type NetworkingConfig struct {
	EndpointsConfig map[string]*EndpointSettings // Endpoint configs for each connecting network
}
