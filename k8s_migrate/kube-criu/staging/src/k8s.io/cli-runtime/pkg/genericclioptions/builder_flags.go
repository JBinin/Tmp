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
Copyright 2018 The Kubernetes Authors.

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

package genericclioptions

import (
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
)

// ResourceBuilderFlags are flags for finding resources
// TODO(juanvallejo): wire --local flag from commands through
type ResourceBuilderFlags struct {
	FileNameFlags *FileNameFlags

	LabelSelector        *string
	FieldSelector        *string
	AllNamespaces        *bool
	All                  *bool
	Local                *bool
	IncludeUninitialized *bool

	Scheme           *runtime.Scheme
	Latest           bool
	StopOnFirstError bool
}

// NewResourceBuilderFlags returns a default ResourceBuilderFlags
func NewResourceBuilderFlags() *ResourceBuilderFlags {
	filenames := []string{}

	return &ResourceBuilderFlags{
		FileNameFlags: &FileNameFlags{
			Usage:     "identifying the resource.",
			Filenames: &filenames,
			Recursive: boolPtr(true),
		},
	}
}

func (o *ResourceBuilderFlags) WithFile(recurse bool, files ...string) *ResourceBuilderFlags {
	o.FileNameFlags = &FileNameFlags{
		Usage:     "identifying the resource.",
		Filenames: &files,
		Recursive: boolPtr(recurse),
	}

	return o
}

func (o *ResourceBuilderFlags) WithLabelSelector(selector string) *ResourceBuilderFlags {
	o.LabelSelector = &selector
	return o
}

func (o *ResourceBuilderFlags) WithFieldSelector(selector string) *ResourceBuilderFlags {
	o.FieldSelector = &selector
	return o
}

func (o *ResourceBuilderFlags) WithAllNamespaces(defaultVal bool) *ResourceBuilderFlags {
	o.AllNamespaces = &defaultVal
	return o
}

func (o *ResourceBuilderFlags) WithAll(defaultVal bool) *ResourceBuilderFlags {
	o.All = &defaultVal
	return o
}

func (o *ResourceBuilderFlags) WithLocal(defaultVal bool) *ResourceBuilderFlags {
	o.Local = &defaultVal
	return o
}

// WithUninitialized is using an alpha feature and may be dropped
func (o *ResourceBuilderFlags) WithUninitialized(defaultVal bool) *ResourceBuilderFlags {
	o.IncludeUninitialized = &defaultVal
	return o
}

func (o *ResourceBuilderFlags) WithScheme(scheme *runtime.Scheme) *ResourceBuilderFlags {
	o.Scheme = scheme
	return o
}

func (o *ResourceBuilderFlags) WithLatest() *ResourceBuilderFlags {
	o.Latest = true
	return o
}

func (o *ResourceBuilderFlags) StopOnError() *ResourceBuilderFlags {
	o.StopOnFirstError = true
	return o
}

// AddFlags registers flags for finding resources
func (o *ResourceBuilderFlags) AddFlags(flagset *pflag.FlagSet) {
	o.FileNameFlags.AddFlags(flagset)

	if o.LabelSelector != nil {
		flagset.StringVarP(o.LabelSelector, "selector", "l", *o.LabelSelector, "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	}
	if o.FieldSelector != nil {
		flagset.StringVar(o.FieldSelector, "field-selector", *o.FieldSelector, "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The server only supports a limited number of field queries per type.")
	}
	if o.AllNamespaces != nil {
		flagset.BoolVar(o.AllNamespaces, "all-namespaces", *o.AllNamespaces, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	}
	if o.All != nil {
		flagset.BoolVar(o.All, "all", *o.All, "Select all resources in the namespace of the specified resource types")
	}
	if o.Local != nil {
		flagset.BoolVar(o.Local, "local", *o.Local, "If true, annotation will NOT contact api-server but run locally.")
	}
	if o.IncludeUninitialized != nil {
		flagset.BoolVar(o.IncludeUninitialized, "include-uninitialized", *o.IncludeUninitialized, `If true, the kubectl command applies to uninitialized objects. If explicitly set to false, this flag overrides other flags that make the kubectl commands apply to uninitialized objects, e.g., "--all". Objects with empty metadata.initializers are regarded as initialized.`)
	}
}

// ToBuilder gives you back a resource finder to visit resources that are located
func (o *ResourceBuilderFlags) ToBuilder(restClientGetter RESTClientGetter, resources []string) ResourceFinder {
	namespace, enforceNamespace, namespaceErr := restClientGetter.ToRawKubeConfigLoader().Namespace()

	builder := resource.NewBuilder(restClientGetter).
		NamespaceParam(namespace).DefaultNamespace()

	if o.Scheme != nil {
		builder.WithScheme(o.Scheme, o.Scheme.PrioritizedVersionsAllGroups()...)
	} else {
		builder.Unstructured()
	}

	if o.FileNameFlags != nil {
		opts := o.FileNameFlags.ToOptions()
		builder.FilenameParam(enforceNamespace, &opts)
	}

	if o.Local == nil || !*o.Local {
		// resource type/name tuples only work non-local
		if o.All != nil {
			builder.ResourceTypeOrNameArgs(*o.All, resources...)
		} else {
			builder.ResourceTypeOrNameArgs(false, resources...)
		}
		// label selectors only work non-local (for now)
		if o.LabelSelector != nil {
			builder.LabelSelectorParam(*o.LabelSelector)
		}
		// field selectors only work non-local (forever)
		if o.FieldSelector != nil {
			builder.FieldSelectorParam(*o.FieldSelector)
		}
		// latest only works non-local (forever)
		if o.Latest {
			builder.Latest()
		}

	} else {
		builder.Local()

		if len(resources) > 0 {
			builder.AddError(resource.LocalResourceError)
		}
	}

	if o.IncludeUninitialized != nil {
		builder.IncludeUninitialized(*o.IncludeUninitialized)
	}

	if !o.StopOnFirstError {
		builder.ContinueOnError()
	}

	return &ResourceFindBuilderWrapper{
		builder: builder.
			Flatten(). // I think we're going to recommend this everywhere
			AddError(namespaceErr),
	}
}

// ResourceFindBuilderWrapper wraps a builder in an interface
type ResourceFindBuilderWrapper struct {
	builder *resource.Builder
}

// Do finds you resources to check
func (b *ResourceFindBuilderWrapper) Do() resource.Visitor {
	return b.builder.Do()
}

// ResourceFinder allows mocking the resource builder
// TODO resource builders needs to become more interfacey
type ResourceFinder interface {
	Do() resource.Visitor
}

// ResourceFinderFunc is a handy way to make a  ResourceFinder
type ResourceFinderFunc func() resource.Visitor

// Do implements ResourceFinder
func (fn ResourceFinderFunc) Do() resource.Visitor {
	return fn()
}

// ResourceFinderForResult skins a visitor for re-use as a ResourceFinder
func ResourceFinderForResult(result resource.Visitor) ResourceFinder {
	return ResourceFinderFunc(func() resource.Visitor {
		return result
	})
}

func strPtr(val string) *string {
	return &val
}

func boolPtr(val bool) *bool {
	return &val
}
