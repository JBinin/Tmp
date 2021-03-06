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

package kubeadm

import "testing"

// kubeadmReset executes "kubeadm reset" and restarts kubelet.
func kubeadmReset() error {
	kubeadmPath := getKubeadmPath()
	_, _, err := RunCmd(kubeadmPath, "reset")
	return err
}

func TestCmdJoinConfig(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--config=foobar", false},
		{"--config=/does/not/exist/foo/bar", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinConfig running 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}

func TestCmdJoinDiscoveryFile(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--discovery-file=foobar", false},
		{"--discovery-file=file:wrong", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinDiscoveryFile running 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}

func TestCmdJoinDiscoveryToken(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--discovery-token=foobar", false},
		{"--discovery-token=token://asdf:asdf", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinDiscoveryToken running 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}

func TestCmdJoinNodeName(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--node-name=foobar", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinNodeName running 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}

func TestCmdJoinTLSBootstrapToken(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--tls-bootstrap-token=foobar", false},
		{"--tls-bootstrap-token=token://asdf:asdf", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinTLSBootstrapToken running 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}

func TestCmdJoinToken(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--token=foobar", false},
		{"--token=token://asdf:asdf", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinToken running 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}

func TestCmdJoinBadArgs(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	kubeadmPath := getKubeadmPath()
	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--discovery-token=abcdef.1234567890123456 --discovery-file=file:///tmp/foo.bar", false}, // DiscoveryToken, DiscoveryFile can't both be set
		{"", false}, // DiscoveryToken or DiscoveryFile must be set
	}

	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinBadArgs 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}

func TestCmdJoinArgsMixed(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var initTest = []struct {
		args     string
		expected bool
	}{
		{"--discovery-token=abcdef.1234567890abcdef --config=/etc/kubernetes/kubeadm.config", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range initTest {
		_, _, actual := RunCmd(kubeadmPath, "join", rt.args, "--ignore-preflight-errors=all")
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed CmdJoinArgsMixed running 'kubeadm join %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
				rt.args,
				actual,
				rt.expected,
				(actual == nil),
			)
		}
		kubeadmReset()
	}
}
