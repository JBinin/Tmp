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
Copyright 2015 The Kubernetes Authors.

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

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/core/event"
)

type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against events.
func NewREST(optsGetter generic.RESTOptionsGetter, ttl uint64) *REST {
	resource := api.Resource("events")
	opts, err := optsGetter.GetRESTOptions(resource)
	if err != nil {
		panic(err) // TODO: Propagate error up
	}

	// We explicitly do NOT do any decoration here - switching on Cacher
	// for events will lead to too high memory consumption.
	opts.Decorator = generic.UndecoratedStorage // TODO use watchCacheSize=-1 to signal UndecoratedStorage

	store := &genericregistry.Store{
		NewFunc:       func() runtime.Object { return &api.Event{} },
		NewListFunc:   func() runtime.Object { return &api.EventList{} },
		PredicateFunc: event.MatchEvent,
		TTLFunc: func(runtime.Object, uint64, bool) (uint64, error) {
			return ttl, nil
		},
		DefaultQualifiedResource: resource,

		CreateStrategy: event.Strategy,
		UpdateStrategy: event.Strategy,
		DeleteStrategy: event.Strategy,

		TableConvertor: printerstorage.TableConvertor{TablePrinter: printers.NewTablePrinter().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: opts, AttrFunc: event.GetAttrs} // Pass in opts to use UndecoratedStorage
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}
	return &REST{store}
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"ev"}
}
