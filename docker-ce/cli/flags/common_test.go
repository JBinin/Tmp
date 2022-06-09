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
package flags

import (
	"path/filepath"
	"testing"

	cliconfig "github.com/docker/docker/cli/config"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/spf13/pflag"
)

func TestCommonOptionsInstallFlags(t *testing.T) {
	flags := pflag.NewFlagSet("testing", pflag.ContinueOnError)
	opts := NewCommonOptions()
	opts.InstallFlags(flags)

	err := flags.Parse([]string{
		"--tlscacert=\"/foo/cafile\"",
		"--tlscert=\"/foo/cert\"",
		"--tlskey=\"/foo/key\"",
	})
	assert.NilError(t, err)
	assert.Equal(t, opts.TLSOptions.CAFile, "/foo/cafile")
	assert.Equal(t, opts.TLSOptions.CertFile, "/foo/cert")
	assert.Equal(t, opts.TLSOptions.KeyFile, "/foo/key")
}

func defaultPath(filename string) string {
	return filepath.Join(cliconfig.Dir(), filename)
}

func TestCommonOptionsInstallFlagsWithDefaults(t *testing.T) {
	flags := pflag.NewFlagSet("testing", pflag.ContinueOnError)
	opts := NewCommonOptions()
	opts.InstallFlags(flags)

	err := flags.Parse([]string{})
	assert.NilError(t, err)
	assert.Equal(t, opts.TLSOptions.CAFile, defaultPath("ca.pem"))
	assert.Equal(t, opts.TLSOptions.CertFile, defaultPath("cert.pem"))
	assert.Equal(t, opts.TLSOptions.KeyFile, defaultPath("key.pem"))
}
