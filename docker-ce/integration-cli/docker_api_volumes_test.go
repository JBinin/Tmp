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
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestVolumesAPIList(c *check.C) {
	prefix, _ := getPrefixAndSlashFromDaemonPlatform()
	dockerCmd(c, "run", "-v", prefix+"/foo", "busybox")

	status, b, err := request.SockRequest("GET", "/volumes", nil, daemonHost())
	c.Assert(err, checker.IsNil)
	c.Assert(status, checker.Equals, http.StatusOK)

	var volumes volumetypes.VolumesListOKBody
	c.Assert(json.Unmarshal(b, &volumes), checker.IsNil)

	c.Assert(len(volumes.Volumes), checker.Equals, 1, check.Commentf("\n%v", volumes.Volumes))
}

func (s *DockerSuite) TestVolumesAPICreate(c *check.C) {
	config := volumetypes.VolumesCreateBody{
		Name: "test",
	}
	status, b, err := request.SockRequest("POST", "/volumes/create", config, daemonHost())
	c.Assert(err, check.IsNil)
	c.Assert(status, check.Equals, http.StatusCreated, check.Commentf(string(b)))

	var vol types.Volume
	err = json.Unmarshal(b, &vol)
	c.Assert(err, checker.IsNil)

	c.Assert(filepath.Base(filepath.Dir(vol.Mountpoint)), checker.Equals, config.Name)
}

func (s *DockerSuite) TestVolumesAPIRemove(c *check.C) {
	prefix, _ := getPrefixAndSlashFromDaemonPlatform()
	dockerCmd(c, "run", "-v", prefix+"/foo", "--name=test", "busybox")

	status, b, err := request.SockRequest("GET", "/volumes", nil, daemonHost())
	c.Assert(err, checker.IsNil)
	c.Assert(status, checker.Equals, http.StatusOK)

	var volumes volumetypes.VolumesListOKBody
	c.Assert(json.Unmarshal(b, &volumes), checker.IsNil)
	c.Assert(len(volumes.Volumes), checker.Equals, 1, check.Commentf("\n%v", volumes.Volumes))

	v := volumes.Volumes[0]
	status, _, err = request.SockRequest("DELETE", "/volumes/"+v.Name, nil, daemonHost())
	c.Assert(err, checker.IsNil)
	c.Assert(status, checker.Equals, http.StatusConflict, check.Commentf("Should not be able to remove a volume that is in use"))

	dockerCmd(c, "rm", "-f", "test")
	status, data, err := request.SockRequest("DELETE", "/volumes/"+v.Name, nil, daemonHost())
	c.Assert(err, checker.IsNil)
	c.Assert(status, checker.Equals, http.StatusNoContent, check.Commentf(string(data)))

}

func (s *DockerSuite) TestVolumesAPIInspect(c *check.C) {
	config := volumetypes.VolumesCreateBody{
		Name: "test",
	}
	status, b, err := request.SockRequest("POST", "/volumes/create", config, daemonHost())
	c.Assert(err, check.IsNil)
	c.Assert(status, check.Equals, http.StatusCreated, check.Commentf(string(b)))

	status, b, err = request.SockRequest("GET", "/volumes", nil, daemonHost())
	c.Assert(err, checker.IsNil)
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf(string(b)))

	var volumes volumetypes.VolumesListOKBody
	c.Assert(json.Unmarshal(b, &volumes), checker.IsNil)
	c.Assert(len(volumes.Volumes), checker.Equals, 1, check.Commentf("\n%v", volumes.Volumes))

	var vol types.Volume
	status, b, err = request.SockRequest("GET", "/volumes/"+config.Name, nil, daemonHost())
	c.Assert(err, checker.IsNil)
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf(string(b)))
	c.Assert(json.Unmarshal(b, &vol), checker.IsNil)
	c.Assert(vol.Name, checker.Equals, config.Name)
}
