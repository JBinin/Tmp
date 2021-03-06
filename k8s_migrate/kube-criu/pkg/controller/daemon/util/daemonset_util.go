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

package util

import (
	"fmt"
	"strconv"

	apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	"k8s.io/kubernetes/pkg/features"
	kubelettypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"
)

// GetTemplateGeneration gets the template generation associated with a v1.DaemonSet by extracting it from the
// deprecated annotation. If no annotation is found nil is returned. If the annotation is found and fails to parse
// nil is returned with an error. If the generation can be parsed from the annotation, a pointer to the parsed int64
// value is returned.
func GetTemplateGeneration(ds *apps.DaemonSet) (*int64, error) {
	annotation, found := ds.Annotations[apps.DeprecatedTemplateGeneration]
	if !found {
		return nil, nil
	}
	generation, err := strconv.ParseInt(annotation, 10, 64)
	if err != nil {
		return nil, err
	}
	return &generation, nil
}

// CreatePodTemplate returns copy of provided template with additional
// label which contains templateGeneration (for backward compatibility),
// hash of provided template and sets default daemon tolerations.
func CreatePodTemplate(template v1.PodTemplateSpec, generation *int64, hash string) v1.PodTemplateSpec {
	newTemplate := *template.DeepCopy()
	// DaemonSet pods shouldn't be deleted by NodeController in case of node problems.
	// Add infinite toleration for taint notReady:NoExecute here
	// to survive taint-based eviction enforced by NodeController
	// when node turns not ready.
	v1helper.AddOrUpdateTolerationInPodSpec(&newTemplate.Spec, &v1.Toleration{
		Key:      algorithm.TaintNodeNotReady,
		Operator: v1.TolerationOpExists,
		Effect:   v1.TaintEffectNoExecute,
	})

	// DaemonSet pods shouldn't be deleted by NodeController in case of node problems.
	// Add infinite toleration for taint unreachable:NoExecute here
	// to survive taint-based eviction enforced by NodeController
	// when node turns unreachable.
	v1helper.AddOrUpdateTolerationInPodSpec(&newTemplate.Spec, &v1.Toleration{
		Key:      algorithm.TaintNodeUnreachable,
		Operator: v1.TolerationOpExists,
		Effect:   v1.TaintEffectNoExecute,
	})

	// According to TaintNodesByCondition feature, all DaemonSet pods should tolerate
	// MemoryPressure and DisPressure taints, and the critical pods should tolerate
	// OutOfDisk taint.
	v1helper.AddOrUpdateTolerationInPodSpec(&newTemplate.Spec, &v1.Toleration{
		Key:      algorithm.TaintNodeDiskPressure,
		Operator: v1.TolerationOpExists,
		Effect:   v1.TaintEffectNoSchedule,
	})

	v1helper.AddOrUpdateTolerationInPodSpec(&newTemplate.Spec, &v1.Toleration{
		Key:      algorithm.TaintNodeMemoryPressure,
		Operator: v1.TolerationOpExists,
		Effect:   v1.TaintEffectNoSchedule,
	})

	// TODO(#48843) OutOfDisk taints will be removed in 1.10
	if utilfeature.DefaultFeatureGate.Enabled(features.ExperimentalCriticalPodAnnotation) &&
		kubelettypes.IsCritical(newTemplate.Namespace, newTemplate.Annotations) {
		v1helper.AddOrUpdateTolerationInPodSpec(&newTemplate.Spec, &v1.Toleration{
			Key:      algorithm.TaintNodeOutOfDisk,
			Operator: v1.TolerationOpExists,
			Effect:   v1.TaintEffectNoExecute,
		})
	}

	if newTemplate.ObjectMeta.Labels == nil {
		newTemplate.ObjectMeta.Labels = make(map[string]string)
	}
	if generation != nil {
		newTemplate.ObjectMeta.Labels[extensions.DaemonSetTemplateGenerationKey] = fmt.Sprint(*generation)
	}
	// TODO: do we need to validate if the DaemonSet is RollingUpdate or not?
	if len(hash) > 0 {
		newTemplate.ObjectMeta.Labels[extensions.DefaultDaemonSetUniqueLabelKey] = hash
	}
	return newTemplate
}

// IsPodUpdated checks if pod contains label value that either matches templateGeneration or hash
func IsPodUpdated(pod *v1.Pod, hash string, dsTemplateGeneration *int64) bool {
	// Compare with hash to see if the pod is updated, need to maintain backward compatibility of templateGeneration
	templateMatches := dsTemplateGeneration != nil &&
		pod.Labels[extensions.DaemonSetTemplateGenerationKey] == fmt.Sprint(dsTemplateGeneration)
	hashMatches := len(hash) > 0 && pod.Labels[extensions.DefaultDaemonSetUniqueLabelKey] == hash
	return hashMatches || templateMatches
}

