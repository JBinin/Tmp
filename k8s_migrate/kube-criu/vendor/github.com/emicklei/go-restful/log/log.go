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
package log

import (
	stdlog "log"
	"os"
)

// StdLogger corresponds to a minimal subset of the interface satisfied by stdlib log.Logger
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

var Logger StdLogger

func init() {
	// default Logger
	SetLogger(stdlog.New(os.Stderr, "[restful] ", stdlog.LstdFlags|stdlog.Lshortfile))
}

// SetLogger sets the logger for this package
func SetLogger(customLogger StdLogger) {
	Logger = customLogger
}

// Print delegates to the Logger
func Print(v ...interface{}) {
	Logger.Print(v...)
}

// Printf delegates to the Logger
func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}
