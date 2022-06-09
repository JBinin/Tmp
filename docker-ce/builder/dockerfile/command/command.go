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
// Package command contains the set of Dockerfile commands.
package command

// Define constants for the command strings
const (
	Add         = "add"
	Arg         = "arg"
	Cmd         = "cmd"
	Copy        = "copy"
	Entrypoint  = "entrypoint"
	Env         = "env"
	Expose      = "expose"
	From        = "from"
	Healthcheck = "healthcheck"
	Label       = "label"
	Maintainer  = "maintainer"
	Onbuild     = "onbuild"
	Run         = "run"
	Shell       = "shell"
	StopSignal  = "stopsignal"
	User        = "user"
	Volume      = "volume"
	Workdir     = "workdir"
)

// Commands is list of all Dockerfile commands
var Commands = map[string]struct{}{
	Add:         {},
	Arg:         {},
	Cmd:         {},
	Copy:        {},
	Entrypoint:  {},
	Env:         {},
	Expose:      {},
	From:        {},
	Healthcheck: {},
	Label:       {},
	Maintainer:  {},
	Onbuild:     {},
	Run:         {},
	Shell:       {},
	StopSignal:  {},
	User:        {},
	Volume:      {},
	Workdir:     {},
}
