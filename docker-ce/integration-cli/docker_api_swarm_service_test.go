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
// +build !windows

package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/go-check/check"
)

func setPortConfig(portConfig []swarm.PortConfig) daemon.ServiceConstructor {
	return func(s *swarm.Service) {
		if s.Spec.EndpointSpec == nil {
			s.Spec.EndpointSpec = &swarm.EndpointSpec{}
		}
		s.Spec.EndpointSpec.Ports = portConfig
	}
}

func (s *DockerSwarmSuite) TestAPIServiceUpdatePort(c *check.C) {
	d := s.AddDaemon(c, true, true)

	// Create a service with a port mapping of 8080:8081.
	portConfig := []swarm.PortConfig{{TargetPort: 8081, PublishedPort: 8080}}
	serviceID := d.CreateService(c, simpleTestService, setInstances(1), setPortConfig(portConfig))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)

	// Update the service: changed the port mapping from 8080:8081 to 8082:8083.
	updatedPortConfig := []swarm.PortConfig{{TargetPort: 8083, PublishedPort: 8082}}
	remoteService := d.GetService(c, serviceID)
	d.UpdateService(c, remoteService, setPortConfig(updatedPortConfig))

	// Inspect the service and verify port mapping.
	updatedService := d.GetService(c, serviceID)
	c.Assert(updatedService.Spec.EndpointSpec, check.NotNil)
	c.Assert(len(updatedService.Spec.EndpointSpec.Ports), check.Equals, 1)
	c.Assert(updatedService.Spec.EndpointSpec.Ports[0].TargetPort, check.Equals, uint32(8083))
	c.Assert(updatedService.Spec.EndpointSpec.Ports[0].PublishedPort, check.Equals, uint32(8082))
}

func (s *DockerSwarmSuite) TestAPISwarmServicesEmptyList(c *check.C) {
	d := s.AddDaemon(c, true, true)

	services := d.ListServices(c)
	c.Assert(services, checker.NotNil)
	c.Assert(len(services), checker.Equals, 0, check.Commentf("services: %#v", services))
}

func (s *DockerSwarmSuite) TestAPISwarmServicesCreate(c *check.C) {
	d := s.AddDaemon(c, true, true)

	instances := 2
	id := d.CreateService(c, simpleTestService, setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)

	service := d.GetService(c, id)
	instances = 5
	d.UpdateService(c, service, setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)

	d.RemoveService(c, service.ID)
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 0)
}

func (s *DockerSwarmSuite) TestAPISwarmServicesMultipleAgents(c *check.C) {
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, false)
	d3 := s.AddDaemon(c, true, false)

	time.Sleep(1 * time.Second) // make sure all daemons are ready to accept tasks

	instances := 9
	id := d1.CreateService(c, simpleTestService, setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckActiveContainerCount, checker.GreaterThan, 0)

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	// reconciliation on d2 node down
	d2.Stop(c)

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	// test downscaling
	instances = 5
	d1.UpdateService(c, d1.GetService(c, id), setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

}

func (s *DockerSwarmSuite) TestAPISwarmServicesCreateGlobal(c *check.C) {
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, false)
	d3 := s.AddDaemon(c, true, false)

	d1.CreateService(c, simpleTestService, setGlobalMode)

	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.Equals, 1)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.Equals, 1)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckActiveContainerCount, checker.Equals, 1)

	d4 := s.AddDaemon(c, true, false)
	d5 := s.AddDaemon(c, true, false)

	waitAndAssert(c, defaultReconciliationTimeout, d4.CheckActiveContainerCount, checker.Equals, 1)
	waitAndAssert(c, defaultReconciliationTimeout, d5.CheckActiveContainerCount, checker.Equals, 1)
}

