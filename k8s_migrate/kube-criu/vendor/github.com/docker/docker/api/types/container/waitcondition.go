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
package container

// WaitCondition is a type used to specify a container state for which
// to wait.
type WaitCondition string

// Possible WaitCondition Values.
//
// WaitConditionNotRunning (default) is used to wait for any of the non-running
// states: "created", "exited", "dead", "removing", or "removed".
//
// WaitConditionNextExit is used to wait for the next time the state changes
// to a non-running state. If the state is currently "created" or "exited",
// this would cause Wait() to block until either the container runs and exits
// or is removed.
//
// WaitConditionRemoved is used to wait for the container to be removed.
const (
	WaitConditionNotRunning WaitCondition = "not-running"
	WaitConditionNextExit   WaitCondition = "next-exit"
	WaitConditionRemoved    WaitCondition = "removed"
)
