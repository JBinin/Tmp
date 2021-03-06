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
package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/pkg/errors"
)

// Swarm is a test daemon with helpers for participating in a swarm.
type Swarm struct {
	*Daemon
	swarm.Info
	Port       int
	ListenAddr string
}

// Init initializes a new swarm cluster.
func (d *Swarm) Init(req swarm.InitRequest) error {
	if req.ListenAddr == "" {
		req.ListenAddr = d.ListenAddr
	}
	status, out, err := d.SockRequest("POST", "/swarm/init", req)
	if status != http.StatusOK {
		return fmt.Errorf("initializing swarm: invalid statuscode %v, %q", status, out)
	}
	if err != nil {
		return fmt.Errorf("initializing swarm: %v", err)
	}
	info, err := d.SwarmInfo()
	if err != nil {
		return err
	}
	d.Info = info
	return nil
}

// Join joins a daemon to an existing cluster.
func (d *Swarm) Join(req swarm.JoinRequest) error {
	if req.ListenAddr == "" {
		req.ListenAddr = d.ListenAddr
	}
	status, out, err := d.SockRequest("POST", "/swarm/join", req)
	if status != http.StatusOK {
		return fmt.Errorf("joining swarm: invalid statuscode %v, %q", status, out)
	}
	if err != nil {
		return fmt.Errorf("joining swarm: %v", err)
	}
	info, err := d.SwarmInfo()
	if err != nil {
		return err
	}
	d.Info = info
	return nil
}

// Leave forces daemon to leave current cluster.
func (d *Swarm) Leave(force bool) error {
	url := "/swarm/leave"
	if force {
		url += "?force=1"
	}
	status, out, err := d.SockRequest("POST", url, nil)
	if status != http.StatusOK {
		return fmt.Errorf("leaving swarm: invalid statuscode %v, %q", status, out)
	}
	if err != nil {
		err = fmt.Errorf("leaving swarm: %v", err)
	}
	return err
}

// SwarmInfo returns the swarm information of the daemon
func (d *Swarm) SwarmInfo() (swarm.Info, error) {
	var info struct {
		Swarm swarm.Info
	}
	status, dt, err := d.SockRequest("GET", "/info", nil)
	if status != http.StatusOK {
		return info.Swarm, fmt.Errorf("get swarm info: invalid statuscode %v", status)
	}
	if err != nil {
		return info.Swarm, fmt.Errorf("get swarm info: %v", err)
	}
	if err := json.Unmarshal(dt, &info); err != nil {
		return info.Swarm, err
	}
	return info.Swarm, nil
}

// Unlock tries to unlock a locked swarm
func (d *Swarm) Unlock(req swarm.UnlockRequest) error {
	status, out, err := d.SockRequest("POST", "/swarm/unlock", req)
	if status != http.StatusOK {
		return fmt.Errorf("unlocking swarm: invalid statuscode %v, %q", status, out)
	}
	if err != nil {
		err = errors.Wrap(err, "unlocking swarm")
	}
	return err
}

// ServiceConstructor defines a swarm service constructor function
type ServiceConstructor func(*swarm.Service)

// NodeConstructor defines a swarm node constructor
type NodeConstructor func(*swarm.Node)

// SpecConstructor defines a swarm spec constructor
type SpecConstructor func(*swarm.Spec)

// CreateService creates a swarm service given the specified service constructor
func (d *Swarm) CreateService(c *check.C, f ...ServiceConstructor) string {
	var service swarm.Service
	for _, fn := range f {
		fn(&service)
	}
	status, out, err := d.SockRequest("POST", "/services/create", service.Spec)

	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusCreated, check.Commentf("output: %q", string(out)))

	var scr types.ServiceCreateResponse
	c.Assert(json.Unmarshal(out, &scr), checker.IsNil)
	return scr.ID
}

// GetService returns the swarm service corresponding to the specified id
func (d *Swarm) GetService(c *check.C, id string) *swarm.Service {
	var service swarm.Service
	status, out, err := d.SockRequest("GET", "/services/"+id, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &service), checker.IsNil)
	return &service
}

// GetServiceTasks returns the swarm tasks for the specified service
func (d *Swarm) GetServiceTasks(c *check.C, service string) []swarm.Task {
	var tasks []swarm.Task

	filterArgs := filters.NewArgs()
	filterArgs.Add("desired-state", "running")
	filterArgs.Add("service", service)
	filters, err := filters.ToParam(filterArgs)
	c.Assert(err, checker.IsNil)

	status, out, err := d.SockRequest("GET", "/tasks?filters="+filters, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &tasks), checker.IsNil)
	return tasks
}

