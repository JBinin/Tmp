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

package app

import (
	apiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
	"k8s.io/kubernetes/pkg/controller"
)

// Config is the main context object for the cloud controller manager.
type Config struct {
	ComponentConfig componentconfig.CloudControllerManagerConfiguration

	SecureServing *apiserver.SecureServingInfo
	// TODO: remove deprecated insecure serving
	InsecureServing *apiserver.DeprecatedInsecureServingInfo
	Authentication  apiserver.AuthenticationInfo
	Authorization   apiserver.AuthorizationInfo

	// the general kube client
	Client *clientset.Clientset

	// the client only used for leader election
	LeaderElectionClient *clientset.Clientset

	// the rest config for the master
	Kubeconfig *restclient.Config

	// the event sink
	EventRecorder record.EventRecorder

	// ClientBuilder will provide a client for this controller to use
	ClientBuilder controller.ControllerClientBuilder

	// VersionedClient will provide a client for informers
	VersionedClient clientset.Interface

	// SharedInformers gives access to informers for the controller.
	SharedInformers informers.SharedInformerFactory
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	cc := completedConfig{c}
	return &CompletedConfig{&cc}
}
