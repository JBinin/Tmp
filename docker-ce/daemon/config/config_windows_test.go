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
// +build windows

package config

import (
	"io/ioutil"
	"testing"
)

func TestDaemonConfigurationMerge(t *testing.T) {
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil {
		t.Fatal(err)
	}

	configFile := f.Name()

	f.Write([]byte(`
		{
			"debug": true,
			"log-opts": {
				"tag": "test_tag"
			}
		}`))

	f.Close()

	c := &Config{
		CommonConfig: CommonConfig{
			AutoRestart: true,
			LogConfig: LogConfig{
				Type:   "syslog",
				Config: map[string]string{"tag": "test"},
			},
		},
	}

	cc, err := MergeDaemonConfigurations(c, nil, configFile)
	if err != nil {
		t.Fatal(err)
	}
	if !cc.Debug {
		t.Fatalf("expected %v, got %v\n", true, cc.Debug)
	}
	if !cc.AutoRestart {
		t.Fatalf("expected %v, got %v\n", true, cc.AutoRestart)
	}
	if cc.LogConfig.Type != "syslog" {
		t.Fatalf("expected syslog config, got %q\n", cc.LogConfig)
	}

	if configValue, OK := cc.LogConfig.Config["tag"]; !OK {
		t.Fatal("expected syslog config attributes, got nil\n")
	} else {
		if configValue != "test_tag" {
			t.Fatalf("expected syslog config attributes 'tag=test_tag', got 'tag=%s'\n", configValue)
		}
	}
}
