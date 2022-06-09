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
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSuite) BenchmarkConcurrentContainerActions(c *check.C) {
	maxConcurrency := runtime.GOMAXPROCS(0)
	numIterations := c.N
	outerGroup := &sync.WaitGroup{}
	outerGroup.Add(maxConcurrency)
	chErr := make(chan error, numIterations*2*maxConcurrency)

	for i := 0; i < maxConcurrency; i++ {
		go func() {
			defer outerGroup.Done()
			innerGroup := &sync.WaitGroup{}
			innerGroup.Add(2)

			go func() {
				defer innerGroup.Done()
				for i := 0; i < numIterations; i++ {
					args := []string{"run", "-d", defaultSleepImage}
					args = append(args, sleepCommandForDaemonPlatform()...)
					out, _, err := dockerCmdWithError(args...)
					if err != nil {
						chErr <- fmt.Errorf(out)
						return
					}

					id := strings.TrimSpace(out)
					tmpDir, err := ioutil.TempDir("", "docker-concurrent-test-"+id)
					if err != nil {
						chErr <- err
						return
					}
					defer os.RemoveAll(tmpDir)
					out, _, err = dockerCmdWithError("cp", id+":/tmp", tmpDir)
					if err != nil {
						chErr <- fmt.Errorf(out)
						return
					}

					out, _, err = dockerCmdWithError("kill", id)
					if err != nil {
						chErr <- fmt.Errorf(out)
					}

					out, _, err = dockerCmdWithError("start", id)
					if err != nil {
						chErr <- fmt.Errorf(out)
					}

					out, _, err = dockerCmdWithError("kill", id)
					if err != nil {
						chErr <- fmt.Errorf(out)
					}

					// don't do an rm -f here since it can potentially ignore errors from the graphdriver
					out, _, err = dockerCmdWithError("rm", id)
					if err != nil {
						chErr <- fmt.Errorf(out)
					}
				}
			}()

			go func() {
				defer innerGroup.Done()
				for i := 0; i < numIterations; i++ {
					out, _, err := dockerCmdWithError("ps")
					if err != nil {
						chErr <- fmt.Errorf(out)
					}
				}
			}()

			innerGroup.Wait()
		}()
	}

	outerGroup.Wait()
	close(chErr)

	for err := range chErr {
		c.Assert(err, checker.IsNil)
	}
}
