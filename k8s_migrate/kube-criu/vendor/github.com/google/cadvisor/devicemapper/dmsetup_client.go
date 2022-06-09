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
// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package devicemapper

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// DmsetupClient is a low-level client for interacting with device mapper via
// the `dmsetup` utility, which is provided by the `device-mapper` package.
type DmsetupClient interface {
	// Table runs `dmsetup table` on the given device name and returns the
	// output or an error.
	Table(deviceName string) ([]byte, error)
	// Message runs `dmsetup message` on the given device, passing the given
	// message to the given sector, and returns the output or an error.
	Message(deviceName string, sector int, message string) ([]byte, error)
	// Status runs `dmsetup status` on the given device and returns the output
	// or an error.
	Status(deviceName string) ([]byte, error)
}

// NewDmSetupClient returns a new DmsetupClient.
func NewDmsetupClient() DmsetupClient {
	return &defaultDmsetupClient{}
}

// defaultDmsetupClient is a functional DmsetupClient
type defaultDmsetupClient struct{}

var _ DmsetupClient = &defaultDmsetupClient{}

func (c *defaultDmsetupClient) Table(deviceName string) ([]byte, error) {
	return c.dmsetup("table", deviceName)
}

func (c *defaultDmsetupClient) Message(deviceName string, sector int, message string) ([]byte, error) {
	return c.dmsetup("message", deviceName, strconv.Itoa(sector), message)
}

func (c *defaultDmsetupClient) Status(deviceName string) ([]byte, error) {
	return c.dmsetup("status", deviceName)
}

func (*defaultDmsetupClient) dmsetup(args ...string) ([]byte, error) {
	glog.V(5).Infof("running dmsetup %v", strings.Join(args, " "))
	return exec.Command("dmsetup", args...).Output()
}
