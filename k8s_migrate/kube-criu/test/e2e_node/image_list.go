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

package e2e_node

import (
	"fmt"
	"os/exec"
	"os/user"
	"time"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/util/sets"
	internalapi "k8s.io/kubernetes/pkg/kubelet/apis/cri"
	runtimeapi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
	commontest "k8s.io/kubernetes/test/e2e/common"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

const (
	// Number of attempts to pull an image.
	maxImagePullRetries = 5
	// Sleep duration between image pull retry attempts.
	imagePullRetryDelay = time.Second
)

// NodeImageWhiteList is a list of images used in node e2e test. These images will be prepulled
// before test running so that the image pulling won't fail in actual test.
var NodeImageWhiteList = sets.NewString(
	"google/cadvisor:latest",
	"k8s.gcr.io/stress:v1",
	busyboxImage,
	"k8s.gcr.io/busybox@sha256:4bdd623e848417d96127e16037743f0cd8b528c026e9175e22a84f639eca58ff",
	"k8s.gcr.io/node-problem-detector:v0.4.1",
	imageutils.GetE2EImage(imageutils.Nginx),
	imageutils.GetE2EImage(imageutils.ServeHostname),
	imageutils.GetE2EImage(imageutils.Netexec),
	imageutils.GetE2EImage(imageutils.Nonewprivs),
	imageutils.GetPauseImageName(),
	framework.GetGPUDevicePluginImage(),
)

func init() {
	// Union NodeImageWhiteList and CommonImageWhiteList into the framework image white list.
	framework.ImageWhiteList = NodeImageWhiteList.Union(commontest.CommonImageWhiteList)
}

// puller represents a generic image puller
type puller interface {
	// Pull pulls an image by name
	Pull(image string) ([]byte, error)
	// Name returns the name of the specific puller implementation
	Name() string
}

type dockerPuller struct {
}

func (dp *dockerPuller) Name() string {
	return "docker"
}

func (dp *dockerPuller) Pull(image string) ([]byte, error) {
	// TODO(random-liu): Use docker client to get rid of docker binary dependency.
	return exec.Command("docker", "pull", image).CombinedOutput()
}

type remotePuller struct {
	imageService internalapi.ImageManagerService
}

func (rp *remotePuller) Name() string {
	return "CRI"
}

func (rp *remotePuller) Pull(image string) ([]byte, error) {
	imageStatus, err := rp.imageService.ImageStatus(&runtimeapi.ImageSpec{Image: image})
	if err == nil && imageStatus != nil {
		return nil, nil
	}
	_, err = rp.imageService.PullImage(&runtimeapi.ImageSpec{Image: image}, nil)
	return nil, err
}

func getPuller() (puller, error) {
	runtime := framework.TestContext.ContainerRuntime
	switch runtime {
	case "docker":
		return &dockerPuller{}, nil
	case "remote":
		_, is, err := getCRIClient()
		if err != nil {
			return nil, err
		}
		return &remotePuller{
			imageService: is,
		}, nil
	}
	return nil, fmt.Errorf("can't prepull images, unknown container runtime %q", runtime)
}

// Pre-fetch all images tests depend on so that we don't fail in an actual test.
func PrePullAllImages() error {
	puller, err := getPuller()
	if err != nil {
		return err
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	images := framework.ImageWhiteList.List()
	glog.V(4).Infof("Pre-pulling images with %s %+v", puller.Name(), images)
	for _, image := range images {
		var (
			err    error
			output []byte
		)
		for i := 0; i < maxImagePullRetries; i++ {
			if i > 0 {
				time.Sleep(imagePullRetryDelay)
			}
			if output, err = puller.Pull(image); err == nil {
				break
			}
			glog.Warningf("Failed to pull %s as user %q, retrying in %s (%d of %d): %v",
				image, usr.Username, imagePullRetryDelay.String(), i+1, maxImagePullRetries, err)
		}
		if err != nil {
			glog.Warningf("Could not pre-pull image %s %v output: %s", image, err, output)
			return err
		}
	}
	return nil
}
