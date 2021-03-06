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

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/pkg/stringid"
	icmd "github.com/docker/docker/pkg/testutil/cmd"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestRmiWithContainerFails(c *check.C) {
	errSubstr := "is using it"

	// create a container
	out, _ := dockerCmd(c, "run", "-d", "busybox", "true")

	cleanedContainerID := strings.TrimSpace(out)

	// try to delete the image
	out, _, err := dockerCmdWithError("rmi", "busybox")
	// Container is using image, should not be able to rmi
	c.Assert(err, checker.NotNil)
	// Container is using image, error message should contain errSubstr
	c.Assert(out, checker.Contains, errSubstr, check.Commentf("Container: %q", cleanedContainerID))

	// make sure it didn't delete the busybox name
	images, _ := dockerCmd(c, "images")
	// The name 'busybox' should not have been removed from images
	c.Assert(images, checker.Contains, "busybox")
}

func (s *DockerSuite) TestRmiTag(c *check.C) {
	imagesBefore, _ := dockerCmd(c, "images", "-a")
	dockerCmd(c, "tag", "busybox", "utest:tag1")
	dockerCmd(c, "tag", "busybox", "utest/docker:tag2")
	dockerCmd(c, "tag", "busybox", "utest:5000/docker:tag3")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+3, check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))
	}
	dockerCmd(c, "rmi", "utest/docker:tag2")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+2, check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))
	}
	dockerCmd(c, "rmi", "utest:5000/docker:tag3")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+1, check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))

	}
	dockerCmd(c, "rmi", "utest:tag1")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n"), check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))

	}
}

func (s *DockerSuite) TestRmiImgIDMultipleTag(c *check.C) {
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "mkdir '/busybox-one'")

	containerID := strings.TrimSpace(out)

	// Wait for it to exit as cannot commit a running container on Windows, and
	// it will take a few seconds to exit
	if testEnv.DaemonPlatform() == "windows" {
		err := waitExited(containerID, 60*time.Second)
		c.Assert(err, check.IsNil)
	}

	dockerCmd(c, "commit", containerID, "busybox-one")

	imagesBefore, _ := dockerCmd(c, "images", "-a")
	dockerCmd(c, "tag", "busybox-one", "busybox-one:tag1")
	dockerCmd(c, "tag", "busybox-one", "busybox-one:tag2")

	imagesAfter, _ := dockerCmd(c, "images", "-a")
	// tag busybox to create 2 more images with same imageID
	c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+2, check.Commentf("docker images shows: %q\n", imagesAfter))

	imgID := inspectField(c, "busybox-one:tag1", "Id")

	// run a container with the image
	out, _ = runSleepingContainerInImage(c, "busybox-one")

	containerID = strings.TrimSpace(out)

	// first checkout without force it fails
	out, _, err := dockerCmdWithError("rmi", imgID)
	expected := fmt.Sprintf("conflict: unable to delete %s (cannot be forced) - image is being used by running container %s", stringid.TruncateID(imgID), stringid.TruncateID(containerID))
	// rmi tagged in multiple repos should have failed without force
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, expected)

	dockerCmd(c, "stop", containerID)
	dockerCmd(c, "rmi", "-f", imgID)

	imagesAfter, _ = dockerCmd(c, "images", "-a")
	// rmi -f failed, image still exists
	c.Assert(imagesAfter, checker.Not(checker.Contains), imgID[:12], check.Commentf("ImageID:%q; ImagesAfter: %q", imgID, imagesAfter))
}

func (s *DockerSuite) TestRmiImgIDForce(c *check.C) {
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "mkdir '/busybox-test'")

	containerID := strings.TrimSpace(out)

	// Wait for it to exit as cannot commit a running container on Windows, and
	// it will take a few seconds to exit
	if testEnv.DaemonPlatform() == "windows" {
		err := waitExited(containerID, 60*time.Second)
		c.Assert(err, check.IsNil)
	}

	dockerCmd(c, "commit", containerID, "busybox-test")

	imagesBefore, _ := dockerCmd(c, "images", "-a")
	dockerCmd(c, "tag", "busybox-test", "utest:tag1")
	dockerCmd(c, "tag", "busybox-test", "utest:tag2")
	dockerCmd(c, "tag", "busybox-test", "utest/docker:tag3")
	dockerCmd(c, "tag", "busybox-test", "utest:5000/docker:tag4")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+4, check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))
	}
	imgID := inspectField(c, "busybox-test", "Id")

	// first checkout without force it fails
	out, _, err := dockerCmdWithError("rmi", imgID)
	// rmi tagged in multiple repos should have failed without force
	c.Assert(err, checker.NotNil)
	// rmi tagged in multiple repos should have failed without force
	c.Assert(out, checker.Contains, "(must be forced) - image is referenced in multiple repositories", check.Commentf("out: %s; err: %v;", out, err))

	dockerCmd(c, "rmi", "-f", imgID)
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		// rmi failed, image still exists
		c.Assert(imagesAfter, checker.Not(checker.Contains), imgID[:12])
	}
}

