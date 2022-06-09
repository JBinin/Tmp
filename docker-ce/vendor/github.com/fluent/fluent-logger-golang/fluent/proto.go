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
//go:generate msgp

package fluent

//msgp:tuple Entry
type Entry struct {
	Time   int64       `msg:"time"`
	Record interface{} `msg:"record"`
}

//msgp:tuple Forward
type Forward struct {
	Tag     string      `msg:"tag"`
	Entries []Entry     `msg:"entries"`
	Option  interface{} `msg:"option"`
}

//msgp:tuple Message
type Message struct {
	Tag    string      `msg:"tag"`
	Time   int64       `msg:"time"`
	Record interface{} `msg:"record"`
	Option interface{} `msg:"option"`
}
