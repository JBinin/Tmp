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

package auth

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kubernetes/plugin/pkg/admission/serviceaccount"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var mountImage = imageutils.GetE2EImage(imageutils.Mounttest)

var _ = SIGDescribe("ServiceAccounts", func() {
	f := framework.NewDefaultFramework("svcaccounts")

	It("should ensure a single API token exists", func() {
		// wait for the service account to reference a single secret
		var secrets []v1.ObjectReference
		framework.ExpectNoError(wait.Poll(time.Millisecond*500, time.Second*10, func() (bool, error) {
			By("waiting for a single token reference")
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				framework.Logf("default service account was not found")
				return false, nil
			}
			if err != nil {
				framework.Logf("error getting default service account: %v", err)
				return false, err
			}
			switch len(sa.Secrets) {
			case 0:
				framework.Logf("default service account has no secret references")
				return false, nil
			case 1:
				framework.Logf("default service account has a single secret reference")
				secrets = sa.Secrets
				return true, nil
			default:
				return false, fmt.Errorf("default service account has too many secret references: %#v", sa.Secrets)
			}
		}))

		// make sure the reference doesn't flutter
		{
			By("ensuring the single token reference persists")
			time.Sleep(2 * time.Second)
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			framework.ExpectNoError(err)
			Expect(sa.Secrets).To(Equal(secrets))
		}

		// delete the referenced secret
		By("deleting the service account token")
		framework.ExpectNoError(f.ClientSet.CoreV1().Secrets(f.Namespace.Name).Delete(secrets[0].Name, nil))

		// wait for the referenced secret to be removed, and another one autocreated
		framework.ExpectNoError(wait.Poll(time.Millisecond*500, framework.ServiceAccountProvisionTimeout, func() (bool, error) {
			By("waiting for a new token reference")
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			if err != nil {
				framework.Logf("error getting default service account: %v", err)
				return false, err
			}
			switch len(sa.Secrets) {
			case 0:
				framework.Logf("default service account has no secret references")
				return false, nil
			case 1:
				if sa.Secrets[0] == secrets[0] {
					framework.Logf("default service account still has the deleted secret reference")
					return false, nil
				}
				framework.Logf("default service account has a new single secret reference")
				secrets = sa.Secrets
				return true, nil
			default:
				return false, fmt.Errorf("default service account has too many secret references: %#v", sa.Secrets)
			}
		}))

		// make sure the reference doesn't flutter
		{
			By("ensuring the single token reference persists")
			time.Sleep(2 * time.Second)
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			framework.ExpectNoError(err)
			Expect(sa.Secrets).To(Equal(secrets))
		}

		// delete the reference from the service account
		By("deleting the reference to the service account token")
		{
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			framework.ExpectNoError(err)
			sa.Secrets = nil
			_, updateErr := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Update(sa)
			framework.ExpectNoError(updateErr)
		}

		// wait for another one to be autocreated
		framework.ExpectNoError(wait.Poll(time.Millisecond*500, framework.ServiceAccountProvisionTimeout, func() (bool, error) {
			By("waiting for a new token to be created and added")
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			if err != nil {
				framework.Logf("error getting default service account: %v", err)
				return false, err
			}
			switch len(sa.Secrets) {
			case 0:
				framework.Logf("default service account has no secret references")
				return false, nil
			case 1:
				framework.Logf("default service account has a new single secret reference")
				secrets = sa.Secrets
				return true, nil
			default:
				return false, fmt.Errorf("default service account has too many secret references: %#v", sa.Secrets)
			}
		}))

		// make sure the reference doesn't flutter
		{
			By("ensuring the single token reference persists")
			time.Sleep(2 * time.Second)
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			framework.ExpectNoError(err)
			Expect(sa.Secrets).To(Equal(secrets))
		}
	})

	/*
	   Release: v1.9
	   Testname: Service Account Tokens Must AutoMount
	   Description: Ensure that Service Account keys are mounted into the Container. Pod
	                contains three containers each will read Service Account token,
	                root CA and default namespace respectively from the default API
	                Token Mount path. All these three files MUST exist and the Service
	                Account mount path MUST be auto mounted to the Container.
	*/
	framework.ConformanceIt("should mount an API token into pods ", func() {
		var tokenContent string
		var rootCAContent string

		// Standard get, update retry loop
		framework.ExpectNoError(wait.Poll(time.Millisecond*500, framework.ServiceAccountProvisionTimeout, func() (bool, error) {
			By("getting the auto-created API token")
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get("default", metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				framework.Logf("default service account was not found")
				return false, nil
			}
			if err != nil {
				framework.Logf("error getting default service account: %v", err)
				return false, err
			}
			if len(sa.Secrets) == 0 {
				framework.Logf("default service account has no secret references")
				return false, nil
			}
			for _, secretRef := range sa.Secrets {
				secret, err := f.ClientSet.CoreV1().Secrets(f.Namespace.Name).Get(secretRef.Name, metav1.GetOptions{})
				if err != nil {
					framework.Logf("Error getting secret %s: %v", secretRef.Name, err)
					continue
				}
				if secret.Type == v1.SecretTypeServiceAccountToken {
					tokenContent = string(secret.Data[v1.ServiceAccountTokenKey])
					rootCAContent = string(secret.Data[v1.ServiceAccountRootCAKey])
					return true, nil
				}
			}

			framework.Logf("default service account has no secret references to valid service account tokens")
			return false, nil
		}))

		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "pod-service-account-" + string(uuid.NewUUID()) + "-",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "token-test",
						Image: mountImage,
						Args: []string{
							fmt.Sprintf("--file_content=%s/%s", serviceaccount.DefaultAPITokenMountPath, v1.ServiceAccountTokenKey),
						},
					},
					{
						Name:  "root-ca-test",
						Image: mountImage,
						Args: []string{
							fmt.Sprintf("--file_content=%s/%s", serviceaccount.DefaultAPITokenMountPath, v1.ServiceAccountRootCAKey),
						},
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		}

		pod.Spec.Containers = append(pod.Spec.Containers, v1.Container{
			Name:  "namespace-test",
			Image: mountImage,
			Args: []string{
				fmt.Sprintf("--file_content=%s/%s", serviceaccount.DefaultAPITokenMountPath, v1.ServiceAccountNamespaceKey),
			},
		})

		f.TestContainerOutput("consume service account token", pod, 0, []string{
			fmt.Sprintf(`content of file "%s/%s": %s`, serviceaccount.DefaultAPITokenMountPath, v1.ServiceAccountTokenKey, tokenContent),
		})
		f.TestContainerOutput("consume service account root CA", pod, 1, []string{
			fmt.Sprintf(`content of file "%s/%s": %s`, serviceaccount.DefaultAPITokenMountPath, v1.ServiceAccountRootCAKey, rootCAContent),
		})

		f.TestContainerOutput("consume service account namespace", pod, 2, []string{
			fmt.Sprintf(`content of file "%s/%s": %s`, serviceaccount.DefaultAPITokenMountPath, v1.ServiceAccountNamespaceKey, f.Namespace.Name),
		})
	})

	/*
	   Release: v1.9
	   Testname: Service account tokens auto mount optionally
	   Description: Ensure that Service Account keys are mounted into the Pod only
	                when AutoMountServiceToken is not set to false. We test the
	                following scenarios here.
	   1. Create Pod, Pod Spec has AutomountServiceAccountToken set to nil
	      a) Service Account with default value,
	      b) Service Account is an configured AutomountServiceAccountToken set to true,
	      c) Service Account is an configured AutomountServiceAccountToken set to false
	   2. Create Pod, Pod Spec has AutomountServiceAccountToken set to true
	      a) Service Account with default value,
	      b) Service Account is configured with AutomountServiceAccountToken set to true,
	      c) Service Account is configured with AutomountServiceAccountToken set to false
	   3. Create Pod, Pod Spec has AutomountServiceAccountToken set to false
	      a) Service Account with default value,
	      b) Service Account is configured with AutomountServiceAccountToken set to true,
	      c) Service Account is configured with AutomountServiceAccountToken set to false

	   The Containers running in these pods MUST verify that the ServiceTokenVolume path is
	   auto mounted only when Pod Spec has AutomountServiceAccountToken not set to false
	   and ServiceAccount object has AutomountServiceAccountToken not set to false, this
	   include test cases 1a,1b,2a,2b and 2c.
	   In the test cases 1c,3a,3b and 3c the ServiceTokenVolume MUST not be auto mounted.
	*/
	framework.ConformanceIt("should allow opting out of API token automount ", func() {

		var err error
		trueValue := true
		falseValue := false
		mountSA := &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "mount"}, AutomountServiceAccountToken: &trueValue}
		nomountSA := &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "nomount"}, AutomountServiceAccountToken: &falseValue}
		mountSA, err = f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Create(mountSA)
		framework.ExpectNoError(err)
		nomountSA, err = f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Create(nomountSA)
		framework.ExpectNoError(err)

		// Standard get, update retry loop
		framework.ExpectNoError(wait.Poll(time.Millisecond*500, framework.ServiceAccountProvisionTimeout, func() (bool, error) {
			By("getting the auto-created API token")
			sa, err := f.ClientSet.CoreV1().ServiceAccounts(f.Namespace.Name).Get(mountSA.Name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				framework.Logf("mount service account was not found")
				return false, nil
			}
			if err != nil {
				framework.Logf("error getting mount service account: %v", err)
				return false, err
			}
			if len(sa.Secrets) == 0 {
				framework.Logf("mount service account has no secret references")
				return false, nil
			}
			for _, secretRef := range sa.Secrets {
				secret, err := f.ClientSet.CoreV1().Secrets(f.Namespace.Name).Get(secretRef.Name, metav1.GetOptions{})
				if err != nil {
					framework.Logf("Error getting secret %s: %v", secretRef.Name, err)
					continue
				}
				if secret.Type == v1.SecretTypeServiceAccountToken {
					return true, nil
				}
			}

			framework.Logf("default service account has no secret references to valid service account tokens")
			return false, nil
		}))

		testcases := []struct {
			PodName            string
			ServiceAccountName string
			AutomountPodSpec   *bool
			ExpectTokenVolume  bool
		}{
			{
				PodName:            "pod-service-account-defaultsa",
				ServiceAccountName: "default",
				AutomountPodSpec:   nil,
				ExpectTokenVolume:  true, // default is true
			},
			{
				PodName:            "pod-service-account-mountsa",
				ServiceAccountName: mountSA.Name,
				AutomountPodSpec:   nil,
				ExpectTokenVolume:  true,
			},
			{
				PodName:            "pod-service-account-nomountsa",
				ServiceAccountName: nomountSA.Name,
				AutomountPodSpec:   nil,
				ExpectTokenVolume:  false,
			},

			// Make sure pod spec trumps when opting in
			{
				PodName:            "pod-service-account-defaultsa-mountspec",
				ServiceAccountName: "default",
				AutomountPodSpec:   &trueValue,
				ExpectTokenVolume:  true,
			},
			{
				PodName:            "pod-service-account-mountsa-mountspec",
				ServiceAccountName: mountSA.Name,
				AutomountPodSpec:   &trueValue,
				ExpectTokenVolume:  true,
			},
			{
				PodName:            "pod-service-account-nomountsa-mountspec",
				ServiceAccountName: nomountSA.Name,
				AutomountPodSpec:   &trueValue,
				ExpectTokenVolume:  true, // pod spec trumps
			},

			// Make sure pod spec trumps when opting out
			{
				PodName:            "pod-service-account-defaultsa-nomountspec",
				ServiceAccountName: "default",
				AutomountPodSpec:   &falseValue,
				ExpectTokenVolume:  false, // pod spec trumps
			},
			{
				PodName:            "pod-service-account-mountsa-nomountspec",
				ServiceAccountName: mountSA.Name,
				AutomountPodSpec:   &falseValue,
				ExpectTokenVolume:  false, // pod spec trumps
			},
			{
				PodName:            "pod-service-account-nomountsa-nomountspec",
				ServiceAccountName: nomountSA.Name,
				AutomountPodSpec:   &falseValue,
				ExpectTokenVolume:  false, // pod spec trumps
			},
		}

		for _, tc := range testcases {
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: tc.PodName},
				Spec: v1.PodSpec{
					Containers:                   []v1.Container{{Name: "token-test", Image: mountImage}},
					RestartPolicy:                v1.RestartPolicyNever,
					ServiceAccountName:           tc.ServiceAccountName,
					AutomountServiceAccountToken: tc.AutomountPodSpec,
				},
			}
			createdPod, err := f.ClientSet.CoreV1().Pods(f.Namespace.Name).Create(pod)
			framework.ExpectNoError(err)
			framework.Logf("created pod %s", tc.PodName)

			hasServiceAccountTokenVolume := false
			for _, c := range createdPod.Spec.Containers {
				for _, vm := range c.VolumeMounts {
					if vm.MountPath == serviceaccount.DefaultAPITokenMountPath {
						hasServiceAccountTokenVolume = true
					}
				}
			}

			if hasServiceAccountTokenVolume != tc.ExpectTokenVolume {
				framework.Failf("%s: expected volume=%v, got %v (%#v)", tc.PodName, tc.ExpectTokenVolume, hasServiceAccountTokenVolume, createdPod)
			} else {
				framework.Logf("pod %s service account token volume mount: %v", tc.PodName, hasServiceAccountTokenVolume)
			}
		}
	})
})
