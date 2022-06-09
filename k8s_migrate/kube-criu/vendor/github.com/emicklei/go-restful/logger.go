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
package restful

// Copyright 2014 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.
import (
	"github.com/emicklei/go-restful/log"
)

var trace bool = false
var traceLogger log.StdLogger

func init() {
	traceLogger = log.Logger // use the package logger by default
}

// TraceLogger enables detailed logging of Http request matching and filter invocation. Default no logger is set.
// You may call EnableTracing() directly to enable trace logging to the package-wide logger.
func TraceLogger(logger log.StdLogger) {
	traceLogger = logger
	EnableTracing(logger != nil)
}

// expose the setter for the global logger on the top-level package
func SetLogger(customLogger log.StdLogger) {
	log.SetLogger(customLogger)
}

// EnableTracing can be used to Trace logging on and off.
func EnableTracing(enabled bool) {
	trace = enabled
}