// See https://github.com/docker/docker/issues/14116
func (s *DockerSuite) TestRmiImageIDForceWithRunningContainersAndMultipleTags(c *check.C) {
	dockerfile := "FROM busybox\nRUN echo test 14116\n"
	buildImageSuccessfully(c, "test-14116", withDockerfile(dockerfile))
	imgID := getIDByName(c, "test-14116")

	newTag := "newtag"
	dockerCmd(c, "tag", imgID, newTag)
	runSleepingContainerInImage(c, imgID)

	out, _, err := dockerCmdWithError("rmi", "-f", imgID)
	// rmi -f should not delete image with running containers
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, "(cannot be forced) - image is being used by running container")
}

func (s *DockerSuite) TestRmiTagWithExistingContainers(c *check.C) {
	container := "test-delete-tag"
	newtag := "busybox:newtag"
	bb := "busybox:latest"
	dockerCmd(c, "tag", bb, newtag)

	dockerCmd(c, "run", "--name", container, bb, "/bin/true")

	out, _ := dockerCmd(c, "rmi", newtag)
	c.Assert(strings.Count(out, "Untagged: "), checker.Equals, 1)
}

func (s *DockerSuite) TestRmiForceWithExistingContainers(c *check.C) {
	image := "busybox-clone"

	icmd.RunCmd(icmd.Cmd{
		Command: []string{dockerBinary, "build", "--no-cache", "-t", image, "-"},
		Stdin: strings.NewReader(`FROM busybox
MAINTAINER foo`),
	}).Assert(c, icmd.Success)

	dockerCmd(c, "run", "--name", "test-force-rmi", image, "/bin/true")

	dockerCmd(c, "rmi", "-f", image)
}

func (s *DockerSuite) TestRmiWithMultipleRepositories(c *check.C) {
	newRepo := "127.0.0.1:5000/busybox"
	oldRepo := "busybox"
	newTag := "busybox:test"
	dockerCmd(c, "tag", oldRepo, newRepo)

	dockerCmd(c, "run", "--name", "test", oldRepo, "touch", "/abcd")

	dockerCmd(c, "commit", "test", newTag)

	out, _ := dockerCmd(c, "rmi", newTag)
	c.Assert(out, checker.Contains, "Untagged: "+newTag)
}

func (s *DockerSuite) TestRmiForceWithMultipleRepositories(c *check.C) {
	imageName := "rmiimage"
	tag1 := imageName + ":tag1"
	tag2 := imageName + ":tag2"

	buildImageSuccessfully(c, tag1, withDockerfile(`FROM busybox
		MAINTAINER "docker"`))
	dockerCmd(c, "tag", tag1, tag2)

	out, _ := dockerCmd(c, "rmi", "-f", tag2)
	c.Assert(out, checker.Contains, "Untagged: "+tag2)
	c.Assert(out, checker.Not(checker.Contains), "Untagged: "+tag1)

	// Check built image still exists
	images, _ := dockerCmd(c, "images", "-a")
	c.Assert(images, checker.Contains, imageName, check.Commentf("Built image missing %q; Images: %q", imageName, images))
}

func (s *DockerSuite) TestRmiBlank(c *check.C) {
	out, _, err := dockerCmdWithError("rmi", " ")
	// Should have failed to delete ' ' image
	c.Assert(err, checker.NotNil)
	// Wrong error message generated
	c.Assert(out, checker.Not(checker.Contains), "no such id", check.Commentf("out: %s", out))
	// Expected error message not generated
	c.Assert(out, checker.Contains, "image name cannot be blank", check.Commentf("out: %s", out))
}

