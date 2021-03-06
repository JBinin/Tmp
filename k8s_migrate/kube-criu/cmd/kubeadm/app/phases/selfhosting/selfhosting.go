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

package selfhosting

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/glog"

	apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
)

const (
	// selfHostingWaitTimeout describes the maximum amount of time a self-hosting wait process should wait before timing out
	selfHostingWaitTimeout = 2 * time.Minute

	// selfHostingFailureThreshold describes how many times kubeadm will retry creating the DaemonSets
	selfHostingFailureThreshold int = 5
)

// CreateSelfHostedControlPlane is responsible for turning a Static Pod-hosted control plane to a self-hosted one
// It achieves that task this way:
// 1. Load the Static Pod specification from disk (from /etc/kubernetes/manifests)
// 2. Extract the PodSpec from that Static Pod specification
// 3. Mutate the PodSpec to be compatible with self-hosting (add the right labels, taints, etc. so it can schedule correctly)
// 4. Build a new DaemonSet object for the self-hosted component in question. Use the above mentioned PodSpec
// 5. Create the DaemonSet resource. Wait until the Pods are running.
// 6. Remove the Static Pod manifest file. The kubelet will stop the original Static Pod-hosted component that was running.
// 7. The self-hosted containers should now step up and take over.
// 8. In order to avoid race conditions, we have to make sure that static pod is deleted correctly before we continue
//      Otherwise, there is a race condition when we proceed without kubelet having restarted the API server correctly and the next .Create call flakes
// 9. Do that for the kube-apiserver, kube-controller-manager and kube-scheduler in a loop
func CreateSelfHostedControlPlane(manifestsDir, kubeConfigDir string, cfg *kubeadmapi.InitConfiguration, client clientset.Interface, waiter apiclient.Waiter, dryRun bool) error {
	glog.V(1).Infoln("creating self hosted control plane")
	// Adjust the timeout slightly to something self-hosting specific
	waiter.SetTimeout(selfHostingWaitTimeout)

	// Here the map of different mutators to use for the control plane's PodSpec is stored
	glog.V(1).Infoln("getting mutators")
	mutators := GetMutatorsFromFeatureGates(cfg.FeatureGates)

	// Some extra work to be done if we should store the control plane certificates in Secrets
	if features.Enabled(cfg.FeatureGates, features.StoreCertsInSecrets) {

		// Upload the certificates and kubeconfig files from disk to the cluster as Secrets
		if err := uploadTLSSecrets(client, cfg.CertificatesDir); err != nil {
			return err
		}
		if err := uploadKubeConfigSecrets(client, kubeConfigDir); err != nil {
			return err
		}
	}

	for _, componentName := range kubeadmconstants.MasterComponents {
		start := time.Now()
		manifestPath := kubeadmconstants.GetStaticPodFilepath(componentName, manifestsDir)

		// Since we want this function to be idempotent; just continue and try the next component if this file doesn't exist
		if _, err := os.Stat(manifestPath); err != nil {
			fmt.Printf("[self-hosted] The Static Pod for the component %q doesn't seem to be on the disk; trying the next one\n", componentName)
			continue
		}

		// Load the Static Pod spec in order to be able to create a self-hosted variant of that file
		podSpec, err := loadPodSpecFromFile(manifestPath)
		if err != nil {
			return err
		}

		// Build a DaemonSet object from the loaded PodSpec
		ds := BuildDaemonSet(componentName, podSpec, mutators)

		// Create or update the DaemonSet in the API Server, and retry selfHostingFailureThreshold times if it errors out
		if err := apiclient.TryRunCommand(func() error {
			return apiclient.CreateOrUpdateDaemonSet(client, ds)
		}, selfHostingFailureThreshold); err != nil {
			return err
		}

		// Wait for the self-hosted component to come up
		if err := waiter.WaitForPodsWithLabel(BuildSelfHostedComponentLabelQuery(componentName)); err != nil {
			return err
		}

		// Remove the old Static Pod manifest if not dryrunning
		if !dryRun {
			if err := os.RemoveAll(manifestPath); err != nil {
				return fmt.Errorf("unable to delete static pod manifest for %s [%v]", componentName, err)
			}
		}

		// Wait for the mirror Pod hash to be removed; otherwise we'll run into race conditions here when the kubelet hasn't had time to
		// remove the Static Pod (or the mirror Pod respectively). This implicitly also tests that the API server endpoint is healthy,
		// because this blocks until the API server returns a 404 Not Found when getting the Static Pod
		staticPodName := fmt.Sprintf("%s-%s", componentName, cfg.NodeRegistration.Name)
		if err := waiter.WaitForPodToDisappear(staticPodName); err != nil {
			return err
		}

		// Just as an extra safety check; make sure the API server is returning ok at the /healthz endpoint (although we know it could return a GET answer for a Pod above)
		if err := waiter.WaitForAPI(); err != nil {
			return err
		}

		fmt.Printf("[self-hosted] self-hosted %s ready after %f seconds\n", componentName, time.Since(start).Seconds())
	}
	return nil
}

// BuildDaemonSet is responsible for mutating the PodSpec and returns a DaemonSet which is suitable for self-hosting
func BuildDaemonSet(name string, podSpec *v1.PodSpec, mutators map[string][]PodSpecMutatorFunc) *apps.DaemonSet {

	// Mutate the PodSpec so it's suitable for self-hosting
	mutatePodSpec(mutators, name, podSpec)

	// Return a DaemonSet based on that Spec
	return &apps.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeadmconstants.AddSelfHostedPrefix(name),
			Namespace: metav1.NamespaceSystem,
			Labels:    BuildSelfhostedComponentLabels(name),
		},
		Spec: apps.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: BuildSelfhostedComponentLabels(name),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: BuildSelfhostedComponentLabels(name),
				},
				Spec: *podSpec,
			},
			UpdateStrategy: apps.DaemonSetUpdateStrategy{
				// Make the DaemonSet utilize the RollingUpdate rollout strategy
				Type: apps.RollingUpdateDaemonSetStrategyType,
			},
		},
	}
}

// BuildSelfhostedComponentLabels returns the labels for a self-hosted component
func BuildSelfhostedComponentLabels(component string) map[string]string {
	return map[string]string{
		"k8s-app": kubeadmconstants.AddSelfHostedPrefix(component),
	}
}

// BuildSelfHostedComponentLabelQuery creates the right query for matching a self-hosted Pod
func BuildSelfHostedComponentLabelQuery(componentName string) string {
	return fmt.Sprintf("k8s-app=%s", kubeadmconstants.AddSelfHostedPrefix(componentName))
}

func loadPodSpecFromFile(filePath string) (*v1.PodSpec, error) {
	podDef, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file path %s: %+v", filePath, err)
	}

	if len(podDef) == 0 {
		return nil, fmt.Errorf("file was empty: %s", filePath)
	}

	codec := clientscheme.Codecs.UniversalDecoder()
	pod := &v1.Pod{}

	if err = runtime.DecodeInto(codec, podDef, pod); err != nil {
		return nil, fmt.Errorf("failed decoding pod: %v", err)
	}

	return &pod.Spec, nil
}
