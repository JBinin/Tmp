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
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// routerSwapper is an http.Handler that allows you to swap
// mux routers.
type routerSwapper struct {
	mu     sync.Mutex
	router *mux.Router
}

// Swap changes the old router with the new one.
func (rs *routerSwapper) Swap(newRouter *mux.Router) {
	rs.mu.Lock()
	rs.router = newRouter
	rs.mu.Unlock()
}

// ServeHTTP makes the routerSwapper to implement the http.Handler interface.
func (rs *routerSwapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs.mu.Lock()
	router := rs.router
	rs.mu.Unlock()
	router.ServeHTTP(w, r)
}
