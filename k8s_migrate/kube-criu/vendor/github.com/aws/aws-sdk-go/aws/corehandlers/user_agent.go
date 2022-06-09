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
package corehandlers

import (
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
)

// SDKVersionUserAgentHandler is a request handler for adding the SDK Version
// to the user agent.
var SDKVersionUserAgentHandler = request.NamedHandler{
	Name: "core.SDKVersionUserAgentHandler",
	Fn: request.MakeAddToUserAgentHandler(aws.SDKName, aws.SDKVersion,
		runtime.Version(), runtime.GOOS, runtime.GOARCH),
}

const execEnvVar = `AWS_EXECUTION_ENV`
const execEnvUAKey = `exec_env`

// AddHostExecEnvUserAgentHander is a request handler appending the SDK's
// execution environment to the user agent.
//
// If the environment variable AWS_EXECUTION_ENV is set, its value will be
// appended to the user agent string.
var AddHostExecEnvUserAgentHander = request.NamedHandler{
	Name: "core.AddHostExecEnvUserAgentHander",
	Fn: func(r *request.Request) {
		v := os.Getenv(execEnvVar)
		if len(v) == 0 {
			return
		}

		request.AddToUserAgent(r, execEnvUAKey+"/"+v)
	},
}
