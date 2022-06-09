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
	"strings"

	"github.com/docker/docker/container"
)

// excludeByIsolation is a platform specific helper function to support PS
// filtering by Isolation. This is a Windows-only concept, so is a no-op on Unix.
func excludeByIsolation(container *container.Container, ctx *listContext) iterationAction {
	i := strings.ToLower(string(container.HostConfig.Isolation))
	if i == "" {
		i = "default"
	}
	if !ctx.filters.Match("isolation", i) {
		return excludeContainer
	}
	return includeContainer
}
