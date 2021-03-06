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

package wardleinitializer_test

import (
	"testing"
	"time"

	"k8s.io/apiserver/pkg/admission"
	"k8s.io/sample-apiserver/pkg/admission/wardleinitializer"
	"k8s.io/sample-apiserver/pkg/client/clientset/internalversion/fake"
	informers "k8s.io/sample-apiserver/pkg/client/informers/internalversion"
)

// TestWantsInternalWardleInformerFactory ensures that the informer factory is injected
// when the WantsInternalWardleInformerFactory interface is implemented by a plugin.
func TestWantsInternalWardleInformerFactory(t *testing.T) {
	cs := &fake.Clientset{}
	sf := informers.NewSharedInformerFactory(cs, time.Duration(1)*time.Second)
	target := wardleinitializer.New(sf)

	wantWardleInformerFactory := &wantInternalWardleInformerFactory{}
	target.Initialize(wantWardleInformerFactory)
	if wantWardleInformerFactory.sf != sf {
		t.Errorf("expected informer factory to be initialized")
	}
}

// wantInternalWardleInformerFactory is a test stub that fulfills the WantsInternalWardleInformerFactory interface
type wantInternalWardleInformerFactory struct {
	sf informers.SharedInformerFactory
}

func (self *wantInternalWardleInformerFactory) SetInternalWardleInformerFactory(sf informers.SharedInformerFactory) {
	self.sf = sf
}
func (self *wantInternalWardleInformerFactory) Admit(a admission.Attributes) error { return nil }
func (self *wantInternalWardleInformerFactory) Handles(o admission.Operation) bool { return false }
func (self *wantInternalWardleInformerFactory) ValidateInitialization() error      { return nil }

var _ admission.Interface = &wantInternalWardleInformerFactory{}
var _ wardleinitializer.WantsInternalWardleInformerFactory = &wantInternalWardleInformerFactory{}
