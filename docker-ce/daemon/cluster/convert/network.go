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
package convert

import (
	"strings"

	basictypes "github.com/docker/docker/api/types"
	networktypes "github.com/docker/docker/api/types/network"
	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

func networkAttachementFromGRPC(na *swarmapi.NetworkAttachment) types.NetworkAttachment {
	if na != nil {
		return types.NetworkAttachment{
			Network:   networkFromGRPC(na.Network),
			Addresses: na.Addresses,
		}
	}
	return types.NetworkAttachment{}
}

func networkFromGRPC(n *swarmapi.Network) types.Network {
	if n != nil {
		network := types.Network{
			ID: n.ID,
			Spec: types.NetworkSpec{
				IPv6Enabled: n.Spec.Ipv6Enabled,
				Internal:    n.Spec.Internal,
				Attachable:  n.Spec.Attachable,
				IPAMOptions: ipamFromGRPC(n.Spec.IPAM),
			},
			IPAMOptions: ipamFromGRPC(n.IPAM),
		}

		// Meta
		network.Version.Index = n.Meta.Version.Index
		network.CreatedAt, _ = gogotypes.TimestampFromProto(n.Meta.CreatedAt)
		network.UpdatedAt, _ = gogotypes.TimestampFromProto(n.Meta.UpdatedAt)

		//Annotations
		network.Spec.Annotations = annotationsFromGRPC(n.Spec.Annotations)

		//DriverConfiguration
		if n.Spec.DriverConfig != nil {
			network.Spec.DriverConfiguration = &types.Driver{
				Name:    n.Spec.DriverConfig.Name,
				Options: n.Spec.DriverConfig.Options,
			}
		}

		//DriverState
		if n.DriverState != nil {
			network.DriverState = types.Driver{
				Name:    n.DriverState.Name,
				Options: n.DriverState.Options,
			}
		}

		return network
	}
	return types.Network{}
}

func ipamFromGRPC(i *swarmapi.IPAMOptions) *types.IPAMOptions {
	var ipam *types.IPAMOptions
	if i != nil {
		ipam = &types.IPAMOptions{}
		if i.Driver != nil {
			ipam.Driver.Name = i.Driver.Name
			ipam.Driver.Options = i.Driver.Options
		}

		for _, config := range i.Configs {
			ipam.Configs = append(ipam.Configs, types.IPAMConfig{
				Subnet:  config.Subnet,
				Range:   config.Range,
				Gateway: config.Gateway,
			})
		}
	}
	return ipam
}

func endpointSpecFromGRPC(es *swarmapi.EndpointSpec) *types.EndpointSpec {
	var endpointSpec *types.EndpointSpec
	if es != nil {
		endpointSpec = &types.EndpointSpec{}
		endpointSpec.Mode = types.ResolutionMode(strings.ToLower(es.Mode.String()))

		for _, portState := range es.Ports {
			endpointSpec.Ports = append(endpointSpec.Ports, types.PortConfig{
				Name:          portState.Name,
				Protocol:      types.PortConfigProtocol(strings.ToLower(swarmapi.PortConfig_Protocol_name[int32(portState.Protocol)])),
				PublishMode:   types.PortConfigPublishMode(strings.ToLower(swarmapi.PortConfig_PublishMode_name[int32(portState.PublishMode)])),
				TargetPort:    portState.TargetPort,
				PublishedPort: portState.PublishedPort,
			})
		}
	}
	return endpointSpec
}

func endpointFromGRPC(e *swarmapi.Endpoint) types.Endpoint {
	endpoint := types.Endpoint{}
	if e != nil {
		if espec := endpointSpecFromGRPC(e.Spec); espec != nil {
			endpoint.Spec = *espec
		}

		for _, portState := range e.Ports {
			endpoint.Ports = append(endpoint.Ports, types.PortConfig{
				Name:          portState.Name,
				Protocol:      types.PortConfigProtocol(strings.ToLower(swarmapi.PortConfig_Protocol_name[int32(portState.Protocol)])),
				PublishMode:   types.PortConfigPublishMode(strings.ToLower(swarmapi.PortConfig_PublishMode_name[int32(portState.PublishMode)])),
				TargetPort:    portState.TargetPort,
				PublishedPort: portState.PublishedPort,
			})
		}

		for _, v := range e.VirtualIPs {
			endpoint.VirtualIPs = append(endpoint.VirtualIPs, types.EndpointVirtualIP{
				NetworkID: v.NetworkID,
				Addr:      v.Addr})
		}

	}

	return endpoint
}

// BasicNetworkFromGRPC converts a grpc Network to a NetworkResource.
func BasicNetworkFromGRPC(n swarmapi.Network) basictypes.NetworkResource {
	spec := n.Spec
	var ipam networktypes.IPAM
	if spec.IPAM != nil {
		if spec.IPAM.Driver != nil {
			ipam.Driver = spec.IPAM.Driver.Name
			ipam.Options = spec.IPAM.Driver.Options
		}
		ipam.Config = make([]networktypes.IPAMConfig, 0, len(spec.IPAM.Configs))
		for _, ic := range spec.IPAM.Configs {
			ipamConfig := networktypes.IPAMConfig{
				Subnet:     ic.Subnet,
				IPRange:    ic.Range,
				Gateway:    ic.Gateway,
				AuxAddress: ic.Reserved,
			}
			ipam.Config = append(ipam.Config, ipamConfig)
		}
	}

	nr := basictypes.NetworkResource{
		ID:         n.ID,
		Name:       n.Spec.Annotations.Name,
		Scope:      "swarm",
		EnableIPv6: spec.Ipv6Enabled,
		IPAM:       ipam,
		Internal:   spec.Internal,
		Attachable: spec.Attachable,
		Labels:     n.Spec.Annotations.Labels,
	}

	if n.DriverState != nil {
		nr.Driver = n.DriverState.Name
		nr.Options = n.DriverState.Options
	}

	return nr
}

// BasicNetworkCreateToGRPC converts a NetworkCreateRequest to a grpc NetworkSpec.
func BasicNetworkCreateToGRPC(create basictypes.NetworkCreateRequest) swarmapi.NetworkSpec {
	ns := swarmapi.NetworkSpec{
		Annotations: swarmapi.Annotations{
			Name:   create.Name,
			Labels: create.Labels,
		},
		DriverConfig: &swarmapi.Driver{
			Name:    create.Driver,
			Options: create.Options,
		},
		Ipv6Enabled: create.EnableIPv6,
		Internal:    create.Internal,
		Attachable:  create.Attachable,
	}
	if create.IPAM != nil {
		driver := create.IPAM.Driver
		if driver == "" {
			driver = "default"
		}
		ns.IPAM = &swarmapi.IPAMOptions{
			Driver: &swarmapi.Driver{
				Name:    driver,
				Options: create.IPAM.Options,
			},
		}
		ipamSpec := make([]*swarmapi.IPAMConfig, 0, len(create.IPAM.Config))
		for _, ipamConfig := range create.IPAM.Config {
			ipamSpec = append(ipamSpec, &swarmapi.IPAMConfig{
				Subnet:  ipamConfig.Subnet,
				Range:   ipamConfig.IPRange,
				Gateway: ipamConfig.Gateway,
			})
		}
		ns.IPAM.Configs = ipamSpec
	}
	return ns
}