// SplitByAvailablePods splits provided daemon set pods by availability
func SplitByAvailablePods(minReadySeconds int32, pods []*v1.Pod) ([]*v1.Pod, []*v1.Pod) {
	unavailablePods := []*v1.Pod{}
	availablePods := []*v1.Pod{}
	for _, pod := range pods {
		if podutil.IsPodAvailable(pod, minReadySeconds, metav1.Now()) {
			availablePods = append(availablePods, pod)
		} else {
			unavailablePods = append(unavailablePods, pod)
		}
	}
	return availablePods, unavailablePods
}

// ReplaceDaemonSetPodNodeNameNodeAffinity replaces the RequiredDuringSchedulingIgnoredDuringExecution
// NodeAffinity of the given affinity with a new NodeAffinity that selects the given nodeName.
// Note that this function assumes that no NodeAffinity conflicts with the selected nodeName.
func ReplaceDaemonSetPodNodeNameNodeAffinity(affinity *v1.Affinity, nodename string) *v1.Affinity {
	nodeSelReq := v1.NodeSelectorRequirement{
		Key:      algorithm.NodeFieldSelectorKeyNodeName,
		Operator: v1.NodeSelectorOpIn,
		Values:   []string{nodename},
	}

	nodeSelector := &v1.NodeSelector{
		NodeSelectorTerms: []v1.NodeSelectorTerm{
			{
				MatchFields: []v1.NodeSelectorRequirement{nodeSelReq},
			},
		},
	}

	if affinity == nil {
		return &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: nodeSelector,
			},
		}
	}

	if affinity.NodeAffinity == nil {
		affinity.NodeAffinity = &v1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: nodeSelector,
		}
		return affinity
	}

	nodeAffinity := affinity.NodeAffinity

	if nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = nodeSelector
		return affinity
	}

	// Replace node selector with the new one.
	nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = []v1.NodeSelectorTerm{
		{
			MatchFields: []v1.NodeSelectorRequirement{nodeSelReq},
		},
	}

	return affinity
}

// AppendNoScheduleTolerationIfNotExist appends unschedulable toleration to `.spec` if not exist; otherwise,
// no changes to `.spec.tolerations`.
func AppendNoScheduleTolerationIfNotExist(tolerations []v1.Toleration) []v1.Toleration {
	unschedulableToleration := v1.Toleration{
		Key:      algorithm.TaintNodeUnschedulable,
		Operator: v1.TolerationOpExists,
		Effect:   v1.TaintEffectNoSchedule,
	}

	unschedulableTaintExist := false

	for _, t := range tolerations {
		if apiequality.Semantic.DeepEqual(t, unschedulableToleration) {
			unschedulableTaintExist = true
			break
		}
	}

	if !unschedulableTaintExist {
		tolerations = append(tolerations, unschedulableToleration)
	}

	return tolerations
}

// GetTargetNodeName get the target node name of DaemonSet pods. If `.spec.NodeName` is not empty (nil),
// return `.spec.NodeName`; otherwise, retrieve node name of pending pods from NodeAffinity. Return error
// if failed to retrieve node name from `.spec.NodeName` and NodeAffinity.
func GetTargetNodeName(pod *v1.Pod) (string, error) {
	if len(pod.Spec.NodeName) != 0 {
		return pod.Spec.NodeName, nil
	}

	// If ScheduleDaemonSetPods was enabled before, retrieve node name of unscheduled pods from NodeAffinity
	if pod.Spec.Affinity == nil ||
		pod.Spec.Affinity.NodeAffinity == nil ||
		pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		return "", fmt.Errorf("no spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution for pod %s/%s",
			pod.Namespace, pod.Name)
	}

	terms := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	if len(terms) < 1 {
		return "", fmt.Errorf("no nodeSelectorTerms in requiredDuringSchedulingIgnoredDuringExecution of pod %s/%s",
			pod.Namespace, pod.Name)
	}

	for _, term := range terms {
		for _, exp := range term.MatchFields {
			if exp.Key == algorithm.NodeFieldSelectorKeyNodeName &&
				exp.Operator == v1.NodeSelectorOpIn {
				if len(exp.Values) != 1 {
					return "", fmt.Errorf("the matchFields value of '%s' is not unique for pod %s/%s",
						algorithm.NodeFieldSelectorKeyNodeName, pod.Namespace, pod.Name)
				}

				return exp.Values[0], nil
			}
		}
	}

	return "", fmt.Errorf("no node name found for pod %s/%s", pod.Namespace, pod.Name)
}
