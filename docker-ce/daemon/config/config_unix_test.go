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
// +build !windows

package config

import (
	"io/ioutil"
	"runtime"

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
			"default-ulimits": {
				"nofile": {
					"Name": "nofile",
					"Hard": 2048,
					"Soft": 1024
				}
			},
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

	if cc.Ulimits == nil {
		t.Fatal("expected default ulimit config, got nil\n")
	} else {
		if _, OK := cc.Ulimits["nofile"]; OK {
			if cc.Ulimits["nofile"].Name != "nofile" ||
				cc.Ulimits["nofile"].Hard != 2048 ||
				cc.Ulimits["nofile"].Soft != 1024 {
				t.Fatalf("expected default ulimit name, hard and soft are nofile, 2048, 1024, got %s, %d, %d\n", cc.Ulimits["nofile"].Name, cc.Ulimits["nofile"].Hard, cc.Ulimits["nofile"].Soft)
			}
		} else {
			t.Fatal("expected default ulimit name nofile, got nil\n")
		}
	}
}

func TestDaemonConfigurationMergeShmSize(t *testing.T) {
	if runtime.GOOS == "solaris" {
		t.Skip("ShmSize not supported on Solaris\n")
	}
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil {
		t.Fatal(err)
	}

	configFile := f.Name()

	f.Write([]byte(`
		{
			"default-shm-size": "1g"
		}`))

	f.Close()

	c := &Config{}
	cc, err := MergeDaemonConfigurations(c, nil, configFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedValue := 1 * 1024 * 1024 * 1024
	if cc.ShmSize.Value() != int64(expectedValue) {
		t.Fatalf("expected default shm size %d, got %d", expectedValue, cc.ShmSize.Value())
	}
}
