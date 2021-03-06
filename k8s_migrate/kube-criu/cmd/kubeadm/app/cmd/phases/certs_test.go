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

package phases

import (
	"fmt"
	"os"
	"strings"
	"testing"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs/pkiutil"
	"k8s.io/kubernetes/pkg/util/node"

	testutil "k8s.io/kubernetes/cmd/kubeadm/test"
	cmdtestutil "k8s.io/kubernetes/cmd/kubeadm/test/cmd"
)

// phaseTestK8sVersion is a fake kubernetes version to use when testing
const phaseTestK8sVersion = "v1.10.0"

func TestCertsSubCommandsHasFlags(t *testing.T) {

	subCmds := getCertsSubCommands(phaseTestK8sVersion)

	commonFlags := []string{
		"cert-dir",
		"config",
	}

	var tests = []struct {
		command         string
		additionalFlags []string
	}{
		{
			command: "all",
			additionalFlags: []string{
				"apiserver-advertise-address",
				"apiserver-cert-extra-sans",
				"service-cidr",
				"service-dns-domain",
			},
		},
		{
			command: "ca",
		},
		{
			command: "apiserver",
			additionalFlags: []string{
				"apiserver-advertise-address",
				"apiserver-cert-extra-sans",
				"service-cidr",
				"service-dns-domain",
			},
		},
		{
			command: "apiserver-kubelet-client",
		},
		{
			command: "etcd-ca",
		},
		{
			command: "etcd-server",
		},
		{
			command: "etcd-peer",
		},
		{
			command: "etcd-healthcheck-client",
		},
		{
			command: "apiserver-etcd-client",
		},
		{
			command: "sa",
		},
		{
			command: "front-proxy-ca",
		},
		{
			command: "front-proxy-client",
		},
	}

	for _, test := range tests {
		expectedFlags := append(commonFlags, test.additionalFlags...)
		cmdtestutil.AssertSubCommandHasFlags(t, subCmds, test.command, expectedFlags...)
	}
}

func TestSubCmdCertsCreateFilesWithFlags(t *testing.T) {

	subCmds := getCertsSubCommands(phaseTestK8sVersion)

	var tests = []struct {
		subCmds       []string
		expectedFiles []string
	}{
		{
			subCmds: []string{"all"},
			expectedFiles: []string{
				kubeadmconstants.CACertName, kubeadmconstants.CAKeyName,
				kubeadmconstants.APIServerCertName, kubeadmconstants.APIServerKeyName,
				kubeadmconstants.APIServerKubeletClientCertName, kubeadmconstants.APIServerKubeletClientKeyName,
				kubeadmconstants.ServiceAccountPrivateKeyName, kubeadmconstants.ServiceAccountPublicKeyName,
				kubeadmconstants.FrontProxyCACertName, kubeadmconstants.FrontProxyCAKeyName,
				kubeadmconstants.FrontProxyClientCertName, kubeadmconstants.FrontProxyClientKeyName,
			},
		},
		{
			subCmds:       []string{"ca", "apiserver", "apiserver-kubelet-client"},
			expectedFiles: []string{kubeadmconstants.CACertName, kubeadmconstants.CAKeyName, kubeadmconstants.APIServerCertName, kubeadmconstants.APIServerKeyName, kubeadmconstants.APIServerKubeletClientCertName, kubeadmconstants.APIServerKubeletClientKeyName},
		},
		{
			subCmds: []string{"etcd-ca", "etcd-server", "etcd-peer", "etcd-healthcheck-client", "apiserver-etcd-client"},
			expectedFiles: []string{
				kubeadmconstants.EtcdCACertName, kubeadmconstants.EtcdCAKeyName,
				kubeadmconstants.EtcdServerCertName, kubeadmconstants.EtcdServerKeyName,
				kubeadmconstants.EtcdPeerCertName, kubeadmconstants.EtcdPeerKeyName,
				kubeadmconstants.EtcdHealthcheckClientCertName, kubeadmconstants.EtcdHealthcheckClientKeyName,
				kubeadmconstants.APIServerEtcdClientCertName, kubeadmconstants.APIServerEtcdClientKeyName,
			},
		},
		{
			subCmds:       []string{"sa"},
			expectedFiles: []string{kubeadmconstants.ServiceAccountPrivateKeyName, kubeadmconstants.ServiceAccountPublicKeyName},
		},
		{
			subCmds:       []string{"front-proxy-ca", "front-proxy-client"},
			expectedFiles: []string{kubeadmconstants.FrontProxyCACertName, kubeadmconstants.FrontProxyCAKeyName, kubeadmconstants.FrontProxyClientCertName, kubeadmconstants.FrontProxyClientKeyName},
		},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.subCmds, ","), func(t *testing.T) {
			// Create temp folder for the test case
			tmpdir := testutil.SetupTempDir(t)
			defer os.RemoveAll(tmpdir)

			// executes given sub commands
			for _, subCmdName := range test.subCmds {
				fmt.Printf("running command %q\n", subCmdName)
				certDirFlag := fmt.Sprintf("--cert-dir=%s", tmpdir)
				cmdtestutil.RunSubCommand(t, subCmds, subCmdName, certDirFlag)
			}

			// verify expected files are there
			testutil.AssertFileExists(t, tmpdir, test.expectedFiles...)
		})
	}
}

