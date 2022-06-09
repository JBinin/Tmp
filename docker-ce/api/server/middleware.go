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
package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/middleware"
)

// handlerWithGlobalMiddlewares wraps the handler function for a request with
// the server's global middlewares. The order of the middlewares is backwards,
// meaning that the first in the list will be evaluated last.
func (s *Server) handlerWithGlobalMiddlewares(handler httputils.APIFunc) httputils.APIFunc {
	next := handler

	for _, m := range s.middlewares {
		next = m.WrapHandler(next)
	}

	if s.cfg.Logging && logrus.GetLevel() == logrus.DebugLevel {
		next = middleware.DebugRequestMiddleware(next)
	}

	return next
}
