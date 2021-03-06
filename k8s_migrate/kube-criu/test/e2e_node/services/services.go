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
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/golang/glog"
	"github.com/kardianos/osext"

	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/test/e2e/framework"
)

// E2EServices starts and stops e2e services in a separate process. The test
// uses it to start and stop all e2e services.
type E2EServices struct {
	// monitorParent determines whether the sub-processes should watch and die with the current
	// process.
	rmDirs        []string
	monitorParent bool
	services      *server
	kubelet       *server
	logs          logFiles
}

// NewE2EServices returns a new E2EServices instance.
func NewE2EServices(monitorParent bool) *E2EServices {
	return &E2EServices{
		monitorParent: monitorParent,
		// Special log files that need to be collected for additional debugging.
		logs: getLogFiles(),
	}
}

// Start starts the e2e services in another process by calling back into the
// test binary.  Returns when all e2e services are ready or an error.
//
// We want to statically link e2e services into the test binary, but we don't
// want their glog output to pollute the test result. So we run the binary in
// run-services-mode to start e2e services in another process.
// The function starts 2 processes:
// * internal e2e services: services which statically linked in the test binary - apiserver, etcd and
// namespace controller.
// * kubelet: kubelet binary is outside. (We plan to move main kubelet start logic out when we have
// standard kubelet launcher)
func (e *E2EServices) Start() error {
	var err error
	if !framework.TestContext.NodeConformance {
		// Start kubelet
		e.kubelet, err = e.startKubelet()
		if err != nil {
			return fmt.Errorf("failed to start kubelet: %v", err)
		}
	}
	e.services, err = e.startInternalServices()
	return err
}

// Stop stops the e2e services.
func (e *E2EServices) Stop() {
	defer func() {
		if !framework.TestContext.NodeConformance {
			// Collect log files.
			e.collectLogFiles()
		}
	}()
	if e.services != nil {
		if err := e.services.kill(); err != nil {
			glog.Errorf("Failed to stop services: %v", err)
		}
	}
	if e.kubelet != nil {
		if err := e.kubelet.kill(); err != nil {
			glog.Errorf("Failed to stop kubelet: %v", err)
		}
	}
	if e.rmDirs != nil {
		for _, d := range e.rmDirs {
			err := os.RemoveAll(d)
			if err != nil {
				glog.Errorf("Failed to delete directory %s: %v", d, err)
			}
		}
	}
}

// RunE2EServices actually start the e2e services. This function is used to
// start e2e services in current process. This is only used in run-services-mode.
func RunE2EServices() {
	// Populate global DefaultFeatureGate with value from TestContext.FeatureGates.
	// This way, statically-linked components see the same feature gate config as the test context.
	utilfeature.DefaultFeatureGate.SetFromMap(framework.TestContext.FeatureGates)
	e := newE2EServices()
	if err := e.run(); err != nil {
		glog.Fatalf("Failed to run e2e services: %v", err)
	}
}

const (
	// services.log is the combined log of all services
	servicesLogFile = "services.log"
	// LOG_VERBOSITY_LEVEL is consistent with the level used in a cluster e2e test.
	LOG_VERBOSITY_LEVEL = "4"
)

// startInternalServices starts the internal services in a separate process.
func (e *E2EServices) startInternalServices() (*server, error) {
	testBin, err := osext.Executable()
	if err != nil {
		return nil, fmt.Errorf("can't get current binary: %v", err)
	}
	// Pass all flags into the child process, so that it will see the same flag set.
	startCmd := exec.Command(testBin, append([]string{"--run-services-mode"}, os.Args[1:]...)...)
	server := newServer("services", startCmd, nil, nil, getServicesHealthCheckURLs(), servicesLogFile, e.monitorParent, false)
	return server, server.start()
}

// collectLogFiles collects logs of interest either via journalctl or by creating sym
// links. Since we scp files from the remote directory, symlinks will be
// treated as normal files and file contents will be copied over.
func (e *E2EServices) collectLogFiles() {
	// Nothing to do if report dir is not specified.
	if framework.TestContext.ReportDir == "" {
		return
	}
	glog.Info("Fetching log files...")
	journaldFound := isJournaldAvailable()
	for targetFileName, log := range e.logs {
		targetLink := path.Join(framework.TestContext.ReportDir, targetFileName)
		if journaldFound {
			// Skip log files that do not have an equivalent in journald-based machines.
			if len(log.JournalctlCommand) == 0 {
				continue
			}
			glog.Infof("Get log file %q with journalctl command %v.", targetFileName, log.JournalctlCommand)
			out, err := exec.Command("journalctl", log.JournalctlCommand...).CombinedOutput()
			if err != nil {
				glog.Errorf("failed to get %q from journald: %v, %v", targetFileName, string(out), err)
			} else {
				if err = ioutil.WriteFile(targetLink, out, 0644); err != nil {
					glog.Errorf("failed to write logs to %q: %v", targetLink, err)
				}
			}
			continue
		}
		for _, file := range log.Files {
			if _, err := os.Stat(file); err != nil {
				// Expected file not found on this distro.
				continue
			}
			if err := copyLogFile(file, targetLink); err != nil {
				glog.Error(err)
			} else {
				break
			}
		}
	}
}

// isJournaldAvailable returns whether the system executing the tests uses
// journald.
func isJournaldAvailable() bool {
	_, err := exec.LookPath("journalctl")
	return err == nil
}

func copyLogFile(src, target string) error {
	// If not a journald based distro, then just symlink files.
	if out, err := exec.Command("cp", src, target).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to copy %q to %q: %v, %v", src, target, out, err)
	}
	if out, err := exec.Command("chmod", "a+r", target).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to make log file %q world readable: %v, %v", target, out, err)
	}
	return nil
}
