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
Copyright 2014 The Kubernetes Authors.

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

package config

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/controller"
)

// ServiceHandler is an abstract interface of objects which receive
// notifications about service object changes.
type ServiceHandler interface {
	// OnServiceAdd is called whenever creation of new service object
	// is observed.
	OnServiceAdd(service *v1.Service)
	// OnServiceUpdate is called whenever modification of an existing
	// service object is observed.
	OnServiceUpdate(oldService, service *v1.Service)
	// OnServiceDelete is called whenever deletion of an existing service
	// object is observed.
	OnServiceDelete(service *v1.Service)
	// OnServiceSynced is called once all the initial even handlers were
	// called and the state is fully propagated to local cache.
	OnServiceSynced()
}

// EndpointsHandler is an abstract interface of objects which receive
// notifications about endpoints object changes.
type EndpointsHandler interface {
	// OnEndpointsAdd is called whenever creation of new endpoints object
	// is observed.
	OnEndpointsAdd(endpoints *v1.Endpoints)
	// OnEndpointsUpdate is called whenever modification of an existing
	// endpoints object is observed.
	OnEndpointsUpdate(oldEndpoints, endpoints *v1.Endpoints)
	// OnEndpointsDelete is called whever deletion of an existing endpoints
	// object is observed.
	OnEndpointsDelete(endpoints *v1.Endpoints)
	// OnEndpointsSynced is called once all the initial event handlers were
	// called and the state is fully propagated to local cache.
	OnEndpointsSynced()
}

// EndpointsConfig tracks a set of endpoints configurations.
// It accepts "set", "add" and "remove" operations of endpoints via channels, and invokes registered handlers on change.
type EndpointsConfig struct {
	lister        listers.EndpointsLister
	listerSynced  cache.InformerSynced
	eventHandlers []EndpointsHandler
}

// NewEndpointsConfig creates a new EndpointsConfig.
func NewEndpointsConfig(endpointsInformer coreinformers.EndpointsInformer, resyncPeriod time.Duration) *EndpointsConfig {
	result := &EndpointsConfig{
		lister:       endpointsInformer.Lister(),
		listerSynced: endpointsInformer.Informer().HasSynced,
	}

	endpointsInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    result.handleAddEndpoints,
			UpdateFunc: result.handleUpdateEndpoints,
			DeleteFunc: result.handleDeleteEndpoints,
		},
		resyncPeriod,
	)

	return result
}

// RegisterEventHandler registers a handler which is called on every endpoints change.
func (c *EndpointsConfig) RegisterEventHandler(handler EndpointsHandler) {
	c.eventHandlers = append(c.eventHandlers, handler)
}

// Run starts the goroutine responsible for calling registered handlers.
func (c *EndpointsConfig) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	glog.Info("Starting endpoints config controller")
	defer glog.Info("Shutting down endpoints config controller")

	if !controller.WaitForCacheSync("endpoints config", stopCh, c.listerSynced) {
		return
	}

	for i := range c.eventHandlers {
		glog.V(3).Infof("Calling handler.OnEndpointsSynced()")
		c.eventHandlers[i].OnEndpointsSynced()
	}

	<-stopCh
}

func (c *EndpointsConfig) handleAddEndpoints(obj interface{}) {
	endpoints, ok := obj.(*v1.Endpoints)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
		return
	}
	for i := range c.eventHandlers {
		glog.V(4).Infof("Calling handler.OnEndpointsAdd")
		c.eventHandlers[i].OnEndpointsAdd(endpoints)
	}
}

func (c *EndpointsConfig) handleUpdateEndpoints(oldObj, newObj interface{}) {
	oldEndpoints, ok := oldObj.(*v1.Endpoints)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", oldObj))
		return
	}
	endpoints, ok := newObj.(*v1.Endpoints)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", newObj))
		return
	}
	for i := range c.eventHandlers {
		glog.V(4).Infof("Calling handler.OnEndpointsUpdate")
		c.eventHandlers[i].OnEndpointsUpdate(oldEndpoints, endpoints)
	}
}

func (c *EndpointsConfig) handleDeleteEndpoints(obj interface{}) {
	endpoints, ok := obj.(*v1.Endpoints)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}
		if endpoints, ok = tombstone.Obj.(*v1.Endpoints); !ok {
			utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}
	}
	for i := range c.eventHandlers {
		glog.V(4).Infof("Calling handler.OnEndpointsDelete")
		c.eventHandlers[i].OnEndpointsDelete(endpoints)
	}
}

// ServiceConfig tracks a set of service configurations.
// It accepts "set", "add" and "remove" operations of services via channels, and invokes registered handlers on change.
type ServiceConfig struct {
	lister        listers.ServiceLister
	listerSynced  cache.InformerSynced
	eventHandlers []ServiceHandler
}

// NewServiceConfig creates a new ServiceConfig.
func NewServiceConfig(serviceInformer coreinformers.ServiceInformer, resyncPeriod time.Duration) *ServiceConfig {
	result := &ServiceConfig{
		lister:       serviceInformer.Lister(),
		listerSynced: serviceInformer.Informer().HasSynced,
	}

	serviceInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    result.handleAddService,
			UpdateFunc: result.handleUpdateService,
			DeleteFunc: result.handleDeleteService,
		},
		resyncPeriod,
	)

	return result
}

// RegisterEventHandler registers a handler which is called on every service change.
func (c *ServiceConfig) RegisterEventHandler(handler ServiceHandler) {
	c.eventHandlers = append(c.eventHandlers, handler)
}

// Run starts the goroutine responsible for calling
// registered handlers.
func (c *ServiceConfig) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	glog.Info("Starting service config controller")
	defer glog.Info("Shutting down service config controller")

	if !controller.WaitForCacheSync("service config", stopCh, c.listerSynced) {
		return
	}

	for i := range c.eventHandlers {
		glog.V(3).Infof("Calling handler.OnServiceSynced()")
		c.eventHandlers[i].OnServiceSynced()
	}

	<-stopCh
}

func (c *ServiceConfig) handleAddService(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
		return
	}
	for i := range c.eventHandlers {
		glog.V(4).Infof("Calling handler.OnServiceAdd")
		c.eventHandlers[i].OnServiceAdd(service)
	}
}

func (c *ServiceConfig) handleUpdateService(oldObj, newObj interface{}) {
	oldService, ok := oldObj.(*v1.Service)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", oldObj))
		return
	}
	service, ok := newObj.(*v1.Service)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", newObj))
		return
	}
	for i := range c.eventHandlers {
		glog.V(4).Infof("Calling handler.OnServiceUpdate")
		c.eventHandlers[i].OnServiceUpdate(oldService, service)
	}
}

func (c *ServiceConfig) handleDeleteService(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}
		if service, ok = tombstone.Obj.(*v1.Service); !ok {
			utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}
	}
	for i := range c.eventHandlers {
		glog.V(4).Infof("Calling handler.OnServiceDelete")
		c.eventHandlers[i].OnServiceDelete(service)
	}
}
