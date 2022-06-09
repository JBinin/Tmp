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

package apps

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

const (
	NginxImageName = "nginx"
	RedisImageName = "redis"
)

var (
	CronJobGroupVersionResourceAlpha = schema.GroupVersionResource{Group: "batch", Version: "v2alpha1", Resource: "cronjobs"}
	CronJobGroupVersionResourceBeta  = schema.GroupVersionResource{Group: "batch", Version: "v1beta1", Resource: "cronjobs"}
	NautilusImage                    = imageutils.GetE2EImage(imageutils.Nautilus)
	KittenImage                      = imageutils.GetE2EImage(imageutils.Kitten)
	NginxImage                       = imageutils.GetE2EImage(imageutils.Nginx)
	NewNginxImage                    = imageutils.GetE2EImage(imageutils.NginxNew)
	RedisImage                       = imageutils.GetE2EImage(imageutils.Redis)
)