// CheckServiceTasksInState returns the number of tasks with a matching state,
// and optional message substring.
func (d *Swarm) CheckServiceTasksInState(service string, state swarm.TaskState, message string) func(*check.C) (interface{}, check.CommentInterface) {
	return func(c *check.C) (interface{}, check.CommentInterface) {
		tasks := d.GetServiceTasks(c, service)
		var count int
		for _, task := range tasks {
			if task.Status.State == state {
				if message == "" || strings.Contains(task.Status.Message, message) {
					count++
				}
			}
		}
		return count, nil
	}
}

// CheckServiceRunningTasks returns the number of running tasks for the specified service
func (d *Swarm) CheckServiceRunningTasks(service string) func(*check.C) (interface{}, check.CommentInterface) {
	return d.CheckServiceTasksInState(service, swarm.TaskStateRunning, "")
}

// CheckServiceUpdateState returns the current update state for the specified service
func (d *Swarm) CheckServiceUpdateState(service string) func(*check.C) (interface{}, check.CommentInterface) {
	return func(c *check.C) (interface{}, check.CommentInterface) {
		service := d.GetService(c, service)
		if service.UpdateStatus == nil {
			return "", nil
		}
		return service.UpdateStatus.State, nil
	}
}

// CheckServiceTasks returns the number of tasks for the specified service
func (d *Swarm) CheckServiceTasks(service string) func(*check.C) (interface{}, check.CommentInterface) {
	return func(c *check.C) (interface{}, check.CommentInterface) {
		tasks := d.GetServiceTasks(c, service)
		return len(tasks), nil
	}
}

// CheckRunningTaskImages returns the number of different images attached to a running task
func (d *Swarm) CheckRunningTaskImages(c *check.C) (interface{}, check.CommentInterface) {
	var tasks []swarm.Task

	filterArgs := filters.NewArgs()
	filterArgs.Add("desired-state", "running")
	filters, err := filters.ToParam(filterArgs)
	c.Assert(err, checker.IsNil)

	status, out, err := d.SockRequest("GET", "/tasks?filters="+filters, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &tasks), checker.IsNil)

	result := make(map[string]int)
	for _, task := range tasks {
		if task.Status.State == swarm.TaskStateRunning {
			result[task.Spec.ContainerSpec.Image]++
		}
	}
	return result, nil
}

// CheckNodeReadyCount returns the number of ready node on the swarm
func (d *Swarm) CheckNodeReadyCount(c *check.C) (interface{}, check.CommentInterface) {
	nodes := d.ListNodes(c)
	var readyCount int
	for _, node := range nodes {
		if node.Status.State == swarm.NodeStateReady {
			readyCount++
		}
	}
	return readyCount, nil
}

// GetTask returns the swarm task identified by the specified id
func (d *Swarm) GetTask(c *check.C, id string) swarm.Task {
	var task swarm.Task

	status, out, err := d.SockRequest("GET", "/tasks/"+id, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &task), checker.IsNil)
	return task
}

// UpdateService updates a swarm service with the specified service constructor
func (d *Swarm) UpdateService(c *check.C, service *swarm.Service, f ...ServiceConstructor) {
	for _, fn := range f {
		fn(service)
	}
	url := fmt.Sprintf("/services/%s/update?version=%d", service.ID, service.Version.Index)
	status, out, err := d.SockRequest("POST", url, service.Spec)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
}

// RemoveService removes the specified service
func (d *Swarm) RemoveService(c *check.C, id string) {
	status, out, err := d.SockRequest("DELETE", "/services/"+id, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
}

// GetNode returns a swarm node identified by the specified id
func (d *Swarm) GetNode(c *check.C, id string) *swarm.Node {
	var node swarm.Node
	status, out, err := d.SockRequest("GET", "/nodes/"+id, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &node), checker.IsNil)
	c.Assert(node.ID, checker.Equals, id)
	return &node
}

// RemoveNode removes the specified node
func (d *Swarm) RemoveNode(c *check.C, id string, force bool) {
	url := "/nodes/" + id
	if force {
		url += "?force=1"
	}

	status, out, err := d.SockRequest("DELETE", url, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
}

// UpdateNode updates a swarm node with the specified node constructor
func (d *Swarm) UpdateNode(c *check.C, id string, f ...NodeConstructor) {
	for i := 0; ; i++ {
		node := d.GetNode(c, id)
		for _, fn := range f {
			fn(node)
		}
		url := fmt.Sprintf("/nodes/%s/update?version=%d", node.ID, node.Version.Index)
		status, out, err := d.SockRequest("POST", url, node.Spec)
		if i < 10 && strings.Contains(string(out), "update out of sequence") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		c.Assert(err, checker.IsNil, check.Commentf(string(out)))
		c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
		return
	}
}

// ListNodes returns the list of the current swarm nodes
func (d *Swarm) ListNodes(c *check.C) []swarm.Node {
	status, out, err := d.SockRequest("GET", "/nodes", nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))

	nodes := []swarm.Node{}
	c.Assert(json.Unmarshal(out, &nodes), checker.IsNil)
	return nodes
}

