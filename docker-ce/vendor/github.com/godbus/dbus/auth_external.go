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
package dbus

import (
	"encoding/hex"
)

// AuthExternal returns an Auth that authenticates as the given user with the
// EXTERNAL mechanism.
func AuthExternal(user string) Auth {
	return authExternal{user}
}

// AuthExternal implements the EXTERNAL authentication mechanism.
type authExternal struct {
	user string
}

func (a authExternal) FirstData() ([]byte, []byte, AuthStatus) {
	b := make([]byte, 2*len(a.user))
	hex.Encode(b, []byte(a.user))
	return []byte("EXTERNAL"), b, AuthOk
}

func (a authExternal) HandleData(b []byte) ([]byte, AuthStatus) {
	return nil, AuthError
}
