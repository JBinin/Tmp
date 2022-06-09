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
package probing

import (
	"encoding/json"
	"net/http"
	"time"
)

func NewHandler() http.Handler {
	return &httpHealth{}
}

type httpHealth struct {
}

type Health struct {
	OK  bool
	Now time.Time
}

func (h *httpHealth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	health := Health{OK: true, Now: time.Now()}
	e := json.NewEncoder(w)
	e.Encode(health)
}
