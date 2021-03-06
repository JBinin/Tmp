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
Copyright 2014 The Kubernetes Authors.

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

package types

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	kubeapi "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/scheduling"
	"k8s.io/kubernetes/pkg/features"

	v1alpha1 "k8s.io/podcheckpoint/pkg/apis/podcheckpointcontroller/v1alpha1"
)

const (
	ConfigSourceAnnotationKey    = "kubernetes.io/config.source"
	ConfigMirrorAnnotationKey    = v1.MirrorPodAnnotationKey
	ConfigFirstSeenAnnotationKey = "kubernetes.io/config.seen"
	ConfigHashAnnotationKey      = "kubernetes.io/config.hash"
	CriticalPodAnnotationKey     = "scheduler.alpha.kubernetes.io/critical-pod"
)

// PodOperation defines what changes will be made on a pod configuration.
type PodOperation int

const (
	// This is the current pod configuration
	SET PodOperation = iota
	// Pods with the given ids are new to this source
	ADD
	// Pods with the given ids are gracefully deleted from this source
	DELETE
	// Pods with the given ids have been removed from this source
	REMOVE
	// Pods with the given ids have been updated in this source
	UPDATE
	// Pods with the given ids have unexpected status in this source,
	// kubelet should reconcile status with this source
	RECONCILE
	// Pods with the given ids have been restored from a checkpoint.
	RESTORE
	//Checkpoint Pods
	CHECKPOINT

	// These constants identify the sources of pods
	// Updates from a file
	FileSource = "file"
	// Updates from querying a web page
	HTTPSource = "http"
	// Updates from Kubernetes API Server
	ApiserverSource = "api"
	// Updates from all sources
	AllSource = "*"

	NamespaceDefault = metav1.NamespaceDefault
)

// PodUpdate defines an operation sent on the channel. You can add or remove single services by
// sending an array of size one and Op == ADD|REMOVE (with REMOVE, only the ID is required).
// For setting the state of the system to a given state for this source configuration, set
// Pods as desired and Op to SET, which will reset the system state to that specified in this
// operation for this source channel. To remove all pods, set Pods to empty object and Op to SET.
//
// Additionally, Pods should never be nil - it should always point to an empty slice. While
// functionally similar, this helps our unit tests properly check that the correct PodUpdates
// are generated.
type PodUpdate struct {
	Pods          []*v1.Pod
	PodCheckpoint *v1alpha1.PodCheckpoint
	Op            PodOperation
	Source        string
}

// Gets all validated sources from the specified sources.
func GetValidatedSources(sources []string) ([]string, error) {
	validated := make([]string, 0, len(sources))
	for _, source := range sources {
		switch source {
		case AllSource:
			return []string{FileSource, HTTPSource, ApiserverSource}, nil
		case FileSource, HTTPSource, ApiserverSource:
			validated = append(validated, source)
			break
		case "":
			break
		default:
			return []string{}, fmt.Errorf("unknown pod source %q", source)
		}
	}
	return validated, nil
}

// GetPodSource returns the source of the pod based on the annotation.
func GetPodSource(pod *v1.Pod) (string, error) {
	if pod.Annotations != nil {
		if source, ok := pod.Annotations[ConfigSourceAnnotationKey]; ok {
			return source, nil
		}
	}
	return "", fmt.Errorf("cannot get source of pod %q", pod.UID)
}

// SyncPodType classifies pod updates, eg: create, update.
type SyncPodType int

const (
	// SyncPodSync is when the pod is synced to ensure desired state
	SyncPodSync SyncPodType = iota
	// SyncPodUpdate is when the pod is updated from source
	SyncPodUpdate
	// SyncPodCreate is when the pod is created from source
	SyncPodCreate
	// SyncPodKill is when the pod is killed based on a trigger internal to the kubelet for eviction.
	// If a SyncPodKill request is made to pod workers, the request is never dropped, and will always be processed.
	SyncPodKill
	// SyncPodCheckpoint is when the pod is checkpointed from source
	SyncPodCheckpoint
)

func (sp SyncPodType) String() string {
	switch sp {
	case SyncPodCreate:
		return "create"
	case SyncPodUpdate:
		return "update"
	case SyncPodSync:
		return "sync"
	case SyncPodKill:
		return "kill"
	case SyncPodCheckpoint:
		return "checkpoint"
	default:
		return "unknown"
	}
}

// IsCriticalPod returns true if the pod bears the critical pod annotation key or if pod's priority is greater than
// or equal to SystemCriticalPriority. Both the rescheduler(deprecated in 1.10) and the kubelet use this function
// to make admission and scheduling decisions.
func IsCriticalPod(pod *v1.Pod) bool {
	if utilfeature.DefaultFeatureGate.Enabled(features.PodPriority) {
		if pod.Spec.Priority != nil && IsCriticalPodBasedOnPriority(*pod.Spec.Priority) {
			return true
		}
	}
	if utilfeature.DefaultFeatureGate.Enabled(features.ExperimentalCriticalPodAnnotation) {
		if IsCritical(pod.Namespace, pod.Annotations) {
			return true
		}
	}
	return false
}

// Preemptable returns true if preemptor pod can preempt preemptee pod:
//   - If preemptor's is greater than preemptee's priority, it's preemptable (return true)
//   - If preemptor (or its priority) is nil and preemptee bears the critical pod annotation key,
//     preemptee can not be preempted (return false)
//   - If preemptor (or its priority) is nil and preemptee's priority is greater than or equal to
//     SystemCriticalPriority, preemptee can not be preempted (return false)
func Preemptable(preemptor, preemptee *v1.Pod) bool {
	if utilfeature.DefaultFeatureGate.Enabled(features.PodPriority) {
		if (preemptor != nil && preemptor.Spec.Priority != nil) &&
			(preemptee != nil && preemptee.Spec.Priority != nil) {
			return *(preemptor.Spec.Priority) > *(preemptee.Spec.Priority)
		}
	}

	return !IsCriticalPod(preemptee)
}

// IsCritical returns true if parameters bear the critical pod annotation
// key. The DaemonSetController use this key directly to make scheduling decisions.
// TODO: @ravig - Deprecated. Remove this when we move to resolving critical pods based on priorityClassName.
func IsCritical(ns string, annotations map[string]string) bool {
	// Critical pods are restricted to "kube-system" namespace as of now.
	if ns != kubeapi.NamespaceSystem {
		return false
	}
	val, ok := annotations[CriticalPodAnnotationKey]
	if ok && val == "" {
		return true
	}
	return false
}

// IsCriticalPodBasedOnPriority checks if the given pod is a critical pod based on priority resolved from pod Spec.
func IsCriticalPodBasedOnPriority(priority int32) bool {
	if priority >= scheduling.SystemCriticalPriority {
		return true
	}
	return false
}

