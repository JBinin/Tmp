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
Copyright 2017 The Kubernetes Authors.

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

package config

import (
	"fmt"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ContainerRuntimeOptions struct {

	// General options.

	// ContainerRuntime is the container runtime to use.
	ContainerRuntime string
	// RuntimeCgroups that container runtime is expected to be isolated in.
	RuntimeCgroups string
	// RedirectContainerStreaming enables container streaming redirect.
	// When RedirectContainerStreaming is false, kubelet will proxy container streaming data
	// between apiserver and container runtime. This approach is more secure, but the proxy
	// introduces some overhead.
	// When RedirectContainerStreaming is true, kubelet will return an http redirect to apiserver,
	// and apiserver will access container runtime directly. This approach is more performant,
	// but less secure because the connection between apiserver and container runtime is not
	// authenticated.
	RedirectContainerStreaming bool

	// Docker-specific options.

	// DockershimRootDirectory is the path to the dockershim root directory. Defaults to
	// /var/lib/dockershim if unset. Exposed for integration testing (e.g. in OpenShift).
	DockershimRootDirectory string
	// Enable dockershim only mode.
	ExperimentalDockershim bool
	// PodSandboxImage is the image whose network/ipc namespaces
	// containers in each pod will use.
	PodSandboxImage string
	// DockerEndpoint is the path to the docker endpoint to communicate with.
	DockerEndpoint string
	// If no pulling progress is made before the deadline imagePullProgressDeadline,
	// the image pulling will be cancelled. Defaults to 1m0s.
	// +optional
	ImagePullProgressDeadline metav1.Duration

	// Network plugin options.

	// networkPluginName is the name of the network plugin to be invoked for
	// various events in kubelet/pod lifecycle
	NetworkPluginName string
	// NetworkPluginMTU is the MTU to be passed to the network plugin,
	// and overrides the default MTU for cases where it cannot be automatically
	// computed (such as IPSEC).
	NetworkPluginMTU int32
	// CNIConfDir is the full path of the directory in which to search for
	// CNI config files
	CNIConfDir string
	// CNIBinDir is the full path of the directory in which to search for
	// CNI plugin binaries
	CNIBinDir string
}

func (s *ContainerRuntimeOptions) AddFlags(fs *pflag.FlagSet) {
	dockerOnlyWarning := "This docker-specific flag only works when container-runtime is set to docker."

	// General settings.
	fs.StringVar(&s.ContainerRuntime, "container-runtime", s.ContainerRuntime, "The container runtime to use. Possible values: 'docker', 'remote', 'rkt (deprecated)'.")
	fs.StringVar(&s.RuntimeCgroups, "runtime-cgroups", s.RuntimeCgroups, "Optional absolute name of cgroups to create and run the runtime in.")
	fs.BoolVar(&s.RedirectContainerStreaming, "redirect-container-streaming", s.RedirectContainerStreaming, "Enables container streaming redirect. If false, kubelet will proxy container streaming data between apiserver and container runtime; if true, kubelet will return an http redirect to apiserver, and apiserver will access container runtime directly. The proxy approach is more secure, but introduces some overhead. The redirect approach is more performant, but less secure because the connection between apiserver and container runtime may not be authenticated.")

	// Docker-specific settings.
	fs.BoolVar(&s.ExperimentalDockershim, "experimental-dockershim", s.ExperimentalDockershim, "Enable dockershim only mode. In this mode, kubelet will only start dockershim without any other functionalities. This flag only serves test purpose, please do not use it unless you are conscious of what you are doing. [default=false]")
	fs.MarkHidden("experimental-dockershim")
	fs.StringVar(&s.DockershimRootDirectory, "experimental-dockershim-root-directory", s.DockershimRootDirectory, "Path to the dockershim root directory.")
	fs.MarkHidden("experimental-dockershim-root-directory")
	fs.StringVar(&s.PodSandboxImage, "pod-infra-container-image", s.PodSandboxImage, fmt.Sprintf("The image whose network/ipc namespaces containers in each pod will use. %s", dockerOnlyWarning))
	fs.StringVar(&s.DockerEndpoint, "docker-endpoint", s.DockerEndpoint, fmt.Sprintf("Use this for the docker endpoint to communicate with. %s", dockerOnlyWarning))
	fs.DurationVar(&s.ImagePullProgressDeadline.Duration, "image-pull-progress-deadline", s.ImagePullProgressDeadline.Duration, fmt.Sprintf("If no pulling progress is made before this deadline, the image pulling will be cancelled. %s", dockerOnlyWarning))

	// Network plugin settings for Docker.
	fs.StringVar(&s.NetworkPluginName, "network-plugin", s.NetworkPluginName, fmt.Sprintf("<Warning: Alpha feature> The name of the network plugin to be invoked for various events in kubelet/pod lifecycle. %s", dockerOnlyWarning))
	fs.StringVar(&s.CNIConfDir, "cni-conf-dir", s.CNIConfDir, fmt.Sprintf("<Warning: Alpha feature> The full path of the directory in which to search for CNI config files. %s", dockerOnlyWarning))
	fs.StringVar(&s.CNIBinDir, "cni-bin-dir", s.CNIBinDir, fmt.Sprintf("<Warning: Alpha feature> A comma-separated list of full paths of directories in which to search for CNI plugin binaries. %s", dockerOnlyWarning))
	fs.Int32Var(&s.NetworkPluginMTU, "network-plugin-mtu", s.NetworkPluginMTU, fmt.Sprintf("<Warning: Alpha feature> The MTU to be passed to the network plugin, to override the default. Set to 0 to use the default 1460 MTU. %s", dockerOnlyWarning))
}
