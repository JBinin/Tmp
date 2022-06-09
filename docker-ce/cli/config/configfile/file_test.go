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
package configfile

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestEncodeAuth(t *testing.T) {
	newAuthConfig := &types.AuthConfig{Username: "ken", Password: "test"}
	authStr := encodeAuth(newAuthConfig)
	decAuthConfig := &types.AuthConfig{}
	var err error
	decAuthConfig.Username, decAuthConfig.Password, err = decodeAuth(authStr)
	if err != nil {
		t.Fatal(err)
	}
	if newAuthConfig.Username != decAuthConfig.Username {
		t.Fatal("Encode Username doesn't match decoded Username")
	}
	if newAuthConfig.Password != decAuthConfig.Password {
		t.Fatal("Encode Password doesn't match decoded Password")
	}
	if authStr != "a2VuOnRlc3Q=" {
		t.Fatal("AuthString encoding isn't correct.")
	}
}
