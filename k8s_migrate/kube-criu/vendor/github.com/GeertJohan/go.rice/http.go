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
package rice

import (
	"net/http"
)

// HTTPBox implements http.FileSystem which allows the use of Box with a http.FileServer.
//   e.g.: http.Handle("/", http.FileServer(rice.MustFindBox("http-files").HTTPBox()))
type HTTPBox struct {
	*Box
}

// HTTPBox creates a new HTTPBox from an existing Box
func (b *Box) HTTPBox() *HTTPBox {
	return &HTTPBox{b}
}

// Open returns a File using the http.File interface
func (hb *HTTPBox) Open(name string) (http.File, error) {
	return hb.Box.Open(name)
}
