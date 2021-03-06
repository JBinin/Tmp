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

package azure

import (
	"bytes"
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/go-autorest/autorest/to"
)

type fakeClient struct {
	results []containerregistry.Registry
}

func (f *fakeClient) List(ctx context.Context) ([]containerregistry.Registry, error) {
	return f.results, nil
}

func Test(t *testing.T) {
	configStr := `
    {
        "aadClientId": "foo",
        "aadClientSecret": "bar"
    }`
	result := []containerregistry.Registry{
		{
			Name: to.StringPtr("foo"),
			RegistryProperties: &containerregistry.RegistryProperties{
				LoginServer: to.StringPtr("*.azurecr.io"),
			},
		},
		{
			Name: to.StringPtr("bar"),
			RegistryProperties: &containerregistry.RegistryProperties{
				LoginServer: to.StringPtr("*.azurecr.cn"),
			},
		},
		{
			Name: to.StringPtr("baz"),
			RegistryProperties: &containerregistry.RegistryProperties{
				LoginServer: to.StringPtr("*.azurecr.de"),
			},
		},
		{
			Name: to.StringPtr("bus"),
			RegistryProperties: &containerregistry.RegistryProperties{
				LoginServer: to.StringPtr("*.azurecr.us"),
			},
		},
	}
	fakeClient := &fakeClient{
		results: result,
	}

	provider := &acrProvider{
		registryClient: fakeClient,
	}
	provider.loadConfig(bytes.NewBufferString(configStr))

	creds := provider.Provide()

	if len(creds) != len(result) {
		t.Errorf("Unexpected list: %v, expected length %d", creds, len(result))
	}
	for _, cred := range creds {
		if cred.Username != "foo" {
			t.Errorf("expected 'foo' for username, saw: %v", cred.Username)
		}
		if cred.Password != "bar" {
			t.Errorf("expected 'bar' for password, saw: %v", cred.Username)
		}
	}
	for _, val := range result {
		registryName := getLoginServer(val)
		if _, found := creds[registryName]; !found {
			t.Errorf("Missing expected registry: %s", registryName)
		}
	}
}
