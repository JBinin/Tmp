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
package image

import (
	"os"
	"testing"

	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/cli/trust"
	"github.com/docker/docker/registry"
)

func unsetENV() {
	os.Unsetenv("DOCKER_CONTENT_TRUST")
	os.Unsetenv("DOCKER_CONTENT_TRUST_SERVER")
}

func TestENVTrustServer(t *testing.T) {
	defer unsetENV()
	indexInfo := &registrytypes.IndexInfo{Name: "testserver"}
	if err := os.Setenv("DOCKER_CONTENT_TRUST_SERVER", "https://notary-test.com:5000"); err != nil {
		t.Fatal("Failed to set ENV variable")
	}
	output, err := trust.Server(indexInfo)
	expectedStr := "https://notary-test.com:5000"
	if err != nil || output != expectedStr {
		t.Fatalf("Expected server to be %s, got %s", expectedStr, output)
	}
}

func TestHTTPENVTrustServer(t *testing.T) {
	defer unsetENV()
	indexInfo := &registrytypes.IndexInfo{Name: "testserver"}
	if err := os.Setenv("DOCKER_CONTENT_TRUST_SERVER", "http://notary-test.com:5000"); err != nil {
		t.Fatal("Failed to set ENV variable")
	}
	_, err := trust.Server(indexInfo)
	if err == nil {
		t.Fatal("Expected error with invalid scheme")
	}
}

func TestOfficialTrustServer(t *testing.T) {
	indexInfo := &registrytypes.IndexInfo{Name: "testserver", Official: true}
	output, err := trust.Server(indexInfo)
	if err != nil || output != registry.NotaryServer {
		t.Fatalf("Expected server to be %s, got %s", registry.NotaryServer, output)
	}
}

func TestNonOfficialTrustServer(t *testing.T) {
	indexInfo := &registrytypes.IndexInfo{Name: "testserver", Official: false}
	output, err := trust.Server(indexInfo)
	expectedStr := "https://" + indexInfo.Name
	if err != nil || output != expectedStr {
		t.Fatalf("Expected server to be %s, got %s", expectedStr, output)
	}
}
