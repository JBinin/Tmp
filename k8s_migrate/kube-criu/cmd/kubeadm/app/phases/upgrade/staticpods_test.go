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

package upgrade

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"k8s.io/apimachinery/pkg/runtime"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1alpha3 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	controlplanephase "k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"
	etcdphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/etcd"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	etcdutil "k8s.io/kubernetes/cmd/kubeadm/app/util/etcd"
)

const (
	waitForHashes        = "wait-for-hashes"
	waitForHashChange    = "wait-for-hash-change"
	waitForPodsWithLabel = "wait-for-pods-with-label"

	testConfiguration = `
apiVersion: kubeadm.k8s.io/v1alpha3
kind: InitConfiguration
api:
  advertiseAddress: 1.2.3.4
  bindPort: 6443
apiServerCertSANs: null
apiServerExtraArgs: null
certificatesDir: %s
controllerManagerExtraArgs: null
etcd:
  local:
    dataDir: %s
    image: ""
featureFlags: null
imageRepository: k8s.gcr.io
kubernetesVersion: %s
networking:
  dnsDomain: cluster.local
  podSubnet: ""
  serviceSubnet: 10.96.0.0/12
nodeRegistration:
  name: foo
  criSocket: ""
schedulerExtraArgs: null
token: ce3aa5.5ec8455bb76b379f
tokenTTL: 24h
unifiedControlPlaneImage: ""
`
)

// fakeWaiter is a fake apiclient.Waiter that returns errors it was initialized with
type fakeWaiter struct {
	errsToReturn map[string]error
}

func NewFakeStaticPodWaiter(errsToReturn map[string]error) apiclient.Waiter {
	return &fakeWaiter{
		errsToReturn: errsToReturn,
	}
}

// WaitForAPI just returns a dummy nil, to indicate that the program should just proceed
func (w *fakeWaiter) WaitForAPI() error {
	return nil
}

// WaitForPodsWithLabel just returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForPodsWithLabel(kvLabel string) error {
	return w.errsToReturn[waitForPodsWithLabel]
}

// WaitForPodToDisappear just returns a dummy nil, to indicate that the program should just proceed
func (w *fakeWaiter) WaitForPodToDisappear(podName string) error {
	return nil
}

// SetTimeout is a no-op; we don't use it in this implementation
func (w *fakeWaiter) SetTimeout(_ time.Duration) {}

// WaitForStaticPodControlPlaneHashes returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForStaticPodControlPlaneHashes(_ string) (map[string]string, error) {
	return map[string]string{}, w.errsToReturn[waitForHashes]
}

// WaitForStaticPodSingleHash returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForStaticPodSingleHash(_ string, _ string) (string, error) {
	return "", w.errsToReturn[waitForHashes]
}

// WaitForStaticPodHashChange returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForStaticPodHashChange(_, _, _ string) error {
	return w.errsToReturn[waitForHashChange]
}

// WaitForHealthyKubelet returns a dummy nil just to implement the interface
func (w *fakeWaiter) WaitForHealthyKubelet(_ time.Duration, _ string) error {
	return nil
}

type fakeStaticPodPathManager struct {
	kubernetesDir     string
	realManifestDir   string
	tempManifestDir   string
	backupManifestDir string
	backupEtcdDir     string
	MoveFileFunc      func(string, string) error
}

func NewFakeStaticPodPathManager(moveFileFunc func(string, string) error) (StaticPodPathManager, error) {
	kubernetesDir, err := ioutil.TempDir("", "kubeadm-pathmanager-")
	if err != nil {
		return nil, fmt.Errorf("couldn't create a temporary directory for the upgrade: %v", err)
	}

	realManifestDir := filepath.Join(kubernetesDir, constants.ManifestsSubDirName)
	if err := os.Mkdir(realManifestDir, 0700); err != nil {
		return nil, fmt.Errorf("couldn't create a realManifestDir for the upgrade: %v", err)
	}

	upgradedManifestDir := filepath.Join(kubernetesDir, "upgraded-manifests")
	if err := os.Mkdir(upgradedManifestDir, 0700); err != nil {
		return nil, fmt.Errorf("couldn't create a upgradedManifestDir for the upgrade: %v", err)
	}

	backupManifestDir := filepath.Join(kubernetesDir, "backup-manifests")
	if err := os.Mkdir(backupManifestDir, 0700); err != nil {
		return nil, fmt.Errorf("couldn't create a backupManifestDir for the upgrade: %v", err)
	}

	backupEtcdDir := filepath.Join(kubernetesDir, "kubeadm-backup-etcd")
	if err := os.Mkdir(backupEtcdDir, 0700); err != nil {
		return nil, err
	}

	return &fakeStaticPodPathManager{
		kubernetesDir:     kubernetesDir,
		realManifestDir:   realManifestDir,
		tempManifestDir:   upgradedManifestDir,
		backupManifestDir: backupManifestDir,
		backupEtcdDir:     backupEtcdDir,
		MoveFileFunc:      moveFileFunc,
	}, nil
}

