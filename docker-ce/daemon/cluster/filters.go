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
package cluster

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/filters"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	swarmapi "github.com/docker/swarmkit/api"
)

func newListNodesFilters(filter filters.Args) (*swarmapi.ListNodesRequest_Filters, error) {
	accepted := map[string]bool{
		"name":       true,
		"id":         true,
		"label":      true,
		"role":       true,
		"membership": true,
	}
	if err := filter.Validate(accepted); err != nil {
		return nil, err
	}
	f := &swarmapi.ListNodesRequest_Filters{
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
	}

	for _, r := range filter.Get("role") {
		if role, ok := swarmapi.NodeRole_value[strings.ToUpper(r)]; ok {
			f.Roles = append(f.Roles, swarmapi.NodeRole(role))
		} else if r != "" {
			return nil, fmt.Errorf("Invalid role filter: '%s'", r)
		}
	}

	for _, a := range filter.Get("membership") {
		if membership, ok := swarmapi.NodeSpec_Membership_value[strings.ToUpper(a)]; ok {
			f.Memberships = append(f.Memberships, swarmapi.NodeSpec_Membership(membership))
		} else if a != "" {
			return nil, fmt.Errorf("Invalid membership filter: '%s'", a)
		}
	}

	return f, nil
}

func newListServicesFilters(filter filters.Args) (*swarmapi.ListServicesRequest_Filters, error) {
	accepted := map[string]bool{
		"name":  true,
		"id":    true,
		"label": true,
	}
	if err := filter.Validate(accepted); err != nil {
		return nil, err
	}
	return &swarmapi.ListServicesRequest_Filters{
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
	}, nil
}

func newListTasksFilters(filter filters.Args, transformFunc func(filters.Args) error) (*swarmapi.ListTasksRequest_Filters, error) {
	accepted := map[string]bool{
		"name":          true,
		"id":            true,
		"label":         true,
		"service":       true,
		"node":          true,
		"desired-state": true,
	}
	if err := filter.Validate(accepted); err != nil {
		return nil, err
	}
	if transformFunc != nil {
		if err := transformFunc(filter); err != nil {
			return nil, err
		}
	}
	f := &swarmapi.ListTasksRequest_Filters{
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
		ServiceIDs:   filter.Get("service"),
		NodeIDs:      filter.Get("node"),
	}

	for _, s := range filter.Get("desired-state") {
		if state, ok := swarmapi.TaskState_value[strings.ToUpper(s)]; ok {
			f.DesiredStates = append(f.DesiredStates, swarmapi.TaskState(state))
		} else if s != "" {
			return nil, fmt.Errorf("Invalid desired-state filter: '%s'", s)
		}
	}

	return f, nil
}

func newListSecretsFilters(filter filters.Args) (*swarmapi.ListSecretsRequest_Filters, error) {
	accepted := map[string]bool{
		"names": true,
		"name":  true,
		"id":    true,
		"label": true,
	}
	if err := filter.Validate(accepted); err != nil {
		return nil, err
	}
	return &swarmapi.ListSecretsRequest_Filters{
		Names:        filter.Get("names"),
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
	}, nil
}
