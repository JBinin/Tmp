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
// Package naming centralizes the naming of SwarmKit objects.
package naming

import (
	"fmt"

	"github.com/docker/swarmkit/api"
)

// Task returns the task name from Annotations.Name,
// and, in case Annotations.Name is missing, fallback
// to construct the name from other information.
func Task(t *api.Task) string {
	if t.Annotations.Name != "" {
		// if set, use the container Annotations.Name field, set in the orchestrator.
		return t.Annotations.Name
	}

	slot := fmt.Sprint(t.Slot)
	if slot == "" || t.Slot == 0 {
		// when no slot id is assigned, we assume that this is node-bound task.
		slot = t.NodeID
	}

	// fallback to service.instance.id.
	return fmt.Sprintf("%s.%s.%s", t.ServiceAnnotations.Name, slot, t.ID)
}

// TODO(stevvooe): Consolidate "Hostname" style validation here.
