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

package constants

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/kubernetes/pkg/util/version"
)

func TestGetStaticPodDirectory(t *testing.T) {
	expected := "/etc/kubernetes/manifests"
	actual := GetStaticPodDirectory()

	if actual != expected {
		t.Errorf(
			"failed GetStaticPodDirectory:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetAdminKubeConfigPath(t *testing.T) {
	expected := filepath.Join(KubernetesDir, AdminKubeConfigFileName)
	actual := GetAdminKubeConfigPath()

	if actual != expected {
		t.Errorf(
			"failed GetAdminKubeConfigPath:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetBootstrapKubeletKubeConfigPath(t *testing.T) {
	expected := "/etc/kubernetes/bootstrap-kubelet.conf"
	actual := GetBootstrapKubeletKubeConfigPath()

	if actual != expected {
		t.Errorf(
			"failed GetBootstrapKubeletKubeConfigPath:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetKubeletKubeConfigPath(t *testing.T) {
	expected := "/etc/kubernetes/kubelet.conf"
	actual := GetKubeletKubeConfigPath()

	if actual != expected {
		t.Errorf(
			"failed GetKubeletKubeConfigPath:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetStaticPodFilepath(t *testing.T) {
	var tests = []struct {
		componentName, manifestsDir, expected string
	}{
		{
			componentName: "kube-apiserver",
			manifestsDir:  "/etc/kubernetes/manifests",
			expected:      "/etc/kubernetes/manifests/kube-apiserver.yaml",
		},
		{
			componentName: "kube-controller-manager",
			manifestsDir:  "/etc/kubernetes/manifests/",
			expected:      "/etc/kubernetes/manifests/kube-controller-manager.yaml",
		},
		{
			componentName: "foo",
			manifestsDir:  "/etc/bar/",
			expected:      "/etc/bar/foo.yaml",
		},
	}
	for _, rt := range tests {
		actual := GetStaticPodFilepath(rt.componentName, rt.manifestsDir)
		if actual != rt.expected {
			t.Errorf(
				"failed GetStaticPodFilepath:\n\texpected: %s\n\t  actual: %s",
				rt.expected,
				actual,
			)
		}
	}
}

func TestAddSelfHostedPrefix(t *testing.T) {
	var tests = []struct {
		componentName, expected string
	}{
		{
			componentName: "kube-apiserver",
			expected:      "self-hosted-kube-apiserver",
		},
		{
			componentName: "kube-controller-manager",
			expected:      "self-hosted-kube-controller-manager",
		},
		{
			componentName: "kube-scheduler",
			expected:      "self-hosted-kube-scheduler",
		},
		{
			componentName: "foo",
			expected:      "self-hosted-foo",
		},
	}
	for _, rt := range tests {
		actual := AddSelfHostedPrefix(rt.componentName)
		if actual != rt.expected {
			t.Errorf(
				"failed AddSelfHostedPrefix:\n\texpected: %s\n\t  actual: %s",
				rt.expected,
				actual,
			)
		}
	}
}

func TestEtcdSupportedVersion(t *testing.T) {
	var tests = []struct {
		kubernetesVersion string
		expectedVersion   *version.Version
		expectedError     error
	}{
		{
			kubernetesVersion: "1.99.0",
			expectedVersion:   nil,
			expectedError:     fmt.Errorf("Unsupported or unknown kubernetes version(1.99.0)"),
		},
		{
			kubernetesVersion: "1.10.0",
			expectedVersion:   version.MustParseSemantic("3.1.12"),
			expectedError:     nil,
		},
		{
			kubernetesVersion: "1.10.2",
			expectedVersion:   version.MustParseSemantic("3.1.12"),
			expectedError:     nil,
		},
		{
			kubernetesVersion: "1.11.0",
			expectedVersion:   version.MustParseSemantic("3.2.18"),
			expectedError:     nil,
		},
		{
			kubernetesVersion: "1.12.1",
			expectedVersion:   version.MustParseSemantic("3.2.18"),
			expectedError:     nil,
		},
	}
	for _, rt := range tests {
		actualVersion, actualError := EtcdSupportedVersion(rt.kubernetesVersion)
		if actualError != nil {
			if rt.expectedError == nil {
				t.Errorf("failed EtcdSupportedVersion:\n\texpected no error, but got: %v", actualError)
			} else if actualError.Error() != rt.expectedError.Error() {
				t.Errorf(
					"failed EtcdSupportedVersion:\n\texpected error: %v\n\t  actual error: %v",
					rt.expectedError,
					actualError,
				)
			}
		} else {
			if rt.expectedError != nil {
				t.Errorf("failed EtcdSupportedVersion:\n\texpected error: %v, but got no error", rt.expectedError)
			} else if strings.Compare(actualVersion.String(), rt.expectedVersion.String()) != 0 {
				t.Errorf(
					"failed EtcdSupportedVersion:\n\texpected version: %s\n\t  actual version: %s",
					rt.expectedVersion.String(),
					actualVersion.String(),
				)
			}
		}
	}
}

func TestGetKubeDNSVersion(t *testing.T) {
	var tests = []struct {
		dns      string
		expected string
	}{
		{
			dns:      KubeDNS,
			expected: KubeDNSVersion,
		},
		{
			dns:      CoreDNS,
			expected: CoreDNSVersion,
		},
	}
	for _, rt := range tests {
		actualDNSVersion := GetDNSVersion(rt.dns)
		if actualDNSVersion != rt.expected {
			t.Errorf(
				"failed GetDNSVersion:\n\texpected: %s\n\t  actual: %s",
				rt.expected,
				actualDNSVersion,
			)
		}
	}
}
