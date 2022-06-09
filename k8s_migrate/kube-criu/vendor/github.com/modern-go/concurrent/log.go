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
package concurrent

import (
	"os"
	"log"
	"io/ioutil"
)

// ErrorLogger is used to print out error, can be set to writer other than stderr
var ErrorLogger = log.New(os.Stderr, "", 0)

// InfoLogger is used to print informational message, default to off
var InfoLogger = log.New(ioutil.Discard, "", 0)
