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
// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package accelerators

import info "github.com/google/cadvisor/info/v1"

// This is supposed to store global state about an accelerator metrics collector.
// cadvisor manager will call Setup() when it starts and Destroy() when it stops.
// For each container detected by the cadvisor manager, it will call
// GetCollector() with the devices cgroup path for that container.
// GetCollector() is supposed to return an object that can update
// accelerator stats for that container.
type AcceleratorManager interface {
	Setup()
	Destroy()
	GetCollector(deviceCgroup string) (AcceleratorCollector, error)
}

type AcceleratorCollector interface {
	UpdateStats(*info.ContainerStats) error
}