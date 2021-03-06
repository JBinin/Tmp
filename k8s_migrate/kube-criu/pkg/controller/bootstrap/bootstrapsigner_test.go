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

package bootstrap

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	bootstrapapi "k8s.io/client-go/tools/bootstrap/token/api"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/controller"
)

func init() {
	spew.Config.DisableMethods = true
}

const testTokenID = "abc123"

func newBootstrapSigner() (*BootstrapSigner, *fake.Clientset, coreinformers.SecretInformer, coreinformers.ConfigMapInformer, error) {
	options := DefaultBootstrapSignerOptions()
	cl := fake.NewSimpleClientset()
	informers := informers.NewSharedInformerFactory(fake.NewSimpleClientset(), controller.NoResyncPeriodFunc())
	secrets := informers.Core().V1().Secrets()
	configMaps := informers.Core().V1().ConfigMaps()
	bsc, err := NewBootstrapSigner(cl, secrets, configMaps, options)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return bsc, cl, secrets, configMaps, nil
}

func newConfigMap(tokenID, signature string) *v1.ConfigMap {
	ret := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       metav1.NamespacePublic,
			Name:            bootstrapapi.ConfigMapClusterInfo,
			ResourceVersion: "1",
		},
		Data: map[string]string{
			bootstrapapi.KubeConfigKey: "payload",
		},
	}
	if len(tokenID) > 0 {
		ret.Data[bootstrapapi.JWSSignatureKeyPrefix+tokenID] = signature
	}
	return ret
}

func TestNoConfigMap(t *testing.T) {
	signer, cl, _, _, err := newBootstrapSigner()
	if err != nil {
		t.Fatalf("error creating BootstrapSigner: %v", err)
	}
	signer.signConfigMap()
	verifyActions(t, []core.Action{}, cl.Actions())
}

func TestSimpleSign(t *testing.T) {
	signer, cl, secrets, configMaps, err := newBootstrapSigner()
	if err != nil {
		t.Fatalf("error creating BootstrapSigner: %v", err)
	}

	cm := newConfigMap("", "")
	configMaps.Informer().GetIndexer().Add(cm)

	secret := newTokenSecret(testTokenID, "tokenSecret")
	addSecretSigningUsage(secret, "true")
	secrets.Informer().GetIndexer().Add(secret)

	signer.signConfigMap()

	expected := []core.Action{
		core.NewUpdateAction(schema.GroupVersionResource{Version: "v1", Resource: "configmaps"},
			api.NamespacePublic,
			newConfigMap(testTokenID, "eyJhbGciOiJIUzI1NiIsImtpZCI6ImFiYzEyMyJ9..QSxpUG7Q542CirTI2ECPSZjvBOJURUW5a7XqFpNI958")),
	}

	verifyActions(t, expected, cl.Actions())
}

func TestNoSignNeeded(t *testing.T) {
	signer, cl, secrets, configMaps, err := newBootstrapSigner()
	if err != nil {
		t.Fatalf("error creating BootstrapSigner: %v", err)
	}

	cm := newConfigMap(testTokenID, "eyJhbGciOiJIUzI1NiIsImtpZCI6ImFiYzEyMyJ9..QSxpUG7Q542CirTI2ECPSZjvBOJURUW5a7XqFpNI958")
	configMaps.Informer().GetIndexer().Add(cm)

	secret := newTokenSecret(testTokenID, "tokenSecret")
	addSecretSigningUsage(secret, "true")
	secrets.Informer().GetIndexer().Add(secret)

	signer.signConfigMap()

	verifyActions(t, []core.Action{}, cl.Actions())
}

func TestUpdateSignature(t *testing.T) {
	signer, cl, secrets, configMaps, err := newBootstrapSigner()
	if err != nil {
		t.Fatalf("error creating BootstrapSigner: %v", err)
	}

	cm := newConfigMap(testTokenID, "old signature")
	configMaps.Informer().GetIndexer().Add(cm)

	secret := newTokenSecret(testTokenID, "tokenSecret")
	addSecretSigningUsage(secret, "true")
	secrets.Informer().GetIndexer().Add(secret)

	signer.signConfigMap()

	expected := []core.Action{
		core.NewUpdateAction(schema.GroupVersionResource{Version: "v1", Resource: "configmaps"},
			api.NamespacePublic,
			newConfigMap(testTokenID, "eyJhbGciOiJIUzI1NiIsImtpZCI6ImFiYzEyMyJ9..QSxpUG7Q542CirTI2ECPSZjvBOJURUW5a7XqFpNI958")),
	}

	verifyActions(t, expected, cl.Actions())
}

func TestRemoveSignature(t *testing.T) {
	signer, cl, _, configMaps, err := newBootstrapSigner()
	if err != nil {
		t.Fatalf("error creating BootstrapSigner: %v", err)
	}

	cm := newConfigMap(testTokenID, "old signature")
	configMaps.Informer().GetIndexer().Add(cm)

	signer.signConfigMap()

	expected := []core.Action{
		core.NewUpdateAction(schema.GroupVersionResource{Version: "v1", Resource: "configmaps"},
			api.NamespacePublic,
			newConfigMap("", "")),
	}

	verifyActions(t, expected, cl.Actions())
}
