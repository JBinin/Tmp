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
/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package discovery

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	apimachineryversion "k8s.io/apimachinery/pkg/version"
)

// MatchesServerVersion queries the server to compares the build version
// (git hash) of the client with the server's build version. It returns an error
// if it failed to contact the server or if the versions are not an exact match.
func MatchesServerVersion(clientVersion apimachineryversion.Info, client DiscoveryInterface) error {
	sVer, err := client.ServerVersion()
	if err != nil {
		return fmt.Errorf("couldn't read version from server: %v\n", err)
	}
	// GitVersion includes GitCommit and GitTreeState, but best to be safe?
	if clientVersion.GitVersion != sVer.GitVersion || clientVersion.GitCommit != sVer.GitCommit || clientVersion.GitTreeState != sVer.GitTreeState {
		return fmt.Errorf("server version (%#v) differs from client version (%#v)!\n", sVer, clientVersion)
	}

	return nil
}

// ServerSupportsVersion returns an error if the server doesn't have the required version
func ServerSupportsVersion(client DiscoveryInterface, requiredGV schema.GroupVersion) error {
	groups, err := client.ServerGroups()
	if err != nil {
		// This is almost always a connection error, and higher level code should treat this as a generic error,
		// not a negotiation specific error.
		return err
	}
	versions := metav1.ExtractGroupVersions(groups)
	serverVersions := sets.String{}
	for _, v := range versions {
		serverVersions.Insert(v)
	}

	if serverVersions.Has(requiredGV.String()) {
		return nil
	}

	// If the server supports no versions, then we should pretend it has the version because of old servers.
	// This can happen because discovery fails due to 403 Forbidden errors
	if len(serverVersions) == 0 {
		return nil
	}

	return fmt.Errorf("server does not support API version %q", requiredGV)
}

// GroupVersionResources converts APIResourceLists to the GroupVersionResources.
func GroupVersionResources(rls []*metav1.APIResourceList) (map[schema.GroupVersionResource]struct{}, error) {
	gvrs := map[schema.GroupVersionResource]struct{}{}
	for _, rl := range rls {
		gv, err := schema.ParseGroupVersion(rl.GroupVersion)
		if err != nil {
			return nil, err
		}
		for i := range rl.APIResources {
			gvrs[schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: rl.APIResources[i].Name}] = struct{}{}
		}
	}
	return gvrs, nil
}

// FilteredBy filters by the given predicate. Empty APIResourceLists are dropped.
func FilteredBy(pred ResourcePredicate, rls []*metav1.APIResourceList) []*metav1.APIResourceList {
	result := []*metav1.APIResourceList{}
	for _, rl := range rls {
		filtered := *rl
		filtered.APIResources = nil
		for i := range rl.APIResources {
			if pred.Match(rl.GroupVersion, &rl.APIResources[i]) {
				filtered.APIResources = append(filtered.APIResources, rl.APIResources[i])
			}
		}
		if filtered.APIResources != nil {
			result = append(result, &filtered)
		}
	}
	return result
}

type ResourcePredicate interface {
	Match(groupVersion string, r *metav1.APIResource) bool
}

type ResourcePredicateFunc func(groupVersion string, r *metav1.APIResource) bool

func (fn ResourcePredicateFunc) Match(groupVersion string, r *metav1.APIResource) bool {
	return fn(groupVersion, r)
}

// SupportsAllVerbs is a predicate matching a resource iff all given verbs are supported.
type SupportsAllVerbs struct {
	Verbs []string
}

func (p SupportsAllVerbs) Match(groupVersion string, r *metav1.APIResource) bool {
	return sets.NewString([]string(r.Verbs)...).HasAll(p.Verbs...)
}