func (spm *fakeStaticPodPathManager) MoveFile(oldPath, newPath string) error {
	return spm.MoveFileFunc(oldPath, newPath)
}

func (spm *fakeStaticPodPathManager) KubernetesDir() string {
	return spm.kubernetesDir
}

func (spm *fakeStaticPodPathManager) RealManifestPath(component string) string {
	return constants.GetStaticPodFilepath(component, spm.realManifestDir)
}
func (spm *fakeStaticPodPathManager) RealManifestDir() string {
	return spm.realManifestDir
}

func (spm *fakeStaticPodPathManager) TempManifestPath(component string) string {
	return constants.GetStaticPodFilepath(component, spm.tempManifestDir)
}
func (spm *fakeStaticPodPathManager) TempManifestDir() string {
	return spm.tempManifestDir
}

func (spm *fakeStaticPodPathManager) BackupManifestPath(component string) string {
	return constants.GetStaticPodFilepath(component, spm.backupManifestDir)
}
func (spm *fakeStaticPodPathManager) BackupManifestDir() string {
	return spm.backupManifestDir
}

func (spm *fakeStaticPodPathManager) BackupEtcdDir() string {
	return spm.backupEtcdDir
}

func (spm *fakeStaticPodPathManager) CleanupDirs() error {
	if err := os.RemoveAll(spm.TempManifestDir()); err != nil {
		return err
	}
	if err := os.RemoveAll(spm.BackupManifestDir()); err != nil {
		return err
	}
	return os.RemoveAll(spm.BackupEtcdDir())
}

type fakeTLSEtcdClient struct{ TLS bool }

func (c fakeTLSEtcdClient) HasTLS() bool {
	return c.TLS
}

func (c fakeTLSEtcdClient) ClusterAvailable() (bool, error) { return true, nil }

func (c fakeTLSEtcdClient) WaitForClusterAvailable(delay time.Duration, retries int, retryInterval time.Duration) (bool, error) {
	return true, nil
}

func (c fakeTLSEtcdClient) GetClusterStatus() (map[string]*clientv3.StatusResponse, error) {
	return map[string]*clientv3.StatusResponse{
		"foo": {
			Version: "3.1.12",
		}}, nil
}

func (c fakeTLSEtcdClient) GetClusterVersions() (map[string]string, error) {
	return map[string]string{
		"foo": "3.1.12",
	}, nil
}

func (c fakeTLSEtcdClient) GetVersion() (string, error) {
	return "3.1.12", nil
}

type fakePodManifestEtcdClient struct{ ManifestDir, CertificatesDir string }

func (c fakePodManifestEtcdClient) HasTLS() bool {
	hasTLS, _ := etcdutil.PodManifestsHaveTLS(c.ManifestDir)
	return hasTLS
}

func (c fakePodManifestEtcdClient) ClusterAvailable() (bool, error) { return true, nil }

func (c fakePodManifestEtcdClient) WaitForClusterAvailable(delay time.Duration, retries int, retryInterval time.Duration) (bool, error) {
	return true, nil
}

func (c fakePodManifestEtcdClient) GetClusterStatus() (map[string]*clientv3.StatusResponse, error) {
	// Make sure the certificates generated from the upgrade are readable from disk
	tlsInfo := transport.TLSInfo{
		CertFile:      filepath.Join(c.CertificatesDir, constants.EtcdCACertName),
		KeyFile:       filepath.Join(c.CertificatesDir, constants.EtcdHealthcheckClientCertName),
		TrustedCAFile: filepath.Join(c.CertificatesDir, constants.EtcdHealthcheckClientKeyName),
	}
	_, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	return map[string]*clientv3.StatusResponse{
		"foo": {Version: "3.1.12"},
	}, nil
}

