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
Copyright 2017 The Kubernetes Authors.

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

package banflunder

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/sample-apiserver/pkg/admission/wardleinitializer"
	"k8s.io/sample-apiserver/pkg/apis/wardle"
	informers "k8s.io/sample-apiserver/pkg/client/informers/internalversion"
	listers "k8s.io/sample-apiserver/pkg/client/listers/wardle/internalversion"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register("BanFlunder", func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

type DisallowFlunder struct {
	*admission.Handler
	lister listers.FischerLister
}

var _ = wardleinitializer.WantsInternalWardleInformerFactory(&DisallowFlunder{})

// Admit ensures that the object in-flight is of kind Flunder.
// In addition checks that the Name is not on the banned list.
// The list is stored in Fischers API objects.
func (d *DisallowFlunder) Admit(a admission.Attributes) error {
	// we are only interested in flunders
	if a.GetKind().GroupKind() != wardle.Kind("Flunder") {
		return nil
	}

	if !d.WaitForReady() {
		return admission.NewForbidden(a, fmt.Errorf("not yet ready to handle request"))
	}

	metaAccessor, err := meta.Accessor(a.GetObject())
	if err != nil {
		return err
	}
	flunderName := metaAccessor.GetName()

	fischers, err := d.lister.List(labels.Everything())
	if err != nil {
		return err
	}

	for _, fischer := range fischers {
		for _, disallowedFlunder := range fischer.DisallowedFlunders {
			if flunderName == disallowedFlunder {
				return errors.NewForbidden(
					a.GetResource().GroupResource(),
					a.GetName(),
					fmt.Errorf("this name may not be used, please change the resource name"),
				)
			}
		}
	}
	return nil
}

// SetInternalWardleInformerFactory gets Lister from SharedInformerFactory.
// The lister knows how to lists Fischers.
func (d *DisallowFlunder) SetInternalWardleInformerFactory(f informers.SharedInformerFactory) {
	d.lister = f.Wardle().InternalVersion().Fischers().Lister()
	d.SetReadyFunc(f.Wardle().InternalVersion().Fischers().Informer().HasSynced)
}

// ValidaValidateInitializationte checks whether the plugin was correctly initialized.
func (d *DisallowFlunder) ValidateInitialization() error {
	if d.lister == nil {
		return fmt.Errorf("missing fischer lister")
	}
	return nil
}

// New creates a new ban flunder admission plugin
func New() (*DisallowFlunder, error) {
	return &DisallowFlunder{
		Handler: admission.NewHandler(admission.Create),
	}, nil
}