func (s *DockerSuite) TestRmiContainerImageNotFound(c *check.C) {
	// Build 2 images for testing.
	imageNames := []string{"test1", "test2"}
	imageIds := make([]string, 2)
	for i, name := range imageNames {
		dockerfile := fmt.Sprintf("FROM busybox\nMAINTAINER %s\nRUN echo %s\n", name, name)
		buildImageSuccessfully(c, name, withoutCache, withDockerfile(dockerfile))
		id := getIDByName(c, name)
		imageIds[i] = id
	}

	// Create a long-running container.
	runSleepingContainerInImage(c, imageNames[0])

	// Create a stopped container, and then force remove its image.
	dockerCmd(c, "run", imageNames[1], "true")
	dockerCmd(c, "rmi", "-f", imageIds[1])

	// Try to remove the image of the running container and see if it fails as expected.
	out, _, err := dockerCmdWithError("rmi", "-f", imageIds[0])
	// The image of the running container should not be removed.
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, "image is being used by running container", check.Commentf("out: %s", out))
}

// #13422
func (s *DockerSuite) TestRmiUntagHistoryLayer(c *check.C) {
	image := "tmp1"
	// Build an image for testing.
	dockerfile := `FROM busybox
MAINTAINER foo
RUN echo 0 #layer0
RUN echo 1 #layer1
RUN echo 2 #layer2
`
	buildImageSuccessfully(c, image, withoutCache, withDockerfile(dockerfile))
	out, _ := dockerCmd(c, "history", "-q", image)
	ids := strings.Split(out, "\n")
	idToTag := ids[2]

	// Tag layer0 to "tmp2".
	newTag := "tmp2"
	dockerCmd(c, "tag", idToTag, newTag)
	// Create a container based on "tmp1".
	dockerCmd(c, "run", "-d", image, "true")

	// See if the "tmp2" can be untagged.
	out, _ = dockerCmd(c, "rmi", newTag)
	// Expected 1 untagged entry
	c.Assert(strings.Count(out, "Untagged: "), checker.Equals, 1, check.Commentf("out: %s", out))

	// Now let's add the tag again and create a container based on it.
	dockerCmd(c, "tag", idToTag, newTag)
	out, _ = dockerCmd(c, "run", "-d", newTag, "true")
	cid := strings.TrimSpace(out)

	// At this point we have 2 containers, one based on layer2 and another based on layer0.
	// Try to untag "tmp2" without the -f flag.
	out, _, err := dockerCmdWithError("rmi", newTag)
	// should not be untagged without the -f flag
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, cid[:12])
	c.Assert(out, checker.Contains, "(must force)")

	// Add the -f flag and test again.
	out, _ = dockerCmd(c, "rmi", "-f", newTag)
	// should be allowed to untag with the -f flag
	c.Assert(out, checker.Contains, fmt.Sprintf("Untagged: %s:latest", newTag))
}

func (*DockerSuite) TestRmiParentImageFail(c *check.C) {
	buildImageSuccessfully(c, "test", withDockerfile(`
	FROM busybox
	RUN echo hello`))

	id := inspectField(c, "busybox", "ID")
	out, _, err := dockerCmdWithError("rmi", id)
	c.Assert(err, check.NotNil)
	if !strings.Contains(out, "image has dependent child images") {
		c.Fatalf("rmi should have failed because it's a parent image, got %s", out)
	}
}

func (s *DockerSuite) TestRmiWithParentInUse(c *check.C) {
	out, _ := dockerCmd(c, "create", "busybox")
	cID := strings.TrimSpace(out)

	out, _ = dockerCmd(c, "commit", cID)
	imageID := strings.TrimSpace(out)

	out, _ = dockerCmd(c, "create", imageID)
	cID = strings.TrimSpace(out)

	out, _ = dockerCmd(c, "commit", cID)
	imageID = strings.TrimSpace(out)

	dockerCmd(c, "rmi", imageID)
}

// #18873
func (s *DockerSuite) TestRmiByIDHardConflict(c *check.C) {
	dockerCmd(c, "create", "busybox")

	imgID := inspectField(c, "busybox:latest", "Id")

	_, _, err := dockerCmdWithError("rmi", imgID[:12])
	c.Assert(err, checker.NotNil)

	// check that tag was not removed
	imgID2 := inspectField(c, "busybox:latest", "Id")
	c.Assert(imgID, checker.Equals, imgID2)
}
