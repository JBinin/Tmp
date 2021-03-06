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
package controlapi

import (
	"regexp"
	"strings"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/allocator/networkallocator"
	"github.com/docker/swarmkit/manager/state/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var isValidName = regexp.MustCompile(`^[a-zA-Z0-9](?:[-_]*[A-Za-z0-9]+)*$`)

func buildFilters(by func(string) store.By, values []string) store.By {
	filters := make([]store.By, 0, len(values))
	for _, v := range values {
		filters = append(filters, by(v))
	}
	return store.Or(filters...)
}

func filterContains(match string, candidates []string) bool {
	if len(candidates) == 0 {
		return true
	}
	for _, c := range candidates {
		if c == match {
			return true
		}
	}
	return false
}

func filterContainsPrefix(match string, candidates []string) bool {
	if len(candidates) == 0 {
		return true
	}
	for _, c := range candidates {
		if strings.HasPrefix(match, c) {
			return true
		}
	}
	return false
}

func filterMatchLabels(match map[string]string, candidates map[string]string) bool {
	if len(candidates) == 0 {
		return true
	}

	for k, v := range candidates {
		c, ok := match[k]
		if !ok {
			return false
		}
		if v != "" && v != c {
			return false
		}
	}
	return true
}

func validateAnnotations(m api.Annotations) error {
	if m.Name == "" {
		return grpc.Errorf(codes.InvalidArgument, "meta: name must be provided")
	}
	if !isValidName.MatchString(m.Name) {
		// if the name doesn't match the regex
		return grpc.Errorf(codes.InvalidArgument, "name must be valid as a DNS name component")
	}
	if len(m.Name) > 63 {
		// DNS labels are limited to 63 characters
		return grpc.Errorf(codes.InvalidArgument, "name must be 63 characters or fewer")
	}
	return nil
}

func validateDriver(driver *api.Driver, pg plugingetter.PluginGetter, pluginType string) error {
	if driver == nil {
		// It is ok to not specify the driver. We will choose
		// a default driver.
		return nil
	}

	if driver.Name == "" {
		return grpc.Errorf(codes.InvalidArgument, "driver name: if driver is specified name is required")
	}

	if strings.ToLower(driver.Name) == networkallocator.DefaultDriver || strings.ToLower(driver.Name) == ipamapi.DefaultIPAM {
		return nil
	}

	p, err := pg.Get(driver.Name, pluginType, plugingetter.Lookup)
	if err != nil {
		return grpc.Errorf(codes.InvalidArgument, "error during lookup of plugin %s", driver.Name)
	}

	if p.IsV1() {
		return grpc.Errorf(codes.InvalidArgument, "legacy plugin %s of type %s is not supported in swarm mode", driver.Name, pluginType)
	}

	return nil
}