func TestSubCmdCertsApiServerForwardsFlags(t *testing.T) {

	subCmds := getCertsSubCommands(phaseTestK8sVersion)

	// Create temp folder for the test case
	tmpdir := testutil.SetupTempDir(t)
	defer os.RemoveAll(tmpdir)

	// creates ca cert
	certDirFlag := fmt.Sprintf("--cert-dir=%s", tmpdir)
	cmdtestutil.RunSubCommand(t, subCmds, "ca", certDirFlag)

	// creates apiserver cert
	apiserverFlags := []string{
		fmt.Sprintf("--cert-dir=%s", tmpdir),
		"--apiserver-cert-extra-sans=foo,boo",
		"--service-cidr=10.0.0.0/24",
		"--service-dns-domain=mycluster.local",
		"--apiserver-advertise-address=1.2.3.4",
	}
	cmdtestutil.RunSubCommand(t, subCmds, "apiserver", apiserverFlags...)

	// asserts created cert has values from CLI flags
	APIserverCert, err := pkiutil.TryLoadCertFromDisk(tmpdir, kubeadmconstants.APIServerCertAndKeyBaseName)
	if err != nil {
		t.Fatalf("Error loading API server certificate: %v", err)
	}

	hostname, err := node.GetHostname("")
	if err != nil {
		t.Fatal(err)
	}

	for i, name := range []string{hostname, "kubernetes", "kubernetes.default", "kubernetes.default.svc", "kubernetes.default.svc.mycluster.local"} {
		if APIserverCert.DNSNames[i] != name {
			t.Errorf("APIserverCert.DNSNames[%d] is %s instead of %s", i, APIserverCert.DNSNames[i], name)
		}
	}
	for i, ip := range []string{"10.0.0.1", "1.2.3.4"} {
		if APIserverCert.IPAddresses[i].String() != ip {
			t.Errorf("APIserverCert.IPAddresses[%d] is %s instead of %s", i, APIserverCert.IPAddresses[i], ip)
		}
	}
}

func TestSubCmdCertsCreateFilesWithConfigFile(t *testing.T) {

	subCmds := getCertsSubCommands(phaseTestK8sVersion)

	var tests = []struct {
		subCmds       []string
		expectedFiles []string
	}{
		{
			subCmds: []string{"all"},
			expectedFiles: []string{
				kubeadmconstants.CACertName, kubeadmconstants.CAKeyName,
				kubeadmconstants.APIServerCertName, kubeadmconstants.APIServerKeyName,
				kubeadmconstants.APIServerKubeletClientCertName, kubeadmconstants.APIServerKubeletClientKeyName,
				kubeadmconstants.ServiceAccountPrivateKeyName, kubeadmconstants.ServiceAccountPublicKeyName,
				kubeadmconstants.FrontProxyCACertName, kubeadmconstants.FrontProxyCAKeyName,
				kubeadmconstants.FrontProxyClientCertName, kubeadmconstants.FrontProxyClientKeyName,
			},
		},
		{
			subCmds: []string{"ca", "apiserver", "apiserver-kubelet-client"},
			expectedFiles: []string{
				kubeadmconstants.CACertName, kubeadmconstants.CAKeyName,
				kubeadmconstants.APIServerCertName, kubeadmconstants.APIServerKeyName,
				kubeadmconstants.APIServerKubeletClientCertName, kubeadmconstants.APIServerKubeletClientKeyName,
			},
		},
		{
			subCmds: []string{"etcd-ca", "etcd-server", "etcd-peer", "etcd-healthcheck-client", "apiserver-etcd-client"},
			expectedFiles: []string{
				kubeadmconstants.EtcdCACertName, kubeadmconstants.EtcdCAKeyName,
				kubeadmconstants.EtcdServerCertName, kubeadmconstants.EtcdServerKeyName,
				kubeadmconstants.EtcdPeerCertName, kubeadmconstants.EtcdPeerKeyName,
				kubeadmconstants.EtcdHealthcheckClientCertName, kubeadmconstants.EtcdHealthcheckClientKeyName,
				kubeadmconstants.APIServerEtcdClientCertName, kubeadmconstants.APIServerEtcdClientKeyName,
			},
		},
		{
			subCmds:       []string{"front-proxy-ca", "front-proxy-client"},
			expectedFiles: []string{kubeadmconstants.FrontProxyCACertName, kubeadmconstants.FrontProxyCAKeyName, kubeadmconstants.FrontProxyClientCertName, kubeadmconstants.FrontProxyClientKeyName},
		},
		{
			subCmds:       []string{"sa"},
			expectedFiles: []string{kubeadmconstants.ServiceAccountPrivateKeyName, kubeadmconstants.ServiceAccountPublicKeyName},
		},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.subCmds, ","), func(t *testing.T) {
			// Create temp folder for the test case

			tmpdir := testutil.SetupTempDir(t)
			defer os.RemoveAll(tmpdir)

			cfg := &kubeadmapi.InitConfiguration{
				API:              kubeadmapi.API{AdvertiseAddress: "1.2.3.4", BindPort: 1234},
				CertificatesDir:  tmpdir,
				NodeRegistration: kubeadmapi.NodeRegistrationOptions{Name: "valid-node-name"},
			}
			configPath := testutil.SetupInitConfigurationFile(t, tmpdir, cfg)

			// executes given sub commands
			for _, subCmdName := range test.subCmds {
				t.Logf("running subcommand %q", subCmdName)
				configFlag := fmt.Sprintf("--config=%s", configPath)
				cmdtestutil.RunSubCommand(t, subCmds, subCmdName, configFlag)
			}

			// verify expected files are there
			testutil.AssertFileExists(t, tmpdir, test.expectedFiles...)
		})
	}
}
