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
Copyright 2015 The Kubernetes Authors.

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

// Reads the pod configuration from the Kubernetes apiserver.
package config

import (
	"fmt"
	"time"
	"os"
	"flag"
	"path/filepath"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	api "k8s.io/kubernetes/pkg/apis/core"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"

	"k8s.io/client-go/tools/clientcmd"
	v1alpha1 "k8s.io/podcheckpoint/pkg/apis/podcheckpointcontroller/v1alpha1"
	podcheckpoint_clientset "k8s.io/podcheckpoint/pkg/client/clientset/versioned"
	informers "k8s.io/podcheckpoint/pkg/client/informers/externalversions"
	"k8s.io/podcheckpoint/pkg/signals"
)

// NewSourceApiserver creates a config source that watches and pulls from the apiserver.
func NewSourceApiserver(c clientset.Interface, nodeName types.NodeName, updates chan<- interface{}, configFile string) {
	lw := cache.NewListWatchFromClient(c.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, fields.OneTermEqualSelector(api.PodHostField, string(nodeName)))
	newSourceApiserverFromLW(lw, updates)
	fmt.Println("Start PodCheckpointInformer!!!")
	newPodCheckpointInformer(c, updates, configFile)
}

// newSourceApiserverFromLW holds creates a config source that watches and pulls from the apiserver.
func newSourceApiserverFromLW(lw cache.ListerWatcher, updates chan<- interface{}) {
	send := func(objs []interface{}) {
		var pods []*v1.Pod
		for _, o := range objs {
			pods = append(pods, o.(*v1.Pod))
		}
		updates <- kubetypes.PodUpdate{Pods: pods, Op: kubetypes.SET, Source: kubetypes.ApiserverSource}
	}
	r := cache.NewReflector(lw, &v1.Pod{}, cache.NewUndeltaStore(send, cache.MetaNamespaceKeyFunc), 0)
	go r.Run(wait.NeverStop)
}

func newPodCheckpointInformer(c clientset.Interface, updates chan<- interface{}, configFile string) {
	stopCh := signals.SetupSignalHandler()
	send := func(obj interface{}) {
		fmt.Printf("exec podcheckpoint add func :%v ", obj)
		podcheckpoint := obj.(*v1alpha1.PodCheckpoint)
		if podcheckpoint.Status.Phase == "" {
			podcheckpoint.Status.Phase = v1alpha1.PodPrepareCheckpoint
		}
		updates <- kubetypes.PodUpdate{PodCheckpoint: podcheckpoint, Op: kubetypes.CHECKPOINT, Source: kubetypes.ApiserverSource}
	}

	var kubeconfig string
	if len(configFile) > 0 {
		kubeconfig = configFile
		fmt.Println("kubeletconfigFile : %+v",kubeconfig)
	} else {
		if home := os.Getenv("HOME"); home != "" {
			flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) Absolute path to the kubeconfig file")
			fmt.Println("kubeletconfig : %+v",kubeconfig)
		} else {
			flag.StringVar(&kubeconfig, "kubeconfig", "", "Absolute path to the kubeconfig file")
		}
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Errorf("Error building kubeconfig: %s", err.Error())
	}
	podcheckpointClient, err := podcheckpoint_clientset.NewForConfig(cfg)
	if err != nil {
		fmt.Errorf("Error building example clientset: %s", err.Error())
	}
	podcheckpointInformerFactory := informers.NewSharedInformerFactory(podcheckpointClient, time.Second*30)
	podcheckpointInformer := podcheckpointInformerFactory.Podcheckpointcontroller().V1alpha1().PodCheckpoints()
	podcheckpointInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: send,
	})
	fmt.Println("Start podcheckpointInformerFactory")
	podcheckpointInformerFactory.Start(stopCh)
}

