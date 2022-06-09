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
package client

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type RancherBaseClientImpl struct {
	Opts    *ClientOpts
	Schemas *Schemas
	Types   map[string]Schema
}

type RancherBaseClient interface {
	Websocket(string, map[string][]string) (*websocket.Conn, *http.Response, error)
	List(string, *ListOpts, interface{}) error
	Post(string, interface{}, interface{}) error
	GetLink(Resource, string, interface{}) error
	Create(string, interface{}, interface{}) error
	Update(string, *Resource, interface{}, interface{}) error
	ById(string, string, interface{}) error
	Delete(*Resource) error
	Reload(*Resource, interface{}) error
	Action(string, string, *Resource, interface{}, interface{}) error

	doGet(string, *ListOpts, interface{}) error
	doList(string, *ListOpts, interface{}) error
	doModify(string, string, interface{}, interface{}) error
	doCreate(string, interface{}, interface{}) error
	doUpdate(string, *Resource, interface{}, interface{}) error
	doById(string, string, interface{}) error
	doResourceDelete(string, *Resource) error
	doAction(string, string, *Resource, interface{}, interface{}) error
}