func (c fakePodManifestEtcdClient) GetClusterVersions() (map[string]string, error) {
	return map[string]string{
		"foo": "3.1.12",
	}, nil
}

func (c fakePodManifestEtcdClient) GetVersion() (string, error) {
	return "3.1.12", nil
}

func TestStaticPodControlPlane(t *testing.T) {
	tests := []struct {
		description          string
		waitErrsToReturn     map[string]error
		moveFileFunc         func(string, string) error
		expectedErr          bool
		manifestShouldChange bool
	}{
		{
			description: "error-free case should succeed",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          false,
			manifestShouldChange: true,
		},
		{
			description: "any wait error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        fmt.Errorf("boo! failed"),
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any wait error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    fmt.Errorf("boo! failed"),
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any wait error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: fmt.Errorf("boo! failed"),
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any path-moving error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				// fail for kube-apiserver move
				if strings.Contains(newPath, "kube-apiserver") {
					return fmt.Errorf("moving the kube-apiserver file failed")
				}
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any path-moving error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				// fail for kube-controller-manager move
				if strings.Contains(newPath, "kube-controller-manager") {
					return fmt.Errorf("moving the kube-apiserver file failed")
				}
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any path-moving error should result in a rollback and an abort; even though this is the last component (kube-apiserver and kube-controller-manager healthy)",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				// fail for kube-scheduler move
				if strings.Contains(newPath, "kube-scheduler") {
					return fmt.Errorf("moving the kube-apiserver file failed")
				}
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
	}

	for _, rt := range tests {
		waiter := NewFakeStaticPodWaiter(rt.waitErrsToReturn)
		pathMgr, err := NewFakeStaticPodPathManager(rt.moveFileFunc)
		if err != nil {
			t.Fatalf("couldn't run NewFakeStaticPodPathManager: %v", err)
		}
		defer os.RemoveAll(pathMgr.(*fakeStaticPodPathManager).KubernetesDir())
		constants.KubernetesDir = pathMgr.(*fakeStaticPodPathManager).KubernetesDir()

		tempCertsDir, err := ioutil.TempDir("", "kubeadm-certs")
		if err != nil {
			t.Fatalf("couldn't create temporary certificates directory: %v", err)
		}
		defer os.RemoveAll(tempCertsDir)
		tmpEtcdDataDir, err := ioutil.TempDir("", "kubeadm-etcd-data")
		if err != nil {
			t.Fatalf("couldn't create temporary etcd data directory: %v", err)
		}
		defer os.RemoveAll(tmpEtcdDataDir)

		oldcfg, err := getConfig("v1.9.0", tempCertsDir, tmpEtcdDataDir)
		if err != nil {
			t.Fatalf("couldn't create config: %v", err)
		}

		tree, err := certsphase.GetCertsWithoutEtcd().AsMap().CertTree()
		if err != nil {
			t.Fatalf("couldn't get cert tree: %v", err)
		}

		if err := tree.CreateTree(oldcfg); err != nil {
			t.Fatalf("couldn't get create cert tree: %v", err)
		}

		t.Logf("Wrote certs to %s\n", oldcfg.CertificatesDir)

		// Initialize the directory with v1.7 manifests; should then be upgraded to v1.8 using the method
		err = controlplanephase.CreateInitStaticPodManifestFiles(pathMgr.RealManifestDir(), oldcfg)
		if err != nil {
			t.Fatalf("couldn't run CreateInitStaticPodManifestFiles: %v", err)
		}
		err = etcdphase.CreateLocalEtcdStaticPodManifestFile(pathMgr.RealManifestDir(), oldcfg)
		if err != nil {
			t.Fatalf("couldn't run CreateLocalEtcdStaticPodManifestFile: %v", err)
		}
		// Get a hash of the v1.7 API server manifest to compare later (was the file re-written)
		oldHash, err := getAPIServerHash(pathMgr.RealManifestDir())
		if err != nil {
			t.Fatalf("couldn't read temp file: %v", err)
		}

		newcfg, err := getConfig("v1.10.0", tempCertsDir, tmpEtcdDataDir)
		if err != nil {
			t.Fatalf("couldn't create config: %v", err)
		}

		actualErr := StaticPodControlPlane(
			waiter,
			pathMgr,
			newcfg,
			true,
			fakeTLSEtcdClient{
				TLS: false,
			},
			fakePodManifestEtcdClient{
				ManifestDir:     pathMgr.RealManifestDir(),
				CertificatesDir: newcfg.CertificatesDir,
			},
		)
		if (actualErr != nil) != rt.expectedErr {
			t.Errorf(
				"failed UpgradeStaticPodControlPlane\n%s\n\texpected error: %t\n\tgot: %t\n\tactual error: %v",
				rt.description,
				rt.expectedErr,
				(actualErr != nil),
				actualErr,
			)
		}

		newHash, err := getAPIServerHash(pathMgr.RealManifestDir())
		if err != nil {
			t.Fatalf("couldn't read temp file: %v", err)
		}

		if (oldHash != newHash) != rt.manifestShouldChange {
			t.Errorf(
				"failed StaticPodControlPlane\n%s\n\texpected manifest change: %t\n\tgot: %t",
				rt.description,
				rt.manifestShouldChange,
				(oldHash != newHash),
			)
		}
		return
	}
}

func getAPIServerHash(dir string) (string, error) {
	manifestPath := constants.GetStaticPodFilepath(constants.KubeAPIServer, dir)

	fileBytes, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(fileBytes)), nil
}

