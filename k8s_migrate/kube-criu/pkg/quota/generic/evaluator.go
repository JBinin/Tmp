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

package generic

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/quota"
)

// InformerForResourceFunc knows how to provision an informer
type InformerForResourceFunc func(schema.GroupVersionResource) (informers.GenericInformer, error)

// ListerFuncForResourceFunc knows how to provision a lister from an informer func
func ListerFuncForResourceFunc(f InformerForResourceFunc) quota.ListerForResourceFunc {
	return func(gvr schema.GroupVersionResource) (cache.GenericLister, error) {
		informer, err := f(gvr)
		if err != nil {
			return nil, err
		}
		return informer.Lister(), nil
	}
}

// ListResourceUsingListerFunc returns a listing function based on the shared informer factory for the specified resource.
func ListResourceUsingListerFunc(l quota.ListerForResourceFunc, resource schema.GroupVersionResource) ListFuncByNamespace {
	return func(namespace string) ([]runtime.Object, error) {
		lister, err := l(resource)
		if err != nil {
			return nil, err
		}
		return lister.ByNamespace(namespace).List(labels.Everything())
	}
}

// ObjectCountQuotaResourceNameFor returns the object count quota name for specified groupResource
func ObjectCountQuotaResourceNameFor(groupResource schema.GroupResource) api.ResourceName {
	if len(groupResource.Group) == 0 {
		return api.ResourceName("count/" + groupResource.Resource)
	}
	return api.ResourceName("count/" + groupResource.Resource + "." + groupResource.Group)
}

// ListFuncByNamespace knows how to list resources in a namespace
type ListFuncByNamespace func(namespace string) ([]runtime.Object, error)

// MatchesScopeFunc knows how to evaluate if an object matches a scope
type MatchesScopeFunc func(scope api.ScopedResourceSelectorRequirement, object runtime.Object) (bool, error)

// UsageFunc knows how to measure usage associated with an object
type UsageFunc func(object runtime.Object) (api.ResourceList, error)

// MatchingResourceNamesFunc is a function that returns the list of resources matched
type MatchingResourceNamesFunc func(input []api.ResourceName) []api.ResourceName

// MatchesNoScopeFunc returns false on all match checks
func MatchesNoScopeFunc(scope api.ScopedResourceSelectorRequirement, object runtime.Object) (bool, error) {
	return false, nil
}

// Matches returns true if the quota matches the specified item.
func Matches(
	resourceQuota *api.ResourceQuota, item runtime.Object,
	matchFunc MatchingResourceNamesFunc, scopeFunc MatchesScopeFunc) (bool, error) {
	if resourceQuota == nil {
		return false, fmt.Errorf("expected non-nil quota")
	}
	// verify the quota matches on at least one resource
	matchResource := len(matchFunc(quota.ResourceNames(resourceQuota.Status.Hard))) > 0
	// by default, no scopes matches all
	matchScope := true
	for _, scope := range getScopeSelectorsFromQuota(resourceQuota) {
		innerMatch, err := scopeFunc(scope, item)
		if err != nil {
			return false, err
		}
		matchScope = matchScope && innerMatch
	}
	return matchResource && matchScope, nil
}

func getScopeSelectorsFromQuota(quota *api.ResourceQuota) []api.ScopedResourceSelectorRequirement {
	selectors := []api.ScopedResourceSelectorRequirement{}
	for _, scope := range quota.Spec.Scopes {
		selectors = append(selectors, api.ScopedResourceSelectorRequirement{
			ScopeName: scope,
			Operator:  api.ScopeSelectorOpExists})
	}
	if quota.Spec.ScopeSelector != nil {
		for _, scopeSelector := range quota.Spec.ScopeSelector.MatchExpressions {
			selectors = append(selectors, scopeSelector)
		}
	}
	return selectors
}

// CalculateUsageStats is a utility function that knows how to calculate aggregate usage.
func CalculateUsageStats(options quota.UsageStatsOptions,
	listFunc ListFuncByNamespace,
	scopeFunc MatchesScopeFunc,
	usageFunc UsageFunc) (quota.UsageStats, error) {
	// default each tracked resource to zero
	result := quota.UsageStats{Used: api.ResourceList{}}
	for _, resourceName := range options.Resources {
		result.Used[resourceName] = resource.Quantity{Format: resource.DecimalSI}
	}
	items, err := listFunc(options.Namespace)
	if err != nil {
		return result, fmt.Errorf("failed to list content: %v", err)
	}
	for _, item := range items {
		// need to verify that the item matches the set of scopes
		matchesScopes := true
		for _, scope := range options.Scopes {
			innerMatch, err := scopeFunc(api.ScopedResourceSelectorRequirement{ScopeName: scope}, item)
			if err != nil {
				return result, nil
			}
			if !innerMatch {
				matchesScopes = false
			}
		}
		if options.ScopeSelector != nil {
			for _, selector := range options.ScopeSelector.MatchExpressions {
				innerMatch, err := scopeFunc(selector, item)
				if err != nil {
					return result, nil
				}
				matchesScopes = matchesScopes && innerMatch
			}
		}
		// only count usage if there was a match
		if matchesScopes {
			usage, err := usageFunc(item)
			if err != nil {
				return result, err
			}
			result.Used = quota.Add(result.Used, usage)
		}
	}
	return result, nil
}

