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
// Copyright 2015 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package invoke

import (
	"os"
	"strings"
)

type CNIArgs interface {
	// For use with os/exec; i.e., return nil to inherit the
	// environment from this process
	AsEnv() []string
}

type inherited struct{}

var inheritArgsFromEnv inherited

func (_ *inherited) AsEnv() []string {
	return nil
}

func ArgsFromEnv() CNIArgs {
	return &inheritArgsFromEnv
}

type Args struct {
	Command       string
	ContainerID   string
	NetNS         string
	PluginArgs    [][2]string
	PluginArgsStr string
	IfName        string
	Path          string
}

// Args implements the CNIArgs interface
var _ CNIArgs = &Args{}

func (args *Args) AsEnv() []string {
	env := os.Environ()
	pluginArgsStr := args.PluginArgsStr
	if pluginArgsStr == "" {
		pluginArgsStr = stringify(args.PluginArgs)
	}

	// Ensure that the custom values are first, so any value present in
	// the process environment won't override them.
	env = append([]string{
		"CNI_COMMAND=" + args.Command,
		"CNI_CONTAINERID=" + args.ContainerID,
		"CNI_NETNS=" + args.NetNS,
		"CNI_ARGS=" + pluginArgsStr,
		"CNI_IFNAME=" + args.IfName,
		"CNI_PATH=" + args.Path,
	}, env...)
	return env
}

// taken from rkt/networking/net_plugin.go
func stringify(pluginArgs [][2]string) string {
	entries := make([]string, len(pluginArgs))

	for i, kv := range pluginArgs {
		entries[i] = strings.Join(kv[:], "=")
	}

	return strings.Join(entries, ";")
}
