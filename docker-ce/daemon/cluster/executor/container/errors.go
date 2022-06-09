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

import (
	"errors"
)

var (
	// ErrImageRequired returned if a task is missing the image definition.
	ErrImageRequired = errors.New("dockerexec: image required")

	// ErrContainerDestroyed returned when a container is prematurely destroyed
	// during a wait call.
	ErrContainerDestroyed = errors.New("dockerexec: container destroyed")

	// ErrContainerUnhealthy returned if controller detects the health check failure
	ErrContainerUnhealthy = errors.New("dockerexec: unhealthy container")
)
