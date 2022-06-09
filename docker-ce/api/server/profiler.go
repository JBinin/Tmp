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
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

const debugPathPrefix = "/debug/"

func profilerSetup(mainRouter *mux.Router) {
	var r = mainRouter.PathPrefix(debugPathPrefix).Subrouter()
	r.HandleFunc("/vars", expVars)
	r.HandleFunc("/pprof/", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.HandleFunc("/pprof/block", pprof.Handler("block").ServeHTTP)
	r.HandleFunc("/pprof/heap", pprof.Handler("heap").ServeHTTP)
	r.HandleFunc("/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.HandleFunc("/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

// Replicated from expvar.go as not public.
func expVars(w http.ResponseWriter, r *http.Request) {
	first := true
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, "{")
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintln(w, ",")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintln(w, "\n}")
}
