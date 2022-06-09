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
package agent

import (
	"errors"
	"fmt"
)

var (
	// ErrClosed is returned when an operation fails because the resource is closed.
	ErrClosed = errors.New("agent: closed")

	errNodeNotRegistered = fmt.Errorf("node not registered")

	errAgentStarted    = errors.New("agent: already started")
	errAgentNotStarted = errors.New("agent: not started")

	errTaskNoContoller          = errors.New("agent: no task controller")
	errTaskNotAssigned          = errors.New("agent: task not assigned")
	errTaskStatusUpdateNoChange = errors.New("agent: no change in task status")
	errTaskUnknown              = errors.New("agent: task unknown")

	errTaskInvalid = errors.New("task: invalid")
)
