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

// This file should be consistent with pkg/api/v1/annotation_key_constants.go.

package core

const (
	// ImagePolicyFailedOpenKey is added to pods created by failing open when the image policy
	// webhook backend fails.
	ImagePolicyFailedOpenKey string = "alpha.image-policy.k8s.io/failed-open"

	// PodPresetOptOutAnnotationKey represents the annotation key for a pod to exempt itself from pod preset manipulation
	PodPresetOptOutAnnotationKey string = "podpreset.admission.kubernetes.io/exclude"

	// MirrorAnnotationKey represents the annotation key set by kubelets when creating mirror pods
	MirrorPodAnnotationKey string = "kubernetes.io/config.mirror"

	// TolerationsAnnotationKey represents the key of tolerations data (json serialized)
	// in the Annotations of a Pod.
	TolerationsAnnotationKey string = "scheduler.alpha.kubernetes.io/tolerations"

	// TaintsAnnotationKey represents the key of taints data (json serialized)
	// in the Annotations of a Node.
	TaintsAnnotationKey string = "scheduler.alpha.kubernetes.io/taints"

	// SeccompPodAnnotationKey represents the key of a seccomp profile applied
	// to all containers of a pod.
	SeccompPodAnnotationKey string = "seccomp.security.alpha.kubernetes.io/pod"

	// SeccompContainerAnnotationKeyPrefix represents the key of a seccomp profile applied
	// to one container of a pod.
	SeccompContainerAnnotationKeyPrefix string = "container.seccomp.security.alpha.kubernetes.io/"

	// SeccompProfileRuntimeDefault represents the default seccomp profile used by container runtime.
	SeccompProfileRuntimeDefault string = "runtime/default"

	// DeprecatedSeccompProfileDockerDefault represents the default seccomp profile used by docker.
	// This is now deprecated and should be replaced by SeccompProfileRuntimeDefault.
	DeprecatedSeccompProfileDockerDefault string = "docker/default"

	// PreferAvoidPodsAnnotationKey represents the key of preferAvoidPods data (json serialized)
	// in the Annotations of a Node.
	PreferAvoidPodsAnnotationKey string = "scheduler.alpha.kubernetes.io/preferAvoidPods"

	// ObjectTTLAnnotations represents a suggestion for kubelet for how long it can cache
	// an object (e.g. secret, config map) before fetching it again from apiserver.
	// This annotation can be attached to node.
	ObjectTTLAnnotationKey string = "node.alpha.kubernetes.io/ttl"

	// BootstrapCheckpointAnnotationKey represents a Resource (Pod) that should be checkpointed by
	// the kubelet prior to running
	BootstrapCheckpointAnnotationKey string = "node.kubernetes.io/bootstrap-checkpoint"

	// annotation key prefix used to identify non-convertible json paths.
	NonConvertibleAnnotationPrefix = "non-convertible.kubernetes.io"

	kubectlPrefix = "kubectl.kubernetes.io/"

	// LastAppliedConfigAnnotation is the annotation used to store the previous
	// configuration of a resource for use in a three way diff by UpdateApplyAnnotation.
	LastAppliedConfigAnnotation = kubectlPrefix + "last-applied-configuration"

	// AnnotationLoadBalancerSourceRangesKey is the key of the annotation on a service to set allowed ingress ranges on their LoadBalancers
	//
	// It should be a comma-separated list of CIDRs, e.g. `0.0.0.0/0` to
	// allow full access (the default) or `18.0.0.0/8,56.0.0.0/8` to allow
	// access only from the CIDRs currently allocated to MIT & the USPS.
	//
	// Not all cloud providers support this annotation, though AWS & GCE do.
	AnnotationLoadBalancerSourceRangesKey = "service.beta.kubernetes.io/load-balancer-source-ranges"
)
