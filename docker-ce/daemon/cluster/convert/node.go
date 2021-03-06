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
	"fmt"
	"strings"

	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

// NodeFromGRPC converts a grpc Node to a Node.
func NodeFromGRPC(n swarmapi.Node) types.Node {
	node := types.Node{
		ID: n.ID,
		Spec: types.NodeSpec{
			Role:         types.NodeRole(strings.ToLower(n.Spec.DesiredRole.String())),
			Availability: types.NodeAvailability(strings.ToLower(n.Spec.Availability.String())),
		},
		Status: types.NodeStatus{
			State:   types.NodeState(strings.ToLower(n.Status.State.String())),
			Message: n.Status.Message,
			Addr:    n.Status.Addr,
		},
	}

	// Meta
	node.Version.Index = n.Meta.Version.Index
	node.CreatedAt, _ = gogotypes.TimestampFromProto(n.Meta.CreatedAt)
	node.UpdatedAt, _ = gogotypes.TimestampFromProto(n.Meta.UpdatedAt)

	//Annotations
	node.Spec.Annotations = annotationsFromGRPC(n.Spec.Annotations)

	//Description
	if n.Description != nil {
		node.Description.Hostname = n.Description.Hostname
		if n.Description.Platform != nil {
			node.Description.Platform.Architecture = n.Description.Platform.Architecture
			node.Description.Platform.OS = n.Description.Platform.OS
		}
		if n.Description.Resources != nil {
			node.Description.Resources.NanoCPUs = n.Description.Resources.NanoCPUs
			node.Description.Resources.MemoryBytes = n.Description.Resources.MemoryBytes
		}
		if n.Description.Engine != nil {
			node.Description.Engine.EngineVersion = n.Description.Engine.EngineVersion
			node.Description.Engine.Labels = n.Description.Engine.Labels
			for _, plugin := range n.Description.Engine.Plugins {
				node.Description.Engine.Plugins = append(node.Description.Engine.Plugins, types.PluginDescription{Type: plugin.Type, Name: plugin.Name})
			}
		}
	}

	//Manager
	if n.ManagerStatus != nil {
		node.ManagerStatus = &types.ManagerStatus{
			Leader:       n.ManagerStatus.Leader,
			Reachability: types.Reachability(strings.ToLower(n.ManagerStatus.Reachability.String())),
			Addr:         n.ManagerStatus.Addr,
		}
	}

	return node
}

// NodeSpecToGRPC converts a NodeSpec to a grpc NodeSpec.
func NodeSpecToGRPC(s types.NodeSpec) (swarmapi.NodeSpec, error) {
	spec := swarmapi.NodeSpec{
		Annotations: swarmapi.Annotations{
			Name:   s.Name,
			Labels: s.Labels,
		},
	}
	if role, ok := swarmapi.NodeRole_value[strings.ToUpper(string(s.Role))]; ok {
		spec.DesiredRole = swarmapi.NodeRole(role)
	} else {
		return swarmapi.NodeSpec{}, fmt.Errorf("invalid Role: %q", s.Role)
	}

	if availability, ok := swarmapi.NodeSpec_Availability_value[strings.ToUpper(string(s.Availability))]; ok {
		spec.Availability = swarmapi.NodeSpec_Availability(availability)
	} else {
		return swarmapi.NodeSpec{}, fmt.Errorf("invalid Availability: %q", s.Availability)
	}

	return spec, nil
}
