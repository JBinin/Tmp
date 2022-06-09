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
package plugin

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/plugin/v2"
)

func TestFilterByCapNeg(t *testing.T) {
	p := v2.Plugin{PluginObj: types.Plugin{Name: "test:latest"}}
	iType := types.PluginInterfaceType{"volumedriver", "docker", "1.0"}
	i := types.PluginConfigInterface{"plugins.sock", []types.PluginInterfaceType{iType}}
	p.PluginObj.Config.Interface = i

	_, err := p.FilterByCap("foobar")
	if err == nil {
		t.Fatalf("expected inadequate error, got %v", err)
	}
}

func TestFilterByCapPos(t *testing.T) {
	p := v2.Plugin{PluginObj: types.Plugin{Name: "test:latest"}}

	iType := types.PluginInterfaceType{"volumedriver", "docker", "1.0"}
	i := types.PluginConfigInterface{"plugins.sock", []types.PluginInterfaceType{iType}}
	p.PluginObj.Config.Interface = i

	_, err := p.FilterByCap("volumedriver")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
