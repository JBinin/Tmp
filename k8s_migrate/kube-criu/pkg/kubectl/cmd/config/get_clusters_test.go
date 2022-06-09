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
Copyright 2014 The Kubernetes Authors.

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

package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type getClustersTest struct {
	config   clientcmdapi.Config
	expected string
}

func TestGetClusters(t *testing.T) {
	conf := clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			"minikube": {Server: "https://192.168.0.99"},
		},
	}
	test := getClustersTest{
		config: conf,
		expected: `NAME
minikube
`,
	}

	test.run(t)
}

func TestGetClustersEmpty(t *testing.T) {
	test := getClustersTest{
		config:   clientcmdapi.Config{},
		expected: "NAME\n",
	}

	test.run(t)
}

func (test getClustersTest) run(t *testing.T) {
	fakeKubeFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(fakeKubeFile.Name())
	err = clientcmd.WriteToFile(test.config, fakeKubeFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pathOptions := clientcmd.NewDefaultPathOptions()
	pathOptions.GlobalFile = fakeKubeFile.Name()
	pathOptions.EnvVar = ""
	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdConfigGetClusters(buf, pathOptions)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing command: %v", err)
	}
	if len(test.expected) != 0 {
		if buf.String() != test.expected {
			t.Errorf("expected %v, but got %v", test.expected, buf.String())
		}
		return
	}
}