func (s *DockerSwarmSuite) TestAPISwarmServicesUpdate(c *check.C) {
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ {
		daemons[i] = s.AddDaemon(c, true, i == 0)
	}
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)

	// service image at start
	image1 := "busybox:latest"
	// target image in update
	image2 := "busybox:test"

	// create a different tag
	for _, d := range daemons {
		out, err := d.Cmd("tag", image1, image2)
		c.Assert(err, checker.IsNil, check.Commentf(out))
	}

	// create service
	instances := 5
	parallelism := 2
	id := daemons[0].CreateService(c, serviceForUpdate, setInstances(instances))

	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image1: instances})

	// issue service update
	service := daemons[0].GetService(c, id)
	daemons[0].UpdateService(c, service, setImage(image2))

	// first batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image1: instances - parallelism, image2: parallelism})

	// 2nd batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image1: instances - 2*parallelism, image2: 2 * parallelism})

	// 3nd batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image2: instances})

	// Roll back to the previous version. This uses the CLI because
	// rollback is a client-side operation.
	out, err := daemons[0].Cmd("service", "update", "--rollback", id)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	// first batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image2: instances - parallelism, image1: parallelism})

	// 2nd batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image2: instances - 2*parallelism, image1: 2 * parallelism})

	// 3nd batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image1: instances})
}

func (s *DockerSwarmSuite) TestAPISwarmServicesFailedUpdate(c *check.C) {
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ {
		daemons[i] = s.AddDaemon(c, true, i == 0)
	}
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)

	// service image at start
	image1 := "busybox:latest"
	// target image in update
	image2 := "busybox:badtag"

	// create service
	instances := 5
	id := daemons[0].CreateService(c, serviceForUpdate, setInstances(instances))

	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image1: instances})

	// issue service update
	service := daemons[0].GetService(c, id)
	daemons[0].UpdateService(c, service, setImage(image2), setFailureAction(swarm.UpdateFailureActionPause), setMaxFailureRatio(0.25), setParallelism(1))

	// should update 2 tasks and then pause
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceUpdateState(id), checker.Equals, swarm.UpdateStatePaused)
	v, _ := daemons[0].CheckServiceRunningTasks(id)(c)
	c.Assert(v, checker.Equals, instances-2)

	// Roll back to the previous version. This uses the CLI because
	// rollback is a client-side operation.
	out, err := daemons[0].Cmd("service", "update", "--rollback", id)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int{image1: instances})
}

func (s *DockerSwarmSuite) TestAPISwarmServiceConstraintRole(c *check.C) {
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ {
		daemons[i] = s.AddDaemon(c, true, i == 0)
	}
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)

	// create service
	constraints := []string{"node.role==worker"}
	instances := 3
	id := daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	// validate tasks are running on worker nodes
	tasks := daemons[0].GetServiceTasks(c, id)
	for _, task := range tasks {
		node := daemons[0].GetNode(c, task.NodeID)
		c.Assert(node.Spec.Role, checker.Equals, swarm.NodeRoleWorker)
	}
	//remove service
	daemons[0].RemoveService(c, id)

	// create service
	constraints = []string{"node.role!=worker"}
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate tasks are running on manager nodes
	for _, task := range tasks {
		node := daemons[0].GetNode(c, task.NodeID)
		c.Assert(node.Spec.Role, checker.Equals, swarm.NodeRoleManager)
	}
	//remove service
	daemons[0].RemoveService(c, id)

	// create service
	constraints = []string{"node.role==nosuchrole"}
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks created
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceTasks(id), checker.Equals, instances)
	// let scheduler try
	time.Sleep(250 * time.Millisecond)
	// validate tasks are not assigned to any node
	tasks = daemons[0].GetServiceTasks(c, id)
	for _, task := range tasks {
		c.Assert(task.NodeID, checker.Equals, "")
	}
}

