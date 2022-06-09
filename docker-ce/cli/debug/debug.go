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
package debug

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// Enable sets the DEBUG env var to true
// and makes the logger to log at debug level.
func Enable() {
	os.Setenv("DEBUG", "1")
	logrus.SetLevel(logrus.DebugLevel)
}

// Disable sets the DEBUG env var to false
// and makes the logger to log at info level.
func Disable() {
	os.Setenv("DEBUG", "")
	logrus.SetLevel(logrus.InfoLevel)
}

// IsEnabled checks whether the debug flag is set or not.
func IsEnabled() bool {
	return os.Getenv("DEBUG") != ""
}
