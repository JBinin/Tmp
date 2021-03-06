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

package framework

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	e2eframework "k8s.io/kubernetes/test/e2e/framework"
	testutils "k8s.io/kubernetes/test/utils"

	"github.com/golang/glog"
)

const (
	retries = 5
)

type IntegrationTestNodePreparer struct {
	client          clientset.Interface
	countToStrategy []testutils.CountToStrategy
	nodeNamePrefix  string
}

func NewIntegrationTestNodePreparer(client clientset.Interface, countToStrategy []testutils.CountToStrategy, nodeNamePrefix string) testutils.TestNodePreparer {
	return &IntegrationTestNodePreparer{
		client:          client,
		countToStrategy: countToStrategy,
		nodeNamePrefix:  nodeNamePrefix,
	}
}

func (p *IntegrationTestNodePreparer) PrepareNodes() error {
	numNodes := 0
	for _, v := range p.countToStrategy {
		numNodes += v.Count
	}

	glog.Infof("Making %d nodes", numNodes)
	baseNode := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: p.nodeNamePrefix,
		},
		Status: v1.NodeStatus{
			Capacity: v1.ResourceList{
				v1.ResourcePods:   *resource.NewQuantity(110, resource.DecimalSI),
				v1.ResourceCPU:    resource.MustParse("4"),
				v1.ResourceMemory: resource.MustParse("32Gi"),
			},
			Phase: v1.NodeRunning,
			Conditions: []v1.NodeCondition{
				{Type: v1.NodeReady, Status: v1.ConditionTrue},
			},
		},
	}
	for i := 0; i < numNodes; i++ {
		var err error
		for retry := 0; retry < retries; retry++ {
			_, err = p.client.CoreV1().Nodes().Create(baseNode)
			if err == nil || !testutils.IsRetryableAPIError(err) {
				break
			}
		}
		if err != nil {
			glog.Fatalf("Error creating node: %v", err)
		}
	}

	nodes := e2eframework.GetReadySchedulableNodesOrDie(p.client)
	index := 0
	sum := 0
	for _, v := range p.countToStrategy {
		sum += v.Count
		for ; index < sum; index++ {
			if err := testutils.DoPrepareNode(p.client, &nodes.Items[index], v.Strategy); err != nil {
				glog.Errorf("Aborting node preparation: %v", err)
				return err
			}
		}
	}
	return nil
}

func (p *IntegrationTestNodePreparer) CleanupNodes() error {
	nodes := e2eframework.GetReadySchedulableNodesOrDie(p.client)
	for i := range nodes.Items {
		if err := p.client.CoreV1().Nodes().Delete(nodes.Items[i].Name, &metav1.DeleteOptions{}); err != nil {
			glog.Errorf("Error while deleting Node: %v", err)
		}
	}
	return nil
}
