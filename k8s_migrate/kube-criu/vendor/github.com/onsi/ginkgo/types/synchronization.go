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
package types

import (
	"encoding/json"
)

type RemoteBeforeSuiteState int

const (
	RemoteBeforeSuiteStateInvalid RemoteBeforeSuiteState = iota

	RemoteBeforeSuiteStatePending
	RemoteBeforeSuiteStatePassed
	RemoteBeforeSuiteStateFailed
	RemoteBeforeSuiteStateDisappeared
)

type RemoteBeforeSuiteData struct {
	Data  []byte
	State RemoteBeforeSuiteState
}

func (r RemoteBeforeSuiteData) ToJSON() []byte {
	data, _ := json.Marshal(r)
	return data
}

type RemoteAfterSuiteData struct {
	CanRun bool
}
