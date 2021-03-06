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
/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package routes

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"sync"

	"github.com/golang/glog"

	"k8s.io/apiserver/pkg/server/mux"
)

var (
	lock            = &sync.RWMutex{}
	registeredFlags = map[string]debugFlag{}
)

// DebugFlags adds handlers for flags under /debug/flags.
type DebugFlags struct {
}

// Install registers the APIServer's flags handler.
func (f DebugFlags) Install(c *mux.PathRecorderMux, flag string, handler func(http.ResponseWriter, *http.Request)) {
	c.UnlistedHandle("/debug/flags", http.HandlerFunc(f.Index))
	c.UnlistedHandlePrefix("/debug/flags/", http.HandlerFunc(f.Index))

	url := path.Join("/debug/flags", flag)
	c.UnlistedHandleFunc(url, handler)

	f.addFlag(flag)
}

// Index responds with the `/debug/flags` request.
// For example, "/debug/flags/v" serves the "--v" flag.
// Index responds to a request for "/debug/flags/" with an HTML page
// listing the available flags.
func (f DebugFlags) Index(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()
	if err := indexTmpl.Execute(w, registeredFlags); err != nil {
		glog.Error(err)
	}
}

var indexTmpl = template.Must(template.New("index").Parse(`<html>
<head>
<title>/debug/flags/</title>
</head>
<body>
/debug/flags/<br>
<br>
flags:<br>
<table>
{{range .}}
<tr>{{.Flag}}<br>
{{end}}
</table>
<br>
full flags configurable<br>
</body>
</html>
`))

type debugFlag struct {
	Flag string
}

func (f DebugFlags) addFlag(flag string) {
	lock.Lock()
	defer lock.Unlock()
	registeredFlags[flag] = debugFlag{flag}
}

// StringFlagSetterFunc is a func used for setting string type flag.
type StringFlagSetterFunc func(string) (string, error)

// StringFlagPutHandler wraps an http Handler to set string type flag.
func StringFlagPutHandler(setter StringFlagSetterFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.Method == "PUT":
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				writePlainText(http.StatusBadRequest, "error reading request body: "+err.Error(), w)
				return
			}
			defer req.Body.Close()
			response, err := setter(string(body))
			if err != nil {
				writePlainText(http.StatusBadRequest, err.Error(), w)
				return
			}
			writePlainText(http.StatusOK, response, w)
			return
		default:
			writePlainText(http.StatusNotAcceptable, "unsupported http method", w)
			return
		}
	})
}

// writePlainText renders a simple string response.
func writePlainText(statusCode int, text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, text)
}
