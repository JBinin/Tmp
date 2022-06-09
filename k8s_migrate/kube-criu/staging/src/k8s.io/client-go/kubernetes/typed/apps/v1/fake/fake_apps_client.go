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
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAppsV1 struct {
	*testing.Fake
}

func (c *FakeAppsV1) ControllerRevisions(namespace string) v1.ControllerRevisionInterface {
	return &FakeControllerRevisions{c, namespace}
}

func (c *FakeAppsV1) DaemonSets(namespace string) v1.DaemonSetInterface {
	return &FakeDaemonSets{c, namespace}
}

func (c *FakeAppsV1) Deployments(namespace string) v1.DeploymentInterface {
	return &FakeDeployments{c, namespace}
}

func (c *FakeAppsV1) ReplicaSets(namespace string) v1.ReplicaSetInterface {
	return &FakeReplicaSets{c, namespace}
}

func (c *FakeAppsV1) StatefulSets(namespace string) v1.StatefulSetInterface {
	return &FakeStatefulSets{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAppsV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