// ListServices return the list of the current swarm services
func (d *Swarm) ListServices(c *check.C) []swarm.Service {
	status, out, err := d.SockRequest("GET", "/services", nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))

	services := []swarm.Service{}
	c.Assert(json.Unmarshal(out, &services), checker.IsNil)
	return services
}

// CreateSecret creates a secret given the specified spec
func (d *Swarm) CreateSecret(c *check.C, secretSpec swarm.SecretSpec) string {
	status, out, err := d.SockRequest("POST", "/secrets/create", secretSpec)

	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusCreated, check.Commentf("output: %q", string(out)))

	var scr types.SecretCreateResponse
	c.Assert(json.Unmarshal(out, &scr), checker.IsNil)
	return scr.ID
}

// ListSecrets returns the list of the current swarm secrets
func (d *Swarm) ListSecrets(c *check.C) []swarm.Secret {
	status, out, err := d.SockRequest("GET", "/secrets", nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))

	secrets := []swarm.Secret{}
	c.Assert(json.Unmarshal(out, &secrets), checker.IsNil)
	return secrets
}

// GetSecret returns a swarm secret identified by the specified id
func (d *Swarm) GetSecret(c *check.C, id string) *swarm.Secret {
	var secret swarm.Secret
	status, out, err := d.SockRequest("GET", "/secrets/"+id, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &secret), checker.IsNil)
	return &secret
}

// DeleteSecret removes the swarm secret identified by the specified id
func (d *Swarm) DeleteSecret(c *check.C, id string) {
	status, out, err := d.SockRequest("DELETE", "/secrets/"+id, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusNoContent, check.Commentf("output: %q", string(out)))
}

// GetSwarm return the current swarm object
func (d *Swarm) GetSwarm(c *check.C) swarm.Swarm {
	var sw swarm.Swarm
	status, out, err := d.SockRequest("GET", "/swarm", nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &sw), checker.IsNil)
	return sw
}

// UpdateSwarm updates the current swarm object with the specified spec constructors
func (d *Swarm) UpdateSwarm(c *check.C, f ...SpecConstructor) {
	sw := d.GetSwarm(c)
	for _, fn := range f {
		fn(&sw.Spec)
	}
	url := fmt.Sprintf("/swarm/update?version=%d", sw.Version.Index)
	status, out, err := d.SockRequest("POST", url, sw.Spec)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
}

// RotateTokens update the swarm to rotate tokens
func (d *Swarm) RotateTokens(c *check.C) {
	var sw swarm.Swarm
	status, out, err := d.SockRequest("GET", "/swarm", nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &sw), checker.IsNil)

	url := fmt.Sprintf("/swarm/update?version=%d&rotateWorkerToken=true&rotateManagerToken=true", sw.Version.Index)
	status, out, err = d.SockRequest("POST", url, sw.Spec)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
}

// JoinTokens returns the current swarm join tokens
func (d *Swarm) JoinTokens(c *check.C) swarm.JoinTokens {
	var sw swarm.Swarm
	status, out, err := d.SockRequest("GET", "/swarm", nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf("output: %q", string(out)))
	c.Assert(json.Unmarshal(out, &sw), checker.IsNil)
	return sw.JoinTokens
}

// CheckLocalNodeState returns the current swarm node state
func (d *Swarm) CheckLocalNodeState(c *check.C) (interface{}, check.CommentInterface) {
	info, err := d.SwarmInfo()
	c.Assert(err, checker.IsNil)
	return info.LocalNodeState, nil
}

// CheckControlAvailable returns the current swarm control available
func (d *Swarm) CheckControlAvailable(c *check.C) (interface{}, check.CommentInterface) {
	info, err := d.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
	return info.ControlAvailable, nil
}

// CheckLeader returns whether there is a leader on the swarm or not
func (d *Swarm) CheckLeader(c *check.C) (interface{}, check.CommentInterface) {
	errList := check.Commentf("could not get node list")
	status, out, err := d.SockRequest("GET", "/nodes", nil)
	if err != nil {
		return err, errList
	}
	if status != http.StatusOK {
		return fmt.Errorf("expected http status OK, got: %d", status), errList
	}

	var ls []swarm.Node
	if err := json.Unmarshal(out, &ls); err != nil {
		return err, errList
	}

	for _, node := range ls {
		if node.ManagerStatus != nil && node.ManagerStatus.Leader {
			return nil, nil
		}
	}
	return fmt.Errorf("no leader"), check.Commentf("could not find leader")
}

// CmdRetryOutOfSequence tries the specified command against the current daemon for 10 times
func (d *Swarm) CmdRetryOutOfSequence(args ...string) (string, error) {
	for i := 0; ; i++ {
		out, err := d.Cmd(args...)
		if err != nil {
			if strings.Contains(out, "update out of sequence") {
				if i < 10 {
					continue
				}
			}
		}
		return out, err
	}
}
