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

package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	certutil "k8s.io/client-go/util/cert"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs/pkiutil"
)

type configMutatorsFunc func(*kubeadmapi.InitConfiguration, *certutil.Config) error

// KubeadmCert represents a certificate that Kubeadm will create to function properly.
type KubeadmCert struct {
	Name     string
	LongName string
	BaseName string
	CAName   string
	// Some attributes will depend on the InitConfiguration, only known at runtime.
	// These functions will be run in series, passed both the InitConfiguration and a cert Config.
	configMutators []configMutatorsFunc
	config         certutil.Config
}

// GetConfig returns the definition for the given cert given the provided InitConfiguration
func (k *KubeadmCert) GetConfig(ic *kubeadmapi.InitConfiguration) (*certutil.Config, error) {
	for _, f := range k.configMutators {
		if err := f(ic, &k.config); err != nil {
			return nil, err
		}
	}

	return &k.config, nil
}

// CreateFromCA makes and writes a certificate using the given CA cert and key.
func (k *KubeadmCert) CreateFromCA(ic *kubeadmapi.InitConfiguration, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	cfg, err := k.GetConfig(ic)
	if err != nil {
		return fmt.Errorf("couldn't create %q certificate: %v", k.Name, err)
	}
	cert, key, err := pkiutil.NewCertAndKey(caCert, caKey, cfg)
	if err != nil {
		return err
	}
	writeCertificateFilesIfNotExist(
		ic.CertificatesDir,
		k.BaseName,
		caCert,
		cert,
		key,
	)

	return nil
}

// CreateAsCA creates a certificate authority, writing the files to disk and also returning the created CA so it can be used to sign child certs.
func (k *KubeadmCert) CreateAsCA(ic *kubeadmapi.InitConfiguration) (*x509.Certificate, *rsa.PrivateKey, error) {
	cfg, err := k.GetConfig(ic)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get configuration for %q CA certificate: %v", k.Name, err)
	}
	caCert, caKey, err := NewCACertAndKey(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't generate %q CA certificate: %v", k.Name, err)
	}

	err = writeCertificateAuthorithyFilesIfNotExist(
		ic.CertificatesDir,
		k.BaseName,
		caCert,
		caKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't write out %q CA certificate: %v", k.Name, err)
	}

	return caCert, caKey, nil
}

// CertificateTree is represents a one-level-deep tree, mapping a CA to the certs that depend on it.
type CertificateTree map[*KubeadmCert]Certificates

