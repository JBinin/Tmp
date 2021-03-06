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
Copyright 2015 The Kubernetes Authors.

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

package kubectl

import (
	"encoding/json"
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/credentialprovider"
	"k8s.io/kubernetes/pkg/kubectl/util/hash"
)

// SecretForDockerRegistryGeneratorV1 supports stable generation of a docker registry secret
type SecretForDockerRegistryGeneratorV1 struct {
	// Name of secret (required)
	Name string
	// FileSources to derive the secret from (optional)
	FileSources []string
	// Username for registry (required)
	Username string
	// Email for registry (optional)
	Email string
	// Password for registry (required)
	Password string
	// Server for registry (required)
	Server string
	// AppendHash; if true, derive a hash from the Secret and append it to the name
	AppendHash bool
}

// Ensure it supports the generator pattern that uses parameter injection
var _ Generator = &SecretForDockerRegistryGeneratorV1{}

// Ensure it supports the generator pattern that uses parameters specified during construction
var _ StructuredGenerator = &SecretForDockerRegistryGeneratorV1{}

// Generate returns a secret using the specified parameters
func (s SecretForDockerRegistryGeneratorV1) Generate(genericParams map[string]interface{}) (runtime.Object, error) {
	err := ValidateParams(s.ParamNames(), genericParams)
	if err != nil {
		return nil, err
	}
	delegate := &SecretForDockerRegistryGeneratorV1{}
	hashParam, found := genericParams["append-hash"]
	if found {
		hashBool, isBool := hashParam.(bool)
		if !isBool {
			return nil, fmt.Errorf("expected bool, found :%v", hashParam)
		}
		delegate.AppendHash = hashBool
		delete(genericParams, "append-hash")
	}
	params := map[string]string{}
	for key, value := range genericParams {
		strVal, isString := value.(string)
		if !isString {
			return nil, fmt.Errorf("expected string, saw %v for '%s'", value, key)
		}
		params[key] = strVal
	}
	delegate.Name = params["name"]
	delegate.Username = params["docker-username"]
	delegate.Email = params["docker-email"]
	delegate.Password = params["docker-password"]
	delegate.Server = params["docker-server"]
	return delegate.StructuredGenerate()
}

// StructuredGenerate outputs a secret object using the configured fields
func (s SecretForDockerRegistryGeneratorV1) StructuredGenerate() (runtime.Object, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	secret := &v1.Secret{}
	secret.Name = s.Name
	secret.Type = v1.SecretTypeDockerConfigJson
	secret.Data = map[string][]byte{}
	if len(s.FileSources) > 0 {
		if err := handleFromFileSources(secret, s.FileSources); err != nil {
			return nil, err
		}
	}
	if len(s.FileSources) == 0 {
		dockercfgJsonContent, err := handleDockerCfgJsonContent(s.Username, s.Password, s.Email, s.Server)
		if err != nil {
			return nil, err
		}
		secret.Data[v1.DockerConfigJsonKey] = dockercfgJsonContent
	}
	if s.AppendHash {
		h, err := hash.SecretHash(secret)
		if err != nil {
			return nil, err
		}
		secret.Name = fmt.Sprintf("%s-%s", secret.Name, h)
	}
	return secret, nil
}

// ParamNames returns the set of supported input parameters when using the parameter injection generator pattern
func (s SecretForDockerRegistryGeneratorV1) ParamNames() []GeneratorParam {
	return []GeneratorParam{
		{"name", true},
		{"from-file", false},
		{"docker-username", true},
		{"docker-email", false},
		{"docker-password", true},
		{"docker-server", true},
		{"append-hash", false},
	}
}

// validate validates required fields are set to support structured generation
func (s SecretForDockerRegistryGeneratorV1) validate() error {
	if len(s.Name) == 0 {
		return fmt.Errorf("name must be specified")
	}

	if len(s.FileSources) == 0 {
		if len(s.Username) == 0 {
			return fmt.Errorf("username must be specified")
		}
		if len(s.Password) == 0 {
			return fmt.Errorf("password must be specified")
		}
		if len(s.Server) == 0 {
			return fmt.Errorf("server must be specified")
		}
	}
	return nil
}

// handleDockerCfgJsonContent serializes a ~/.docker/config.json file
func handleDockerCfgJsonContent(username, password, email, server string) ([]byte, error) {
	dockercfgAuth := credentialprovider.DockerConfigEntry{
		Username: username,
		Password: password,
		Email:    email,
	}

	dockerCfgJson := credentialprovider.DockerConfigJson{
		Auths: map[string]credentialprovider.DockerConfigEntry{server: dockercfgAuth},
	}

	return json.Marshal(dockerCfgJson)
}
