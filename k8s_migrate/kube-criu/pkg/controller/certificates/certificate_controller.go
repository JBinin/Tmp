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

// Package certificates implements an abstract controller that is useful for
// building controllers that manage CSRs
package certificates

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"golang.org/x/time/rate"

	certificates "k8s.io/api/certificates/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	certificatesinformers "k8s.io/client-go/informers/certificates/v1beta1"
	clientset "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	certificateslisters "k8s.io/client-go/listers/certificates/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/controller"
)

type CertificateController struct {
	kubeClient clientset.Interface

	csrLister  certificateslisters.CertificateSigningRequestLister
	csrsSynced cache.InformerSynced

	handler func(*certificates.CertificateSigningRequest) error

	queue workqueue.RateLimitingInterface
}

func NewCertificateController(
	kubeClient clientset.Interface,
	csrInformer certificatesinformers.CertificateSigningRequestInformer,
	handler func(*certificates.CertificateSigningRequest) error,
) *CertificateController {
	// Send events to the apiserver
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	cc := &CertificateController{
		kubeClient: kubeClient,
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.NewMaxOfRateLimiter(
			workqueue.NewItemExponentialFailureRateLimiter(200*time.Millisecond, 1000*time.Second),
			// 10 qps, 100 bucket size.  This is only for retry speed and its only the overall factor (not per item)
			&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
		), "certificate"),
		handler: handler,
	}

	// Manage the addition/update of certificate requests
	csrInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			csr := obj.(*certificates.CertificateSigningRequest)
			glog.V(4).Infof("Adding certificate request %s", csr.Name)
			cc.enqueueCertificateRequest(obj)
		},
		UpdateFunc: func(old, new interface{}) {
			oldCSR := old.(*certificates.CertificateSigningRequest)
			glog.V(4).Infof("Updating certificate request %s", oldCSR.Name)
			cc.enqueueCertificateRequest(new)
		},
		DeleteFunc: func(obj interface{}) {
			csr, ok := obj.(*certificates.CertificateSigningRequest)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					glog.V(2).Infof("Couldn't get object from tombstone %#v", obj)
					return
				}
				csr, ok = tombstone.Obj.(*certificates.CertificateSigningRequest)
				if !ok {
					glog.V(2).Infof("Tombstone contained object that is not a CSR: %#v", obj)
					return
				}
			}
			glog.V(4).Infof("Deleting certificate request %s", csr.Name)
			cc.enqueueCertificateRequest(obj)
		},
	})
	cc.csrLister = csrInformer.Lister()
	cc.csrsSynced = csrInformer.Informer().HasSynced
	return cc
}

// Run the main goroutine responsible for watching and syncing jobs.
func (cc *CertificateController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer cc.queue.ShutDown()

	glog.Infof("Starting certificate controller")
	defer glog.Infof("Shutting down certificate controller")

	if !controller.WaitForCacheSync("certificate", stopCh, cc.csrsSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(cc.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker runs a thread that dequeues CSRs, handles them, and marks them done.
func (cc *CertificateController) worker() {
	for cc.processNextWorkItem() {
	}
}

// processNextWorkItem deals with one key off the queue.  It returns false when it's time to quit.
func (cc *CertificateController) processNextWorkItem() bool {
	cKey, quit := cc.queue.Get()
	if quit {
		return false
	}
	defer cc.queue.Done(cKey)

	if err := cc.syncFunc(cKey.(string)); err != nil {
		cc.queue.AddRateLimited(cKey)
		if _, ignorable := err.(ignorableError); !ignorable {
			utilruntime.HandleError(fmt.Errorf("Sync %v failed with : %v", cKey, err))
		} else {
			glog.V(4).Infof("Sync %v failed with : %v", cKey, err)
		}
		return true
	}

	cc.queue.Forget(cKey)
	return true

}

func (cc *CertificateController) enqueueCertificateRequest(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for object %+v: %v", obj, err))
		return
	}
	cc.queue.Add(key)
}

// maybeSignCertificate will inspect the certificate request and, if it has
// been approved and meets policy expectations, generate an X509 cert using the
// cluster CA assets. If successful it will update the CSR approve subresource
// with the signed certificate.
func (cc *CertificateController) syncFunc(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing certificate request %q (%v)", key, time.Since(startTime))
	}()
	csr, err := cc.csrLister.Get(key)
	if errors.IsNotFound(err) {
		glog.V(3).Infof("csr has been deleted: %v", key)
		return nil
	}
	if err != nil {
		return err
	}

	if csr.Status.Certificate != nil {
		// no need to do anything because it already has a cert
		return nil
	}

	// need to operate on a copy so we don't mutate the csr in the shared cache
	csr = csr.DeepCopy()

	return cc.handler(csr)
}

// IgnorableError returns an error that we shouldn't handle (i.e. log) because
// it's spammy and usually user error. Instead we will log these errors at a
// higher log level. We still need to throw these errors to signal that the
// sync should be retried.
func IgnorableError(s string, args ...interface{}) ignorableError {
	return ignorableError(fmt.Sprintf(s, args...))
}

type ignorableError string

func (e ignorableError) Error() string {
	return string(e)
}
