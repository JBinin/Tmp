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

package setdefault

import (
	"fmt"
	"io"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	admission "k8s.io/apiserver/pkg/admission"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/helper"
	"k8s.io/kubernetes/pkg/apis/storage"
	storageutil "k8s.io/kubernetes/pkg/apis/storage/util"
	informers "k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion"
	storagelisters "k8s.io/kubernetes/pkg/client/listers/storage/internalversion"
	kubeapiserveradmission "k8s.io/kubernetes/pkg/kubeapiserver/admission"
)

const (
	// PluginName is the name of this admission controller plugin
	PluginName = "DefaultStorageClass"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		plugin := newPlugin()
		return plugin, nil
	})
}

// claimDefaulterPlugin holds state for and implements the admission plugin.
type claimDefaulterPlugin struct {
	*admission.Handler

	lister storagelisters.StorageClassLister
}

var _ admission.Interface = &claimDefaulterPlugin{}
var _ admission.MutationInterface = &claimDefaulterPlugin{}
var _ = kubeapiserveradmission.WantsInternalKubeInformerFactory(&claimDefaulterPlugin{})

// newPlugin creates a new admission plugin.
func newPlugin() *claimDefaulterPlugin {
	return &claimDefaulterPlugin{
		Handler: admission.NewHandler(admission.Create),
	}
}

func (a *claimDefaulterPlugin) SetInternalKubeInformerFactory(f informers.SharedInformerFactory) {
	informer := f.Storage().InternalVersion().StorageClasses()
	a.lister = informer.Lister()
	a.SetReadyFunc(informer.Informer().HasSynced)
}

// ValidateInitialization ensures lister is set.
func (a *claimDefaulterPlugin) ValidateInitialization() error {
	if a.lister == nil {
		return fmt.Errorf("missing lister")
	}
	return nil
}

// Admit sets the default value of a PersistentVolumeClaim's storage class, in case the user did
// not provide a value.
//
// 1.  Find available StorageClasses.
// 2.  Figure which is the default
// 3.  Write to the PVClaim
func (a *claimDefaulterPlugin) Admit(attr admission.Attributes) error {
	if attr.GetResource().GroupResource() != api.Resource("persistentvolumeclaims") {
		return nil
	}

	if len(attr.GetSubresource()) != 0 {
		return nil
	}

	pvc, ok := attr.GetObject().(*api.PersistentVolumeClaim)
	// if we can't convert then we don't handle this object so just return
	if !ok {
		return nil
	}

	if helper.PersistentVolumeClaimHasClass(pvc) {
		// The user asked for a class.
		return nil
	}

	glog.V(4).Infof("no storage class for claim %s (generate: %s)", pvc.Name, pvc.GenerateName)

	def, err := getDefaultClass(a.lister)
	if err != nil {
		return admission.NewForbidden(attr, err)
	}
	if def == nil {
		// No default class selected, do nothing about the PVC.
		return nil
	}

	glog.V(4).Infof("defaulting storage class for claim %s (generate: %s) to %s", pvc.Name, pvc.GenerateName, def.Name)
	pvc.Spec.StorageClassName = &def.Name
	return nil
}

// getDefaultClass returns the default StorageClass from the store, or nil.
func getDefaultClass(lister storagelisters.StorageClassLister) (*storage.StorageClass, error) {
	list, err := lister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	defaultClasses := []*storage.StorageClass{}
	for _, class := range list {
		if storageutil.IsDefaultAnnotation(class.ObjectMeta) {
			defaultClasses = append(defaultClasses, class)
			glog.V(4).Infof("getDefaultClass added: %s", class.Name)
		}
	}

	if len(defaultClasses) == 0 {
		return nil, nil
	}
	if len(defaultClasses) > 1 {
		glog.V(4).Infof("getDefaultClass %d defaults found", len(defaultClasses))
		return nil, errors.NewInternalError(fmt.Errorf("%d default StorageClasses were found", len(defaultClasses)))
	}
	return defaultClasses[0], nil
}
