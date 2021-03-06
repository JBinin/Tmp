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
	"net/http"
	"os"

	"github.com/docker/docker/integration-cli/registry"
	"github.com/go-check/check"
)

func makefile(contents string) (string, func(), error) {
	cleanup := func() {

	}

	f, err := ioutil.TempFile(".", "tmp")
	if err != nil {
		return "", cleanup, err
	}
	err = ioutil.WriteFile(f.Name(), []byte(contents), os.ModePerm)
	if err != nil {
		return "", cleanup, err
	}

	cleanup = func() {
		err := os.Remove(f.Name())
		if err != nil {
			fmt.Println("Error removing tmpfile")
		}
	}
	return f.Name(), cleanup, nil

}

// TestV2Only ensures that a daemon in v2-only mode does not
// attempt to contact any v1 registry endpoints.
func (s *DockerRegistrySuite) TestV2Only(c *check.C) {
	reg, err := registry.NewMock(c)
	defer reg.Close()
	c.Assert(err, check.IsNil)

	reg.RegisterHandler("/v2/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	reg.RegisterHandler("/v1/.*", func(w http.ResponseWriter, r *http.Request) {
		c.Fatal("V1 registry contacted")
	})

	repoName := fmt.Sprintf("%s/busybox", reg.URL())

	s.d.Start(c, "--insecure-registry", reg.URL(), "--disable-legacy-registry=true")

	dockerfileName, cleanup, err := makefile(fmt.Sprintf("FROM %s/busybox", reg.URL()))
	c.Assert(err, check.IsNil, check.Commentf("Unable to create test dockerfile"))
	defer cleanup()

	s.d.Cmd("build", "--file", dockerfileName, ".")

	s.d.Cmd("run", repoName)
	s.d.Cmd("login", "-u", "richard", "-p", "testtest", "-e", "testuser@testdomain.com", reg.URL())
	s.d.Cmd("tag", "busybox", repoName)
	s.d.Cmd("push", repoName)
	s.d.Cmd("pull", repoName)
}

// TestV1 starts a daemon in 'normal' mode
// and ensure v1 endpoints are hit for the following operations:
// login, push, pull, build & run
func (s *DockerRegistrySuite) TestV1(c *check.C) {
	reg, err := registry.NewMock(c)
	defer reg.Close()
	c.Assert(err, check.IsNil)

	v2Pings := 0
	reg.RegisterHandler("/v2/", func(w http.ResponseWriter, r *http.Request) {
		v2Pings++
		// V2 ping 404 causes fallback to v1
		w.WriteHeader(404)
	})

	v1Pings := 0
	reg.RegisterHandler("/v1/_ping", func(w http.ResponseWriter, r *http.Request) {
		v1Pings++
	})

	v1Logins := 0
	reg.RegisterHandler("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		v1Logins++
	})

	v1Repo := 0
	reg.RegisterHandler("/v1/repositories/busybox/", func(w http.ResponseWriter, r *http.Request) {
		v1Repo++
	})

	reg.RegisterHandler("/v1/repositories/busybox/images", func(w http.ResponseWriter, r *http.Request) {
		v1Repo++
	})

	s.d.Start(c, "--insecure-registry", reg.URL(), "--disable-legacy-registry=false")

	dockerfileName, cleanup, err := makefile(fmt.Sprintf("FROM %s/busybox", reg.URL()))
	c.Assert(err, check.IsNil, check.Commentf("Unable to create test dockerfile"))
	defer cleanup()

	s.d.Cmd("build", "--file", dockerfileName, ".")
	c.Assert(v1Repo, check.Equals, 1, check.Commentf("Expected v1 repository access after build"))

	repoName := fmt.Sprintf("%s/busybox", reg.URL())
	s.d.Cmd("run", repoName)
	c.Assert(v1Repo, check.Equals, 2, check.Commentf("Expected v1 repository access after run"))

	s.d.Cmd("login", "-u", "richard", "-p", "testtest", reg.URL())
	c.Assert(v1Logins, check.Equals, 1, check.Commentf("Expected v1 login attempt"))

	s.d.Cmd("tag", "busybox", repoName)
	s.d.Cmd("push", repoName)

	c.Assert(v1Repo, check.Equals, 2)

	s.d.Cmd("pull", repoName)
	c.Assert(v1Repo, check.Equals, 3, check.Commentf("Expected v1 repository access after pull"))
}
