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
Copyright (c) 2017 VMware, Inc. All Rights Reserved.

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

package simulator

import (
	"sync"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

var recentTaskMax = 200 // the VC limit

type TaskManager struct {
	mo.TaskManager
	sync.Mutex
}

func NewTaskManager(ref types.ManagedObjectReference) object.Reference {
	s := &TaskManager{}
	s.Self = ref
	Map.AddHandler(s)
	return s
}

func (m *TaskManager) PutObject(obj mo.Reference) {
	ref := obj.Reference()
	if ref.Type != "Task" {
		return
	}

	m.Lock()
	recent := append(m.RecentTask, ref)
	if len(recent) > recentTaskMax {
		recent = recent[1:]
	}

	Map.Update(m, []types.PropertyChange{{Name: "recentTask", Val: recent}})
	m.Unlock()
}

func (*TaskManager) RemoveObject(types.ManagedObjectReference) {}

func (*TaskManager) UpdateObject(mo.Reference, []types.PropertyChange) {}
