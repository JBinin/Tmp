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

package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	extensions "k8s.io/api/extensions/v1beta1"
)

func TestIngressToManifest(t *testing.T) {
	ing := &extensions.Ingress{}
	// Create a temp dir.
	tmpDir, err := ioutil.TempDir("", "kubemci")
	if err != nil {
		t.Fatalf("unexpected error in creating temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	ingPath := filepath.Join(tmpDir, "ing.yaml")

	// Write the ingress to a file and ensure that there is no error.
	if err := IngressToManifest(ing, ingPath); err != nil {
		t.Fatalf("Error in creating file: %s", err)
	}
	// Writing it again should not return an error.
	if err := IngressToManifest(ing, ingPath); err != nil {
		t.Fatalf("Error in creating file: %s", err)
	}
}
