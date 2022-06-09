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
package equality

import (
	"reflect"

	"github.com/docker/swarmkit/api"
)

// TasksEqualStable returns true if the tasks are functionally equal, ignoring status,
// version and other superfluous fields.
//
// This used to decide whether or not to propagate a task update to a controller.
func TasksEqualStable(a, b *api.Task) bool {
	// shallow copy
	copyA, copyB := *a, *b

	copyA.Status, copyB.Status = api.TaskStatus{}, api.TaskStatus{}
	copyA.Meta, copyB.Meta = api.Meta{}, api.Meta{}

	return reflect.DeepEqual(&copyA, &copyB)
}

// TaskStatusesEqualStable compares the task status excluding timestamp fields.
func TaskStatusesEqualStable(a, b *api.TaskStatus) bool {
	copyA, copyB := *a, *b

	copyA.Timestamp, copyB.Timestamp = nil, nil
	return reflect.DeepEqual(&copyA, &copyB)
}
