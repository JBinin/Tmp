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
// +build !windows,!solaris

// TODO: Create new file for Solaris which tests config parameters
// as described in daemon/config_solaris.go

package main

import (
	"testing"

	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/docker/docker/pkg/testutil/tempfile"
)

func TestLoadDaemonCliConfigWithDaemonFlags(t *testing.T) {
	content := `{"log-opts": {"max-size": "1k"}}`
	tempFile := tempfile.NewTempFile(t, "config", content)
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Name())
	opts.common.Debug = true
	opts.common.LogLevel = "info"
	assert.NilError(t, opts.flags.Set("selinux-enabled", "true"))

	loadedConfig, err := loadDaemonCliConfig(opts)
	assert.NilError(t, err)
	assert.NotNil(t, loadedConfig)

	assert.Equal(t, loadedConfig.Debug, true)
	assert.Equal(t, loadedConfig.LogLevel, "info")
	assert.Equal(t, loadedConfig.EnableSelinuxSupport, true)
	assert.Equal(t, loadedConfig.LogConfig.Type, "json-file")
	assert.Equal(t, loadedConfig.LogConfig.Config["max-size"], "1k")
}

func TestLoadDaemonConfigWithNetwork(t *testing.T) {
	content := `{"bip": "127.0.0.2", "ip": "127.0.0.1"}`
	tempFile := tempfile.NewTempFile(t, "config", content)
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Name())
	loadedConfig, err := loadDaemonCliConfig(opts)
	assert.NilError(t, err)
	assert.NotNil(t, loadedConfig)

	assert.Equal(t, loadedConfig.IP, "127.0.0.2")
	assert.Equal(t, loadedConfig.DefaultIP.String(), "127.0.0.1")
}

func TestLoadDaemonConfigWithMapOptions(t *testing.T) {
	content := `{
		"cluster-store-opts": {"kv.cacertfile": "/var/lib/docker/discovery_certs/ca.pem"},
		"log-opts": {"tag": "test"}
}`
	tempFile := tempfile.NewTempFile(t, "config", content)
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Name())
	loadedConfig, err := loadDaemonCliConfig(opts)
	assert.NilError(t, err)
	assert.NotNil(t, loadedConfig)
	assert.NotNil(t, loadedConfig.ClusterOpts)

	expectedPath := "/var/lib/docker/discovery_certs/ca.pem"
	assert.Equal(t, loadedConfig.ClusterOpts["kv.cacertfile"], expectedPath)
	assert.NotNil(t, loadedConfig.LogConfig.Config)
	assert.Equal(t, loadedConfig.LogConfig.Config["tag"], "test")
}

func TestLoadDaemonConfigWithTrueDefaultValues(t *testing.T) {
	content := `{ "userland-proxy": false }`
	tempFile := tempfile.NewTempFile(t, "config", content)
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Name())
	loadedConfig, err := loadDaemonCliConfig(opts)
	assert.NilError(t, err)
	assert.NotNil(t, loadedConfig)
	assert.NotNil(t, loadedConfig.ClusterOpts)

	assert.Equal(t, loadedConfig.EnableUserlandProxy, false)

	// make sure reloading doesn't generate configuration
	// conflicts after normalizing boolean values.
	reload := func(reloadedConfig *config.Config) {
		assert.Equal(t, reloadedConfig.EnableUserlandProxy, false)
	}
	assert.NilError(t, config.Reload(opts.configFile, opts.flags, reload))
}

func TestLoadDaemonConfigWithTrueDefaultValuesLeaveDefaults(t *testing.T) {
	tempFile := tempfile.NewTempFile(t, "config", `{}`)
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Name())
	loadedConfig, err := loadDaemonCliConfig(opts)
	assert.NilError(t, err)
	assert.NotNil(t, loadedConfig)
	assert.NotNil(t, loadedConfig.ClusterOpts)

	assert.Equal(t, loadedConfig.EnableUserlandProxy, true)
}

func TestLoadDaemonConfigWithLegacyRegistryOptions(t *testing.T) {
	content := `{"disable-legacy-registry": true}`
	tempFile := tempfile.NewTempFile(t, "config", content)
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Name())
	loadedConfig, err := loadDaemonCliConfig(opts)
	assert.NilError(t, err)
	assert.NotNil(t, loadedConfig)
	assert.Equal(t, loadedConfig.V2Only, true)
}
