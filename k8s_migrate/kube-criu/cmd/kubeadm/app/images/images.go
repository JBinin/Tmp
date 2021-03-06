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

package images

import (
	"fmt"
	"runtime"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
)

// GetGenericImage generates and returns a platform agnostic image (backed by manifest list)
func GetGenericImage(prefix, image, tag string) string {
	return fmt.Sprintf("%s/%s:%s", prefix, image, tag)
}

// GetGenericArchImage generates and returns an image based on the current runtime arch
func GetGenericArchImage(prefix, image, tag string) string {
	return fmt.Sprintf("%s/%s-%s:%s", prefix, image, runtime.GOARCH, tag)
}

// GetKubeControlPlaneImage generates and returns the image for the core Kubernetes components or returns the unified control plane image if specified
func GetKubeControlPlaneImage(image string, cfg *kubeadmapi.InitConfiguration) string {
	if cfg.UnifiedControlPlaneImage != "" {
		return cfg.UnifiedControlPlaneImage
	}
	repoPrefix := cfg.GetControlPlaneImageRepository()
	kubernetesImageTag := kubeadmutil.KubernetesVersionToImageTag(cfg.KubernetesVersion)
	return GetGenericArchImage(repoPrefix, image, kubernetesImageTag)
}

// GetEtcdImage generates and returns the image for etcd or returns cfg.Etcd.Local.Image if specified
func GetEtcdImage(cfg *kubeadmapi.InitConfiguration) string {
	if cfg.Etcd.Local != nil && cfg.Etcd.Local.Image != "" {
		return cfg.Etcd.Local.Image
	}
	etcdImageTag := constants.DefaultEtcdVersion
	etcdImageVersion, err := constants.EtcdSupportedVersion(cfg.KubernetesVersion)
	if err == nil {
		etcdImageTag = etcdImageVersion.String()
	}
	return GetGenericArchImage(cfg.ImageRepository, constants.Etcd, etcdImageTag)
}

// GetAllImages returns a list of container images kubeadm expects to use on a control plane node
func GetAllImages(cfg *kubeadmapi.InitConfiguration) []string {
	imgs := []string{}
	imgs = append(imgs, GetKubeControlPlaneImage(constants.KubeAPIServer, cfg))
	imgs = append(imgs, GetKubeControlPlaneImage(constants.KubeControllerManager, cfg))
	imgs = append(imgs, GetKubeControlPlaneImage(constants.KubeScheduler, cfg))
	imgs = append(imgs, GetKubeControlPlaneImage(constants.KubeProxy, cfg))

	// pause, etcd and kube-dns are not available on the ci image repository so use the default image repository.
	imgs = append(imgs, GetGenericImage(cfg.ImageRepository, "pause", "3.1"))

	// if etcd is not external then add the image as it will be required
	if cfg.Etcd.Local != nil {
		imgs = append(imgs, GetEtcdImage(cfg))
	}

	// Append the appropriate DNS images
	if features.Enabled(cfg.FeatureGates, features.CoreDNS) {
		imgs = append(imgs, GetGenericImage(cfg.ImageRepository, constants.CoreDNS, constants.CoreDNSVersion))
	} else {
		imgs = append(imgs, GetGenericArchImage(cfg.ImageRepository, "k8s-dns-kube-dns", constants.KubeDNSVersion))
		imgs = append(imgs, GetGenericArchImage(cfg.ImageRepository, "k8s-dns-sidecar", constants.KubeDNSVersion))
		imgs = append(imgs, GetGenericArchImage(cfg.ImageRepository, "k8s-dns-dnsmasq-nanny", constants.KubeDNSVersion))
	}

	return imgs
}
