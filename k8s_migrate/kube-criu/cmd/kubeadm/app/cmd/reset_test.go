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

package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	kubeadmapiv1alpha3 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/validation"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/preflight"
	"k8s.io/utils/exec"
	fakeexec "k8s.io/utils/exec/testing"
)

func assertExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("file/directory does not exist; error: %s", err)
		t.Errorf("file/directory does not exist: %s", path)
	}
}

func assertNotExists(t *testing.T, path string) {
	if _, err := os.Stat(path); err == nil {
		t.Errorf("file/dir exists: %s", path)
	}
}

// assertDirEmpty verifies a directory either does not exist, or is empty.
func assertDirEmpty(t *testing.T, path string) {
	dac := preflight.DirAvailableCheck{Path: path}
	_, errors := dac.Check()
	if len(errors) != 0 {
		t.Errorf("directory not empty: [%v]", errors)
	}
}

func TestNewReset(t *testing.T) {
	var in io.Reader
	certsDir := kubeadmapiv1alpha3.DefaultCertificatesDir
	criSocketPath := kubeadmapiv1alpha3.DefaultCRISocket
	forceReset := true

	ignorePreflightErrors := []string{"all"}
	ignorePreflightErrorsSet, _ := validation.ValidateIgnorePreflightErrors(ignorePreflightErrors)
	NewReset(in, ignorePreflightErrorsSet, forceReset, certsDir, criSocketPath)

	ignorePreflightErrors = []string{}
	ignorePreflightErrorsSet, _ = validation.ValidateIgnorePreflightErrors(ignorePreflightErrors)
	NewReset(in, ignorePreflightErrorsSet, forceReset, certsDir, criSocketPath)
}

func TestNewCmdReset(t *testing.T) {
	var out io.Writer
	var in io.Reader
	cmd := NewCmdReset(in, out)

	tmpDir, err := ioutil.TempDir("", "kubeadm-reset-test")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	args := []string{"--ignore-preflight-errors=all", "--cert-dir=" + tmpDir, "--force"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("Cannot execute reset command: %v", err)
	}
}

func TestConfigDirCleaner(t *testing.T) {
	tests := map[string]struct {
		resetDir        string
		setupDirs       []string
		setupFiles      []string
		verifyExists    []string
		verifyNotExists []string
	}{
		"simple reset": {
			setupDirs: []string{
				"manifests",
				"pki",
			},
			setupFiles: []string{
				"manifests/etcd.yaml",
				"manifests/kube-apiserver.yaml",
				"pki/ca.pem",
				kubeadmconstants.AdminKubeConfigFileName,
				kubeadmconstants.KubeletKubeConfigFileName,
			},
			verifyExists: []string{
				"manifests",
				"pki",
			},
		},
		"partial reset": {
			setupDirs: []string{
				"pki",
			},
			setupFiles: []string{
				"pki/ca.pem",
				kubeadmconstants.KubeletKubeConfigFileName,
			},
			verifyExists: []string{
				"pki",
			},
			verifyNotExists: []string{
				"manifests",
			},
		},
		"preserve unrelated file foo": {
			setupDirs: []string{
				"manifests",
				"pki",
			},
			setupFiles: []string{
				"manifests/etcd.yaml",
				"manifests/kube-apiserver.yaml",
				"pki/ca.pem",
				kubeadmconstants.AdminKubeConfigFileName,
				kubeadmconstants.KubeletKubeConfigFileName,
				"foo",
			},
			verifyExists: []string{
				"manifests",
				"pki",
				"foo",
			},
		},
		"preserve hidden files and directories": {
			setupDirs: []string{
				"manifests",
				"pki",
				".mydir",
			},
			setupFiles: []string{
				"manifests/etcd.yaml",
				"manifests/kube-apiserver.yaml",
				"pki/ca.pem",
				kubeadmconstants.AdminKubeConfigFileName,
				kubeadmconstants.KubeletKubeConfigFileName,
				".mydir/.myfile",
			},
			verifyExists: []string{
				"manifests",
				"pki",
				".mydir",
				".mydir/.myfile",
			},
		},
		"no-op reset": {
			verifyNotExists: []string{
				"pki",
				"manifests",
			},
		},
		"not a directory": {
			resetDir: "test-path",
			setupFiles: []string{
				"test-path",
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test: %s", name)

		// Create a temporary directory for our fake config dir:
		tmpDir, err := ioutil.TempDir("", "kubeadm-reset-test")
		if err != nil {
			t.Errorf("Unable to create temporary directory: %s", err)
		}

		for _, createDir := range test.setupDirs {
			err := os.Mkdir(filepath.Join(tmpDir, createDir), 0700)
			if err != nil {
				t.Errorf("Unable to setup test config directory: %s", err)
			}
		}

		for _, createFile := range test.setupFiles {
			fullPath := filepath.Join(tmpDir, createFile)
			f, err := os.Create(fullPath)
			if err != nil {
				t.Errorf("Unable to create test file: %s", err)
			}
			f.Close()
		}

		if test.resetDir == "" {
			test.resetDir = "pki"
		}
		resetConfigDir(tmpDir, filepath.Join(tmpDir, test.resetDir))

		// Verify the files we cleanup implicitly in every test:
		assertExists(t, tmpDir)
		assertNotExists(t, filepath.Join(tmpDir, kubeadmconstants.AdminKubeConfigFileName))
		assertNotExists(t, filepath.Join(tmpDir, kubeadmconstants.KubeletKubeConfigFileName))
		assertDirEmpty(t, filepath.Join(tmpDir, "manifests"))
		assertDirEmpty(t, filepath.Join(tmpDir, "pki"))

		// Verify the files as requested by the test:
		for _, path := range test.verifyExists {
			assertExists(t, filepath.Join(tmpDir, path))
		}
		for _, path := range test.verifyNotExists {
			assertNotExists(t, filepath.Join(tmpDir, path))
		}

		os.RemoveAll(tmpDir)
	}
}

func TestRemoveContainers(t *testing.T) {
	fcmd := fakeexec.FakeCmd{
		CombinedOutputScript: []fakeexec.FakeCombinedOutputAction{
			func() ([]byte, error) { return []byte("id1\nid2"), nil },
			func() ([]byte, error) { return []byte(""), nil },
			func() ([]byte, error) { return []byte(""), nil },
			func() ([]byte, error) { return []byte(""), nil },
			func() ([]byte, error) { return []byte(""), nil },
		},
	}
	fexec := fakeexec.FakeExec{
		CommandScript: []fakeexec.FakeCommandAction{
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
		},
		LookPathFunc: func(cmd string) (string, error) { return "/usr/bin/crictl", nil },
	}

	removeContainers(&fexec, "unix:///var/run/crio/crio.sock")
}
