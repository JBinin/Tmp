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
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"golang.org/x/net/context"
)

// getContainersByName inspects container's configuration and serializes it as json.
func (s *containerRouter) getContainersByName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	displaySize := httputils.BoolValue(r, "size")

	version := httputils.VersionFromContext(ctx)
	json, err := s.backend.ContainerInspect(vars["name"], displaySize, version)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, json)
}
