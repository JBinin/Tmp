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
package health

import (
	"encoding/json"
	"net/http"

	"github.com/cloudflare/cfssl/api"
)

// Response contains the response to the /health API
type Response struct {
	Healthy bool `json:"healthy"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) error {
	response := api.NewSuccessResponse(&Response{Healthy: true})
	return json.NewEncoder(w).Encode(response)
}

// NewHealthCheck creates a new handler to serve health checks.
func NewHealthCheck() http.Handler {
	return api.HTTPHandler{
		Handler: api.HandlerFunc(healthHandler),
		Methods: []string{"GET"},
	}
}
