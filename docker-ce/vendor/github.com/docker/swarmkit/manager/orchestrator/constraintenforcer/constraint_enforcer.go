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
package constraintenforcer

import (
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/constraint"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
)

// ConstraintEnforcer watches for updates to nodes and shuts down tasks that no
// longer satisfy scheduling constraints or resource limits.
type ConstraintEnforcer struct {
	store    *store.MemoryStore
	stopChan chan struct{}
	doneChan chan struct{}
}

// New creates a new ConstraintEnforcer.
func New(store *store.MemoryStore) *ConstraintEnforcer {
	return &ConstraintEnforcer{
		store:    store,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

// Run is the ConstraintEnforcer's main loop.
func (ce *ConstraintEnforcer) Run() {
	defer close(ce.doneChan)

	watcher, cancelWatch := state.Watch(ce.store.WatchQueue(), state.EventUpdateNode{})
	defer cancelWatch()

	var (
		nodes []*api.Node
		err   error
	)
	ce.store.View(func(readTx store.ReadTx) {
		nodes, err = store.FindNodes(readTx, store.All)
	})
	if err != nil {
		log.L.WithError(err).Error("failed to check nodes for noncompliant tasks")
	} else {
		for _, node := range nodes {
			ce.rejectNoncompliantTasks(node)
		}
	}

	for {
		select {
		case event := <-watcher:
			node := event.(state.EventUpdateNode).Node
			ce.rejectNoncompliantTasks(node)
		case <-ce.stopChan:
			return
		}
	}
}

func (ce *ConstraintEnforcer) rejectNoncompliantTasks(node *api.Node) {
	// If the availability is "drain", the orchestrator will
	// shut down all tasks.
	// If the availability is "pause", we shouldn't touch
	// the tasks on this node.
	if node.Spec.Availability != api.NodeAvailabilityActive {
		return
	}

	var (
		tasks []*api.Task
		err   error
	)

	ce.store.View(func(tx store.ReadTx) {
		tasks, err = store.FindTasks(tx, store.ByNodeID(node.ID))
	})

	if err != nil {
		log.L.WithError(err).Errorf("failed to list tasks for node ID %s", node.ID)
	}

	var availableMemoryBytes, availableNanoCPUs int64
	if node.Description != nil && node.Description.Resources != nil {
		availableMemoryBytes = node.Description.Resources.MemoryBytes
		availableNanoCPUs = node.Description.Resources.NanoCPUs
	}

	removeTasks := make(map[string]*api.Task)

	// TODO(aaronl): The set of tasks removed will be
	// nondeterministic because it depends on the order of
	// the slice returned from FindTasks. We could do
	// a separate pass over the tasks for each type of
	// resource, and sort by the size of the reservation
	// to remove the most resource-intensive tasks.
	for _, t := range tasks {
		if t.DesiredState < api.TaskStateAssigned || t.DesiredState > api.TaskStateRunning {
			continue
		}

		// Ensure that the task still meets scheduling
		// constraints.
		if t.Spec.Placement != nil && len(t.Spec.Placement.Constraints) != 0 {
			constraints, _ := constraint.Parse(t.Spec.Placement.Constraints)
			if !constraint.NodeMatches(constraints, node) {
				removeTasks[t.ID] = t
				continue
			}
		}

		// Ensure that the task assigned to the node
		// still satisfies the resource limits.
		if t.Spec.Resources != nil && t.Spec.Resources.Reservations != nil {
			if t.Spec.Resources.Reservations.MemoryBytes > availableMemoryBytes {
				removeTasks[t.ID] = t
				continue
			}
			if t.Spec.Resources.Reservations.NanoCPUs > availableNanoCPUs {
				removeTasks[t.ID] = t
				continue
			}
			availableMemoryBytes -= t.Spec.Resources.Reservations.MemoryBytes
			availableNanoCPUs -= t.Spec.Resources.Reservations.NanoCPUs
		}
	}

	if len(removeTasks) != 0 {
		_, err := ce.store.Batch(func(batch *store.Batch) error {
			for _, t := range removeTasks {
				err := batch.Update(func(tx store.Tx) error {
					t = store.GetTask(tx, t.ID)
					if t == nil || t.DesiredState > api.TaskStateRunning {
						return nil
					}

					// We set the observed state to
					// REJECTED, rather than the desired
					// state. Desired state is owned by the
					// orchestrator, and setting it directly
					// will bypass actions such as
					// restarting the task on another node
					// (if applicable).
					t.Status.State = api.TaskStateRejected
					t.Status.Message = "assigned node no longer meets constraints"
					t.Status.Timestamp = ptypes.MustTimestampProto(time.Now())
					return store.UpdateTask(tx, t)
				})
				if err != nil {
					log.L.WithError(err).Errorf("failed to shut down task %s", t.ID)
				}
			}
			return nil
		})

		if err != nil {
			log.L.WithError(err).Errorf("failed to shut down tasks")
		}
	}
}

// Stop stops the ConstraintEnforcer and waits for the main loop to exit.
func (ce *ConstraintEnforcer) Stop() {
	close(ce.stopChan)
	<-ce.doneChan
}
