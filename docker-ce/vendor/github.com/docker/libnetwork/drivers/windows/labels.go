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
package windows

const (
	// NetworkName label for bridge driver
	NetworkName = "com.docker.network.windowsshim.networkname"

	// HNSID of the discovered network
	HNSID = "com.docker.network.windowsshim.hnsid"

	// RoutingDomain of the network
	RoutingDomain = "com.docker.network.windowsshim.routingdomain"

	// Interface of the network
	Interface = "com.docker.network.windowsshim.interface"

	// QosPolicies of the endpoint
	QosPolicies = "com.docker.endpoint.windowsshim.qospolicies"

	// VLAN of the network
	VLAN = "com.docker.network.windowsshim.vlanid"

	// VSID of the network
	VSID = "com.docker.network.windowsshim.vsid"

	// DNSSuffix of the network
	DNSSuffix = "com.docker.network.windowsshim.dnssuffix"

	// DNSServers of the network
	DNSServers = "com.docker.network.windowsshim.dnsservers"

	// SourceMac of the network
	SourceMac = "com.docker.network.windowsshim.sourcemac"

	// DisableICC label
	DisableICC = "com.docker.network.windowsshim.disableicc"

	// DisableDNS label
	DisableDNS = "com.docker.network.windowsshim.disable_dns"
)
