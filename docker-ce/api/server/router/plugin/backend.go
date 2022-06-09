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
	"io"
	"net/http"

	"github.com/docker/distribution/reference"
	enginetypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// Backend for Plugin
type Backend interface {
	Disable(name string, config *enginetypes.PluginDisableConfig) error
	Enable(name string, config *enginetypes.PluginEnableConfig) error
	List(filters.Args) ([]enginetypes.Plugin, error)
	Inspect(name string) (*enginetypes.Plugin, error)
	Remove(name string, config *enginetypes.PluginRmConfig) error
	Set(name string, args []string) error
	Privileges(ctx context.Context, ref reference.Named, metaHeaders http.Header, authConfig *enginetypes.AuthConfig) (enginetypes.PluginPrivileges, error)
	Pull(ctx context.Context, ref reference.Named, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, privileges enginetypes.PluginPrivileges, outStream io.Writer) error
	Push(ctx context.Context, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, outStream io.Writer) error
	Upgrade(ctx context.Context, ref reference.Named, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, privileges enginetypes.PluginPrivileges, outStream io.Writer) error
	CreateFromContext(ctx context.Context, tarCtx io.ReadCloser, options *enginetypes.PluginCreateOptions) error
}
