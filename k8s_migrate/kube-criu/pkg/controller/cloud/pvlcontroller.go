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

package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	"k8s.io/kubernetes/pkg/features"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	volumeutil "k8s.io/kubernetes/pkg/volume/util"

	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const initializerName = "pvlabel.kubernetes.io"

// PersistentVolumeLabelController handles adding labels to persistent volumes when they are created
type PersistentVolumeLabelController struct {
	cloud         cloudprovider.Interface
	kubeClient    kubernetes.Interface
	pvlController cache.Controller
	pvlIndexer    cache.Indexer
	volumeLister  corelisters.PersistentVolumeLister

	syncHandler func(key string) error

	// queue is where incoming work is placed to de-dup and to allow "easy" rate limited requeues on errors
	queue workqueue.RateLimitingInterface
}

// NewPersistentVolumeLabelController creates a PersistentVolumeLabelController object
func NewPersistentVolumeLabelController(
	kubeClient kubernetes.Interface,
	cloud cloudprovider.Interface) *PersistentVolumeLabelController {

	pvlc := &PersistentVolumeLabelController{
		cloud:      cloud,
		kubeClient: kubeClient,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pvLabels"),
	}
	pvlc.syncHandler = pvlc.addLabelsAndAffinity
	pvlc.pvlIndexer, pvlc.pvlController = cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.IncludeUninitialized = true
				return kubeClient.CoreV1().PersistentVolumes().List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.IncludeUninitialized = true
				return kubeClient.CoreV1().PersistentVolumes().Watch(options)
			},
		},
		&v1.PersistentVolume{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					pvlc.queue.Add(key)
				}
			},
		},
		cache.Indexers{},
	)
	pvlc.volumeLister = corelisters.NewPersistentVolumeLister(pvlc.pvlIndexer)

	return pvlc
}

// Run starts a controller that adds labels to persistent volumes
func (pvlc *PersistentVolumeLabelController) Run(threadiness int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer pvlc.queue.ShutDown()

	glog.Infof("Starting PersistentVolumeLabelController")
	defer glog.Infof("Shutting down PersistentVolumeLabelController")

	go pvlc.pvlController.Run(stopCh)

	if !controller.WaitForCacheSync("persistent volume label", stopCh, pvlc.pvlController.HasSynced) {
		return
	}

	// start up your worker threads based on threadiness.  Some controllers have multiple kinds of workers
	for i := 0; i < threadiness; i++ {
		// runWorker will loop until "something bad" happens.  The .Until will then rekick the worker
		// after one second
		go wait.Until(pvlc.runWorker, time.Second, stopCh)
	}

	// wait until we're told to stop
	<-stopCh
}

func (pvlc *PersistentVolumeLabelController) runWorker() {
	// hot loop until we're told to stop.  processNextWorkItem will automatically wait until there's work
	// available, so we don't worry about secondary waits
	for pvlc.processNextWorkItem() {
	}
}

// processNextWorkItem deals with one key off the queue.  It returns false when it's time to quit.
func (pvlc *PersistentVolumeLabelController) processNextWorkItem() bool {
	// pull the next work item from queue.  It should be a key we use to lookup something in a cache
	keyObj, quit := pvlc.queue.Get()
	if quit {
		return false
	}
	// you always have to indicate to the queue that you've completed a piece of work
	defer pvlc.queue.Done(keyObj)

	key := keyObj.(string)
	// do your work on the key.  This method will contains your "do stuff" logic
	err := pvlc.syncHandler(key)
	if err == nil {
		// if you had no error, tell the queue to stop tracking history for your key.  This will
		// reset things like failure counts for per-item rate limiting
		pvlc.queue.Forget(key)
		return true
	}

	// there was a failure so be sure to report it.  This method allows for pluggable error handling
	// which can be used for things like cluster-monitoring
	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", key, err))

	// since we failed, we should requeue the item to work on later.  This method will add a backoff
	// to avoid hotlooping on particular items (they're probably still not going to work right away)
	// and overall controller protection (everything I've done is broken, this controller needs to
	// calm down or it can starve other useful work) cases.
	pvlc.queue.AddRateLimited(key)

	return true
}

// AddLabels adds appropriate labels to persistent volumes and sets the
// volume as available if successful.
func (pvlc *PersistentVolumeLabelController) addLabelsAndAffinity(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return fmt.Errorf("error getting name of volume %q to get volume from informer: %v", key, err)
	}
	volume, err := pvlc.volumeLister.Get(name)
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("error getting volume %s from informer: %v", name, err)
	}

	return pvlc.addLabelsAndAffinityToVolume(volume)
}

func (pvlc *PersistentVolumeLabelController) addLabelsAndAffinityToVolume(vol *v1.PersistentVolume) error {
	var volumeLabels map[string]string
	// Only add labels if the next pending initializer.
	if needsInitialization(vol.Initializers, initializerName) {
		if labeler, ok := (pvlc.cloud).(cloudprovider.PVLabeler); ok {
			labels, err := labeler.GetLabelsForVolume(context.TODO(), vol)
			if err != nil {
				return fmt.Errorf("error querying volume %v: %v", vol.Spec, err)
			}
			volumeLabels = labels
		} else {
			glog.V(4).Info("cloud provider does not support PVLabeler")
		}
		return pvlc.updateVolume(vol, volumeLabels)
	}
	return nil
}

