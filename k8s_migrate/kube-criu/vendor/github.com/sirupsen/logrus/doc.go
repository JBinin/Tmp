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
/*
Package logrus is a structured logger for Go, completely API compatible with the standard library logger.


The simplest way to use Logrus is simply the package-level exported logger:

  package main

  import (
    log "github.com/sirupsen/logrus"
  )

  func main() {
    log.WithFields(log.Fields{
      "animal": "walrus",
      "number": 1,
      "size":   10,
    }).Info("A walrus appears")
  }

Output:
  time="2015-09-07T08:48:33Z" level=info msg="A walrus appears" animal=walrus number=1 size=10

For a full guide visit https://github.com/sirupsen/logrus
*/
package logrus