// TODO: Make this test function use the rest of the "official" API machinery helper funcs we have inside of kubeadm
func getConfig(version, certsDir, etcdDataDir string) (*kubeadmapi.InitConfiguration, error) {
	externalcfg := &kubeadmapiv1alpha3.InitConfiguration{}
	internalcfg := &kubeadmapi.InitConfiguration{}
	if err := runtime.DecodeInto(kubeadmscheme.Codecs.UniversalDecoder(), []byte(fmt.Sprintf(testConfiguration, certsDir, etcdDataDir, version)), externalcfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %v", err)
	}
	kubeadmscheme.Scheme.Convert(externalcfg, internalcfg, nil)
	return internalcfg, nil
}

func getTempDir(t *testing.T, name string) (string, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		t.Fatalf("couldn't make temporary directory: %v", err)
	}

	return dir, func() {
		os.RemoveAll(dir)
	}
}

func TestCleanupDirs(t *testing.T) {
	tests := []struct {
		name                   string
		keepManifest, keepEtcd bool
	}{
		{
			name:         "save manifest backup",
			keepManifest: true,
		},
		{
			name:         "save both etcd and manifest",
			keepManifest: true,
			keepEtcd:     true,
		},
		{
			name: "save nothing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			realManifestDir, cleanup := getTempDir(t, "realManifestDir")
			defer cleanup()

			tempManifestDir, cleanup := getTempDir(t, "tempManifestDir")
			defer cleanup()

			backupManifestDir, cleanup := getTempDir(t, "backupManifestDir")
			defer cleanup()

			backupEtcdDir, cleanup := getTempDir(t, "backupEtcdDir")
			defer cleanup()

			mgr := NewKubeStaticPodPathManager(realManifestDir, tempManifestDir, backupManifestDir, backupEtcdDir, test.keepManifest, test.keepEtcd)
			err := mgr.CleanupDirs()
			if err != nil {
				t.Errorf("unexpected error cleaning up: %v", err)
			}

			if _, err := os.Stat(tempManifestDir); !os.IsNotExist(err) {
				t.Errorf("%q should not have existed", tempManifestDir)
			}
			_, err = os.Stat(backupManifestDir)
			if test.keepManifest {
				if err != nil {
					t.Errorf("unexpected error getting backup manifest dir")
				}
			} else {
				if !os.IsNotExist(err) {
					t.Error("expected backup manifest to not exist")
				}
			}

			_, err = os.Stat(backupEtcdDir)
			if test.keepEtcd {
				if err != nil {
					t.Errorf("unexpected error getting backup etcd dir")
				}
			} else {
				if !os.IsNotExist(err) {
					t.Error("expected backup etcd dir to not exist")
				}
			}
		})
	}
}