func (s *DockerSwarmSuite) TestAPISwarmServiceConstraintLabel(c *check.C) {
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ {
		daemons[i] = s.AddDaemon(c, true, i == 0)
	}
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)
	nodes := daemons[0].ListNodes(c)
	c.Assert(len(nodes), checker.Equals, nodeCount)

	// add labels to nodes
	daemons[0].UpdateNode(c, nodes[0].ID, func(n *swarm.Node) {
		n.Spec.Annotations.Labels = map[string]string{
			"security": "high",
		}
	})
	for i := 1; i < nodeCount; i++ {
		daemons[0].UpdateNode(c, nodes[i].ID, func(n *swarm.Node) {
			n.Spec.Annotations.Labels = map[string]string{
				"security": "low",
			}
		})
	}

	// create service
	instances := 3
	constraints := []string{"node.labels.security==high"}
	id := daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks := daemons[0].GetServiceTasks(c, id)
	// validate all tasks are running on nodes[0]
	for _, task := range tasks {
		c.Assert(task.NodeID, checker.Equals, nodes[0].ID)
	}
	//remove service
	daemons[0].RemoveService(c, id)

	// create service
	constraints = []string{"node.labels.security!=high"}
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate all tasks are NOT running on nodes[0]
	for _, task := range tasks {
		c.Assert(task.NodeID, checker.Not(checker.Equals), nodes[0].ID)
	}
	//remove service
	daemons[0].RemoveService(c, id)

	constraints = []string{"node.labels.security==medium"}
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks created
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceTasks(id), checker.Equals, instances)
	// let scheduler try
	time.Sleep(250 * time.Millisecond)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate tasks are not assigned
	for _, task := range tasks {
		c.Assert(task.NodeID, checker.Equals, "")
	}
	//remove service
	daemons[0].RemoveService(c, id)

	// multiple constraints
	constraints = []string{
		"node.labels.security==high",
		fmt.Sprintf("node.id==%s", nodes[1].ID),
	}
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks created
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceTasks(id), checker.Equals, instances)
	// let scheduler try
	time.Sleep(250 * time.Millisecond)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate tasks are not assigned
	for _, task := range tasks {
		c.Assert(task.NodeID, checker.Equals, "")
	}
	// make nodes[1] fulfills the constraints
	daemons[0].UpdateNode(c, nodes[1].ID, func(n *swarm.Node) {
		n.Spec.Annotations.Labels = map[string]string{
			"security": "high",
		}
	})
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks = daemons[0].GetServiceTasks(c, id)
	for _, task := range tasks {
		c.Assert(task.NodeID, checker.Equals, nodes[1].ID)
	}
}

func (s *DockerSwarmSuite) TestAPISwarmServicesStateReporting(c *check.C) {
	testRequires(c, SameHostDaemon)
	testRequires(c, DaemonIsLinux)

	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, true)
	d3 := s.AddDaemon(c, true, false)

	time.Sleep(1 * time.Second) // make sure all daemons are ready to accept

	instances := 9
	d1.CreateService(c, simpleTestService, setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	getContainers := func() map[string]*daemon.Swarm {
		m := make(map[string]*daemon.Swarm)
		for _, d := range []*daemon.Swarm{d1, d2, d3} {
			for _, id := range d.ActiveContainers() {
				m[id] = d
			}
		}
		return m
	}

	containers := getContainers()
	c.Assert(containers, checker.HasLen, instances)
	var toRemove string
	for i := range containers {
		toRemove = i
	}

	_, err := containers[toRemove].Cmd("stop", toRemove)
	c.Assert(err, checker.IsNil)

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	containers2 := getContainers()
	c.Assert(containers2, checker.HasLen, instances)
	for i := range containers {
		if i == toRemove {
			c.Assert(containers2[i], checker.IsNil)
		} else {
			c.Assert(containers2[i], checker.NotNil)
		}
	}

	containers = containers2
	for i := range containers {
		toRemove = i
	}

	// try with killing process outside of docker
	pidStr, err := containers[toRemove].Cmd("inspect", "-f", "{{.State.Pid}}", toRemove)
	c.Assert(err, checker.IsNil)
	pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
	c.Assert(err, checker.IsNil)
	c.Assert(syscall.Kill(pid, syscall.SIGKILL), checker.IsNil)

	time.Sleep(time.Second) // give some time to handle the signal

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	containers2 = getContainers()
	c.Assert(containers2, checker.HasLen, instances)
	for i := range containers {
		if i == toRemove {
			c.Assert(containers2[i], checker.IsNil)
		} else {
			c.Assert(containers2[i], checker.NotNil)
		}
	}
}