// objectCountEvaluator provides an implementation for quota.Evaluator
// that associates usage of the specified resource based on the number of items
// returned by the specified listing function.
type objectCountEvaluator struct {
	// GroupResource that this evaluator tracks.
	// It is used to construct a generic object count quota name
	groupResource schema.GroupResource
	// A function that knows how to list resources by namespace.
	// TODO move to dynamic client in future
	listFuncByNamespace ListFuncByNamespace
	// Names associated with this resource in the quota for generic counting.
	resourceNames []api.ResourceName
}

// Constraints returns an error if the configured resource name is not in the required set.
func (o *objectCountEvaluator) Constraints(required []api.ResourceName, item runtime.Object) error {
	// no-op for object counting
	return nil
}

// Handles returns true if the object count evaluator needs to track this attributes.
func (o *objectCountEvaluator) Handles(a admission.Attributes) bool {
	operation := a.GetOperation()
	return operation == admission.Create
}

// Matches returns true if the evaluator matches the specified quota with the provided input item
func (o *objectCountEvaluator) Matches(resourceQuota *api.ResourceQuota, item runtime.Object) (bool, error) {
	return Matches(resourceQuota, item, o.MatchingResources, MatchesNoScopeFunc)
}

// MatchingResources takes the input specified list of resources and returns the set of resources it matches.
func (o *objectCountEvaluator) MatchingResources(input []api.ResourceName) []api.ResourceName {
	return quota.Intersection(input, o.resourceNames)
}

// MatchingScopes takes the input specified list of scopes and input object. Returns the set of scopes resource matches.
func (o *objectCountEvaluator) MatchingScopes(item runtime.Object, scopes []api.ScopedResourceSelectorRequirement) ([]api.ScopedResourceSelectorRequirement, error) {
	return []api.ScopedResourceSelectorRequirement{}, nil
}

// UncoveredQuotaScopes takes the input matched scopes which are limited by configuration and the matched quota scopes.
// It returns the scopes which are in limited scopes but dont have a corresponding covering quota scope
func (o *objectCountEvaluator) UncoveredQuotaScopes(limitedScopes []api.ScopedResourceSelectorRequirement, matchedQuotaScopes []api.ScopedResourceSelectorRequirement) ([]api.ScopedResourceSelectorRequirement, error) {
	return []api.ScopedResourceSelectorRequirement{}, nil
}

// Usage returns the resource usage for the specified object
func (o *objectCountEvaluator) Usage(object runtime.Object) (api.ResourceList, error) {
	quantity := resource.NewQuantity(1, resource.DecimalSI)
	resourceList := api.ResourceList{}
	for _, resourceName := range o.resourceNames {
		resourceList[resourceName] = *quantity
	}
	return resourceList, nil
}

// GroupResource tracked by this evaluator
func (o *objectCountEvaluator) GroupResource() schema.GroupResource {
	return o.groupResource
}

// UsageStats calculates aggregate usage for the object.
func (o *objectCountEvaluator) UsageStats(options quota.UsageStatsOptions) (quota.UsageStats, error) {
	return CalculateUsageStats(options, o.listFuncByNamespace, MatchesNoScopeFunc, o.Usage)
}

// Verify implementation of interface at compile time.
var _ quota.Evaluator = &objectCountEvaluator{}

// NewObjectCountEvaluator returns an evaluator that can perform generic
// object quota counting.  It allows an optional alias for backwards compatibility
// purposes for the legacy object counting names in quota.  Unless its supporting
// backward compatibility, alias should not be used.
func NewObjectCountEvaluator(
	groupResource schema.GroupResource, listFuncByNamespace ListFuncByNamespace,
	alias api.ResourceName) quota.Evaluator {

	resourceNames := []api.ResourceName{ObjectCountQuotaResourceNameFor(groupResource)}
	if len(alias) > 0 {
		resourceNames = append(resourceNames, alias)
	}

	return &objectCountEvaluator{
		groupResource:       groupResource,
		listFuncByNamespace: listFuncByNamespace,
		resourceNames:       resourceNames,
	}
}