func (pvlc *PersistentVolumeLabelController) createPatch(vol *v1.PersistentVolume, volLabels map[string]string) ([]byte, error) {
	volName := vol.Name
	newVolume := vol.DeepCopyObject().(*v1.PersistentVolume)
	populateAffinity := utilfeature.DefaultFeatureGate.Enabled(features.VolumeScheduling) && len(volLabels) != 0

	if newVolume.Labels == nil {
		newVolume.Labels = make(map[string]string)
	}

	requirements := make([]v1.NodeSelectorRequirement, 0)
	for k, v := range volLabels {
		newVolume.Labels[k] = v
		// Set NodeSelectorRequirements based on the labels
		if populateAffinity {
			var values []string
			if k == kubeletapis.LabelZoneFailureDomain {
				zones, err := volumeutil.LabelZonesToSet(v)
				if err != nil {
					return nil, fmt.Errorf("failed to convert label string for Zone: %s to a Set", v)
				}
				values = zones.List()
			} else {
				values = []string{v}
			}
			requirements = append(requirements, v1.NodeSelectorRequirement{Key: k, Operator: v1.NodeSelectorOpIn, Values: values})
		}
	}
	if populateAffinity {
		if newVolume.Spec.NodeAffinity == nil {
			newVolume.Spec.NodeAffinity = new(v1.VolumeNodeAffinity)
		}
		if newVolume.Spec.NodeAffinity.Required == nil {
			newVolume.Spec.NodeAffinity.Required = new(v1.NodeSelector)
		}
		if len(newVolume.Spec.NodeAffinity.Required.NodeSelectorTerms) == 0 {
			// Need atleast one term pre-allocated whose MatchExpressions can be appended to
			newVolume.Spec.NodeAffinity.Required.NodeSelectorTerms = make([]v1.NodeSelectorTerm, 1)
		}
		// Populate NodeAffinity with requirements if there are no conflicting keys found
		if v1helper.NodeSelectorRequirementKeysExistInNodeSelectorTerms(requirements, newVolume.Spec.NodeAffinity.Required.NodeSelectorTerms) {
			glog.V(4).Infof("NodeSelectorRequirements for cloud labels %v conflict with existing NodeAffinity %v. Skipping addition of NodeSelectorRequirements for cloud labels.",
				requirements, newVolume.Spec.NodeAffinity)
		} else {
			for _, req := range requirements {
				for i := range newVolume.Spec.NodeAffinity.Required.NodeSelectorTerms {
					newVolume.Spec.NodeAffinity.Required.NodeSelectorTerms[i].MatchExpressions = append(newVolume.Spec.NodeAffinity.Required.NodeSelectorTerms[i].MatchExpressions, req)
				}
			}
		}
	}
	newVolume.Initializers = removeInitializer(newVolume.Initializers, initializerName)
	glog.V(4).Infof("removed initializer on PersistentVolume %s", newVolume.Name)

	oldData, err := json.Marshal(vol)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal old persistentvolume %#v for persistentvolume %q: %v", vol, volName, err)
	}

	newData, err := json.Marshal(newVolume)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new persistentvolume %#v for persistentvolume %q: %v", newVolume, volName, err)
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.PersistentVolume{})
	if err != nil {
		return nil, fmt.Errorf("failed to create patch for persistentvolume %q: %v", volName, err)
	}
	return patch, nil
}

func (pvlc *PersistentVolumeLabelController) updateVolume(vol *v1.PersistentVolume, volLabels map[string]string) error {
	volName := vol.Name
	glog.V(4).Infof("updating PersistentVolume %s", volName)
	patchBytes, err := pvlc.createPatch(vol, volLabels)
	if err != nil {
		return err
	}

	_, err = pvlc.kubeClient.CoreV1().PersistentVolumes().Patch(string(volName), types.StrategicMergePatchType, patchBytes)
	if err != nil {
		return fmt.Errorf("failed to update PersistentVolume %s: %v", volName, err)
	}
	glog.V(4).Infof("updated PersistentVolume %s", volName)

	return nil
}

func removeInitializer(initializers *metav1.Initializers, name string) *metav1.Initializers {
	if initializers == nil {
		return nil
	}

	var updated []metav1.Initializer
	for _, pending := range initializers.Pending {
		if pending.Name != name {
			updated = append(updated, pending)
		}
	}
	if len(updated) == len(initializers.Pending) {
		return initializers
	}
	if len(updated) == 0 {
		return nil
	}

	return &metav1.Initializers{Pending: updated}
}

// needsInitialization checks whether or not the PVL is the next pending initializer.
func needsInitialization(initializers *metav1.Initializers, name string) bool {
	if initializers == nil {
		return false
	}

	if len(initializers.Pending) == 0 {
		return false
	}

	// There is at least one initializer still pending so check to
	// see if the PVL is the next in line.
	return initializers.Pending[0].Name == name
}
