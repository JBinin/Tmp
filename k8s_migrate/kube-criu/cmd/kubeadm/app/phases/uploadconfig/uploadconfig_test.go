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

package uploadconfig

import (
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetfake "k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1alpha3 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
)

func TestUploadConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		errOnCreate    error
		errOnUpdate    error
		updateExisting bool
		errExpected    bool
		verifyResult   bool
	}{
		{
			name:         "basic validation with correct key",
			verifyResult: true,
		},
		{
			name:           "update existing should report no error",
			updateExisting: true,
			verifyResult:   true,
		},
		{
			name:        "unexpected errors for create should be returned",
			errOnCreate: apierrors.NewUnauthorized(""),
			errExpected: true,
		},
		{
			name:           "update existing show report error if unexpected error for update is returned",
			errOnUpdate:    apierrors.NewUnauthorized(""),
			updateExisting: true,
			errExpected:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			cfg := &kubeadmapi.InitConfiguration{
				KubernetesVersion: "v1.10.3",
				BootstrapTokens: []kubeadmapi.BootstrapToken{
					{
						Token: &kubeadmapi.BootstrapTokenString{
							ID:     "abcdef",
							Secret: "abcdef0123456789",
						},
					},
				},
				NodeRegistration: kubeadmapi.NodeRegistrationOptions{
					Name:      "node-foo",
					CRISocket: "/var/run/custom-cri.sock",
				},
			}
			client := clientsetfake.NewSimpleClientset()
			if tt.errOnCreate != nil {
				client.PrependReactor("create", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
					return true, nil, tt.errOnCreate
				})
			}
			// For idempotent test, we check the result of the second call.
			if err := UploadConfiguration(cfg, client); !tt.updateExisting && (err != nil) != tt.errExpected {
				t2.Fatalf("UploadConfiguration() error = %v, wantErr %v", err, tt.errExpected)
			}
			if tt.updateExisting {
				if tt.errOnUpdate != nil {
					client.PrependReactor("update", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
						return true, nil, tt.errOnUpdate
					})
				}
				if err := UploadConfiguration(cfg, client); (err != nil) != tt.errExpected {
					t2.Fatalf("UploadConfiguration() error = %v", err)
				}
			}
			if tt.verifyResult {
				masterCfg, err := client.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(kubeadmconstants.InitConfigurationConfigMap, metav1.GetOptions{})
				if err != nil {
					t2.Fatalf("Fail to query ConfigMap error = %v", err)
				}
				configData := masterCfg.Data[kubeadmconstants.InitConfigurationConfigMapKey]
				if configData == "" {
					t2.Fatalf("Fail to find ConfigMap key")
				}

				decodedCfg := &kubeadmapi.InitConfiguration{}
				if err := runtime.DecodeInto(kubeadmscheme.Codecs.UniversalDecoder(), []byte(configData), decodedCfg); err != nil {
					t2.Fatalf("unable to decode config from bytes: %v", err)
				}

				if decodedCfg.KubernetesVersion != cfg.KubernetesVersion {
					t2.Errorf("Decoded value doesn't match, decoded = %#v, expected = %#v", decodedCfg.KubernetesVersion, cfg.KubernetesVersion)
				}

				// If the decoded cfg has a BootstrapTokens array, verify the sensitive information we had isn't still there.
				if len(decodedCfg.BootstrapTokens) > 0 && decodedCfg.BootstrapTokens[0].Token != nil && decodedCfg.BootstrapTokens[0].Token.String() == cfg.BootstrapTokens[0].Token.String() {
					t2.Errorf("Decoded value contains .BootstrapTokens (sensitive info), decoded = %#v, expected = empty", decodedCfg.BootstrapTokens)
				}

				// Make sure no information from NodeRegistrationOptions was uploaded.
				if decodedCfg.NodeRegistration.Name == cfg.NodeRegistration.Name || decodedCfg.NodeRegistration.CRISocket != kubeadmapiv1alpha3.DefaultCRISocket {
					t2.Errorf("Decoded value contains .NodeRegistration (node-specific info shouldn't be uploaded), decoded = %#v, expected = empty", decodedCfg.NodeRegistration)
				}
			}
		})
	}
}
