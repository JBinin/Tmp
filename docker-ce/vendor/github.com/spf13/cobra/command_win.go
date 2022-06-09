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

package cobra

import (
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
)

var preExecHookFn = preExecHook

// enables an information splash screen on Windows if the CLI is started from explorer.exe.
var MousetrapHelpText string = `This is a command line tool

You need to open cmd.exe and run it from there.
`

func preExecHook(c *Command) {
	if mousetrap.StartedByExplorer() {
		c.Print(MousetrapHelpText)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
}