// CreateTree creates the CAs, certs signed by the CAs, and writes them all to disk.
func (t CertificateTree) CreateTree(ic *kubeadmapi.InitConfiguration) error {
	for ca, leaves := range t {
		cfg, err := ca.GetConfig(ic)
		if err != nil {
			return err
		}

		caCert, caKey, err := NewCACertAndKey(cfg)
		if err != nil {
			return err
		}

		for _, leaf := range leaves {
			if err := leaf.CreateFromCA(ic, caCert, caKey); err != nil {
				return err
			}
		}

		err = writeCertificateAuthorithyFilesIfNotExist(
			ic.CertificatesDir,
			ca.BaseName,
			caCert,
			caKey,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// CertificateMap is a flat map of certificates, keyed by Name.
type CertificateMap map[string]*KubeadmCert

// CertTree returns a one-level-deep tree, mapping a CA cert to an array of certificates that should be signed by it.
func (m CertificateMap) CertTree() (CertificateTree, error) {
	caMap := make(CertificateTree)

	for _, cert := range m {
		if cert.CAName == "" {
			if _, ok := caMap[cert]; !ok {
				caMap[cert] = []*KubeadmCert{}
			}
		} else {
			ca, ok := m[cert.CAName]
			if !ok {
				return nil, fmt.Errorf("Certificate %q references unknown CA %q", cert.Name, cert.CAName)
			}
			caMap[ca] = append(caMap[ca], cert)
		}
	}

	return caMap, nil
}

// Certificates is a list of Certificates that Kubeadm should create.
type Certificates []*KubeadmCert

// AsMap returns the list of certificates as a map, keyed by name.
func (c Certificates) AsMap() CertificateMap {
	certMap := make(map[string]*KubeadmCert)
	for _, cert := range c {
		certMap[cert.Name] = cert
	}

	return certMap
}

// GetDefaultCertList returns  all of the certificates kubeadm requires to function.
func GetDefaultCertList() Certificates {
	return Certificates{
		&KubeadmCertRootCA,
		&KubeadmCertAPIServer,
		&KubeadmCertKubeletClient,
		// Front Proxy certs
		&KubeadmCertFrontProxyCA,
		&KubeadmCertFrontProxyClient,
		// etcd certs
		&KubeadmCertEtcdCA,
		&KubeadmCertEtcdServer,
		&KubeadmCertEtcdPeer,
		&KubeadmCertEtcdHealthcheck,
		&KubeadmCertEtcdAPIClient,
	}
}

// GetCertsWithoutEtcd returns all of the certificates kubeadm needs when etcd is hosted externally.
func GetCertsWithoutEtcd() Certificates {
	return Certificates{
		&KubeadmCertRootCA,
		&KubeadmCertAPIServer,
		&KubeadmCertKubeletClient,
		// Front Proxy certs
		&KubeadmCertFrontProxyCA,
		&KubeadmCertFrontProxyClient,
	}
}

var (
	// KubeadmCertRootCA is the definition of the Kubernetes Root CA for the API Server and kubelet.
	KubeadmCertRootCA = KubeadmCert{
		Name:     "ca",
		LongName: "self-signed kubernetes CA to provision identities for other kuberenets components",
		BaseName: kubeadmconstants.CACertAndKeyBaseName,
		config: certutil.Config{
			CommonName: "kubernetes",
		},
	}
	// KubeadmCertAPIServer is the definition of the cert used to serve the kubernetes API.
	KubeadmCertAPIServer = KubeadmCert{
		Name:     "apiserver",
		LongName: "certificate for serving the kubernetes API",
		BaseName: kubeadmconstants.APIServerCertAndKeyBaseName,
		CAName:   "ca",
		config: certutil.Config{
			CommonName: kubeadmconstants.APIServerCertCommonName,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		configMutators: []configMutatorsFunc{
			makeAltNamesMutator(pkiutil.GetAPIServerAltNames),
		},
	}
	// KubeadmCertKubeletClient is the definition of the cert used by the API server to access the kubelet.
	KubeadmCertKubeletClient = KubeadmCert{
		Name:     "apiserver-kubelet-client",
		LongName: "Client certificate for the API server to connect to kubelet",
		BaseName: kubeadmconstants.APIServerKubeletClientCertAndKeyBaseName,
		CAName:   "ca",
		config: certutil.Config{
			CommonName:   kubeadmconstants.APIServerKubeletClientCertCommonName,
			Organization: []string{kubeadmconstants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}

	// KubeadmCertFrontProxyCA is the definition of the CA used for the front end proxy.
	KubeadmCertFrontProxyCA = KubeadmCert{
		Name:     "front-proxy-ca",
		LongName: "self-signed CA to provision identities for front proxy",
		BaseName: kubeadmconstants.FrontProxyCACertAndKeyBaseName,
		config: certutil.Config{
			CommonName: "front-proxy-ca",
		},
	}

	// KubeadmCertFrontProxyClient is the definition of the cert used by the API server to access the front proxy.
	KubeadmCertFrontProxyClient = KubeadmCert{
		Name:     "front-proxy-client",
		BaseName: kubeadmconstants.FrontProxyClientCertAndKeyBaseName,
		LongName: "client for the front proxy",
		CAName:   "front-proxy-ca",
		config: certutil.Config{
			CommonName: kubeadmconstants.FrontProxyClientCertCommonName,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}

	// KubeadmCertEtcdCA is the definition of the root CA used by the hosted etcd server.
	KubeadmCertEtcdCA = KubeadmCert{
		Name:     "etcd-ca",
		LongName: "self-signed CA to provision identities for etcd",
		BaseName: kubeadmconstants.EtcdCACertAndKeyBaseName,
		config: certutil.Config{
			CommonName: "etcd-ca",
		},
	}
	// KubeadmCertEtcdServer is the definition of the cert used to serve etcd to clients.
	KubeadmCertEtcdServer = KubeadmCert{
		Name:     "etcd-server",
		LongName: "certificate for serving etcd",
		BaseName: kubeadmconstants.EtcdServerCertAndKeyBaseName,
		CAName:   "etcd-ca",
		config: certutil.Config{
			// TODO: etcd 3.2 introduced an undocumented requirement for ClientAuth usage on the
			// server cert: https://github.com/coreos/etcd/issues/9785#issuecomment-396715692
			// Once the upstream issue is resolved, this should be returned to only allowing
			// ServerAuth usage.
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		configMutators: []configMutatorsFunc{
			makeAltNamesMutator(pkiutil.GetEtcdAltNames),
			setCommonNameToNodeName(),
		},
	}
	// KubeadmCertEtcdPeer is the definition of the cert used by etcd peers to access each other.
	KubeadmCertEtcdPeer = KubeadmCert{
		Name:     "etcd-peer",
		LongName: "credentials for etcd nodes to communicate with each other",
		BaseName: kubeadmconstants.EtcdPeerCertAndKeyBaseName,
		CAName:   "etcd-ca",
		config: certutil.Config{
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		configMutators: []configMutatorsFunc{
			makeAltNamesMutator(pkiutil.GetEtcdPeerAltNames),
			setCommonNameToNodeName(),
		},
	}
	// KubeadmCertEtcdHealthcheck is the definition of the cert used by Kubernetes to check the health of the etcd server.
	KubeadmCertEtcdHealthcheck = KubeadmCert{
		Name:     "etcd-healthcheck-client",
		LongName: "client certificate for liveness probes to healtcheck etcd",
		BaseName: kubeadmconstants.EtcdHealthcheckClientCertAndKeyBaseName,
		CAName:   "etcd-ca",
		config: certutil.Config{
			CommonName:   kubeadmconstants.EtcdHealthcheckClientCertCommonName,
			Organization: []string{kubeadmconstants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
	// KubeadmCertEtcdAPIClient is the definition of the cert used by the API server to access etcd.
	KubeadmCertEtcdAPIClient = KubeadmCert{
		Name:     "apiserver-etcd-client",
		LongName: "client apiserver uses to access etcd",
		BaseName: kubeadmconstants.APIServerEtcdClientCertAndKeyBaseName,
		CAName:   "etcd-ca",
		config: certutil.Config{
			CommonName:   kubeadmconstants.APIServerEtcdClientCertCommonName,
			Organization: []string{kubeadmconstants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
)

func makeAltNamesMutator(f func(*kubeadmapi.InitConfiguration) (*certutil.AltNames, error)) configMutatorsFunc {
	return func(mc *kubeadmapi.InitConfiguration, cc *certutil.Config) error {
		altNames, err := f(mc)
		if err != nil {
			return err
		}
		cc.AltNames = *altNames
		return nil
	}
}

func setCommonNameToNodeName() configMutatorsFunc {
	return func(mc *kubeadmapi.InitConfiguration, cc *certutil.Config) error {
		cc.CommonName = mc.NodeRegistration.Name
		return nil
	}
}
