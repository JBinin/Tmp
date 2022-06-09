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
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-check/check"
)

func (s *DockerSuite) BenchmarkLogsCLIRotateFollow(c *check.C) {
	out, _ := dockerCmd(c, "run", "-d", "--log-opt", "max-size=1b", "--log-opt", "max-file=10", "busybox", "sh", "-c", "while true; do usleep 50000; echo hello; done")
	id := strings.TrimSpace(out)
	ch := make(chan error, 1)
	go func() {
		ch <- nil
		out, _, _ := dockerCmdWithError("logs", "-f", id)
		// if this returns at all, it's an error
		ch <- fmt.Errorf(out)
	}()

	<-ch
	select {
	case <-time.After(30 * time.Second):
		// ran for 30 seconds with no problem
		return
	case err := <-ch:
		if err != nil {
			c.Fatal(err)
		}
	}
}
