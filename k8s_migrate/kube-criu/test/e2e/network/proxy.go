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

// OWNER = sig/network

package network

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/net"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	testutils "k8s.io/kubernetes/test/utils"
	imageutils "k8s.io/kubernetes/test/utils/image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	// Try all the proxy tests this many times (to catch even rare flakes).
	proxyAttempts = 20
	// Only print this many characters of the response (to keep the logs
	// legible).
	maxDisplayBodyLen = 100

	// We have seen one of these calls take just over 15 seconds, so putting this at 30.
	proxyHTTPCallTimeout = 30 * time.Second
)

var _ = SIGDescribe("Proxy", func() {
	version := "v1"
	Context("version "+version, func() {
		options := framework.FrameworkOptions{
			ClientQPS: -1.0,
		}
		f := framework.NewFramework("proxy", options, nil)
		prefix := "/api/" + version

		/*
			Release : v1.9
			Testname: Proxy, logs port endpoint
			Description: Select any node in the cluster to invoke /proxy/nodes/<nodeip>:10250/logs endpoint. This endpoint MUST be reachable.
		*/
		framework.ConformanceIt("should proxy logs on node with explicit kubelet port using proxy subresource ", func() { nodeProxyTest(f, prefix+"/nodes/", ":10250/proxy/logs/") })

		/*
			Release : v1.9
			Testname: Proxy, logs endpoint
			Description:  Select any node in the cluster to invoke /proxy/nodes/<nodeip>//logs endpoint. This endpoint MUST be reachable.
		*/
		framework.ConformanceIt("should proxy logs on node using proxy subresource ", func() { nodeProxyTest(f, prefix+"/nodes/", "/proxy/logs/") })

		// using the porter image to serve content, access the content
		// (of multiple pods?) from multiple (endpoints/services?)
		/*
			Release : v1.9
			Testname: Proxy, logs service endpoint
			Description: Select any node in the cluster to invoke  /logs endpoint  using the /nodes/proxy subresource from the kubelet port. This endpoint MUST be reachable.
		*/
		framework.ConformanceIt("should proxy through a service and a pod ", func() {
			start := time.Now()
			labels := map[string]string{"proxy-service-target": "true"}
			service, err := f.ClientSet.CoreV1().Services(f.Namespace.Name).Create(&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "proxy-service-",
				},
				Spec: v1.ServiceSpec{
					Selector: labels,
					Ports: []v1.ServicePort{
						{
							Name:       "portname1",
							Port:       80,
							TargetPort: intstr.FromString("dest1"),
						},
						{
							Name:       "portname2",
							Port:       81,
							TargetPort: intstr.FromInt(162),
						},
						{
							Name:       "tlsportname1",
							Port:       443,
							TargetPort: intstr.FromString("tlsdest1"),
						},
						{
							Name:       "tlsportname2",
							Port:       444,
							TargetPort: intstr.FromInt(462),
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Make an RC with a single pod. The 'porter' image is
			// a simple server which serves the values of the
			// environmental variables below.
			By("starting an echo server on multiple ports")
			pods := []*v1.Pod{}
			cfg := testutils.RCConfig{
				Client:         f.ClientSet,
				InternalClient: f.InternalClientset,
				Image:          imageutils.GetE2EImage(imageutils.Porter),
				Name:           service.Name,
				Namespace:      f.Namespace.Name,
				Replicas:       1,
				PollInterval:   time.Second,
				Env: map[string]string{
					"SERVE_PORT_80":   `<a href="/rewriteme">test</a>`,
					"SERVE_PORT_1080": `<a href="/rewriteme">test</a>`,
					"SERVE_PORT_160":  "foo",
					"SERVE_PORT_162":  "bar",

					"SERVE_TLS_PORT_443": `<a href="/tlsrewriteme">test</a>`,
					"SERVE_TLS_PORT_460": `tls baz`,
					"SERVE_TLS_PORT_462": `tls qux`,
				},
				Ports: map[string]int{
					"dest1": 160,
					"dest2": 162,

					"tlsdest1": 460,
					"tlsdest2": 462,
				},
				ReadinessProbe: &v1.Probe{
					Handler: v1.Handler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(80),
						},
					},
					InitialDelaySeconds: 1,
					TimeoutSeconds:      5,
					PeriodSeconds:       10,
				},
				Labels:      labels,
				CreatedPods: &pods,
			}
			Expect(framework.RunRC(cfg)).NotTo(HaveOccurred())
			defer framework.DeleteRCAndWaitForGC(f.ClientSet, f.Namespace.Name, cfg.Name)

			Expect(framework.WaitForEndpoint(f.ClientSet, f.Namespace.Name, service.Name)).NotTo(HaveOccurred())

			// table constructors
			// Try proxying through the service and directly to through the pod.
			subresourceServiceProxyURL := func(scheme, port string) string {
				return prefix + "/namespaces/" + f.Namespace.Name + "/services/" + net.JoinSchemeNamePort(scheme, service.Name, port) + "/proxy"
			}
			subresourcePodProxyURL := func(scheme, port string) string {
				return prefix + "/namespaces/" + f.Namespace.Name + "/pods/" + net.JoinSchemeNamePort(scheme, pods[0].Name, port) + "/proxy"
			}

			// construct the table
			expectations := map[string]string{
				subresourceServiceProxyURL("", "portname1") + "/":         "foo",
				subresourceServiceProxyURL("http", "portname1") + "/":     "foo",
				subresourceServiceProxyURL("", "portname2") + "/":         "bar",
				subresourceServiceProxyURL("http", "portname2") + "/":     "bar",
				subresourceServiceProxyURL("https", "tlsportname1") + "/": "tls baz",
				subresourceServiceProxyURL("https", "tlsportname2") + "/": "tls qux",

				subresourcePodProxyURL("", "") + "/":         `<a href="` + subresourcePodProxyURL("", "") + `/rewriteme">test</a>`,
				subresourcePodProxyURL("", "1080") + "/":     `<a href="` + subresourcePodProxyURL("", "1080") + `/rewriteme">test</a>`,
				subresourcePodProxyURL("http", "1080") + "/": `<a href="` + subresourcePodProxyURL("http", "1080") + `/rewriteme">test</a>`,
				subresourcePodProxyURL("", "160") + "/":      "foo",
				subresourcePodProxyURL("http", "160") + "/":  "foo",
				subresourcePodProxyURL("", "162") + "/":      "bar",
				subresourcePodProxyURL("http", "162") + "/":  "bar",

				subresourcePodProxyURL("https", "443") + "/": `<a href="` + subresourcePodProxyURL("https", "443") + `/tlsrewriteme">test</a>`,
				subresourcePodProxyURL("https", "460") + "/": "tls baz",
				subresourcePodProxyURL("https", "462") + "/": "tls qux",

				// TODO: below entries don't work, but I believe we should make them work.
				// podPrefix + ":dest1": "foo",
				// podPrefix + ":dest2": "bar",
			}

			wg := sync.WaitGroup{}
			errs := []string{}
			errLock := sync.Mutex{}
			recordError := func(s string) {
				errLock.Lock()
				defer errLock.Unlock()
				errs = append(errs, s)
			}
			d := time.Since(start)
			framework.Logf("setup took %v, starting test cases", d)
			numberTestCases := len(expectations)
			totalAttempts := numberTestCases * proxyAttempts
			By(fmt.Sprintf("running %v cases, %v attempts per case, %v total attempts", numberTestCases, proxyAttempts, totalAttempts))

			for i := 0; i < proxyAttempts; i++ {
				wg.Add(numberTestCases)
				for path, val := range expectations {
					go func(i int, path, val string) {
						defer wg.Done()
						// this runs the test case
						body, status, d, err := doProxy(f, path, i)

						if err != nil {
							if serr, ok := err.(*errors.StatusError); ok {
								recordError(fmt.Sprintf("%v (%v; %v): path %v gave status error: %+v",
									i, status, d, path, serr.Status()))
							} else {
								recordError(fmt.Sprintf("%v: path %v gave error: %v", i, path, err))
							}
							return
						}
						if status != http.StatusOK {
							recordError(fmt.Sprintf("%v: path %v gave status: %v", i, path, status))
						}
						if e, a := val, string(body); e != a {
							recordError(fmt.Sprintf("%v: path %v: wanted %v, got %v", i, path, e, a))
						}
						if d > proxyHTTPCallTimeout {
							recordError(fmt.Sprintf("%v: path %v took %v > %v", i, path, d, proxyHTTPCallTimeout))
						}
					}(i, path, val)
				}
				wg.Wait()
			}

			if len(errs) != 0 {
				body, err := f.ClientSet.CoreV1().Pods(f.Namespace.Name).GetLogs(pods[0].Name, &v1.PodLogOptions{}).Do().Raw()
				if err != nil {
					framework.Logf("Error getting logs for pod %s: %v", pods[0].Name, err)
				} else {
					framework.Logf("Pod %s has the following error logs: %s", pods[0].Name, body)
				}

				framework.Failf(strings.Join(errs, "\n"))
			}
		})
	})
})

func doProxy(f *framework.Framework, path string, i int) (body []byte, statusCode int, d time.Duration, err error) {
	// About all of the proxy accesses in this file:
	// * AbsPath is used because it preserves the trailing '/'.
	// * Do().Raw() is used (instead of DoRaw()) because it will turn an
	//   error from apiserver proxy into an actual error, and there is no
	//   chance of the things we are talking to being confused for an error
	//   that apiserver would have emitted.
	start := time.Now()
	body, err = f.ClientSet.CoreV1().RESTClient().Get().AbsPath(path).Do().StatusCode(&statusCode).Raw()
	d = time.Since(start)
	if len(body) > 0 {
		framework.Logf("(%v) %v: %s (%v; %v)", i, path, truncate(body, maxDisplayBodyLen), statusCode, d)
	} else {
		framework.Logf("%v: %s (%v; %v)", path, "no body", statusCode, d)
	}
	return
}

func truncate(b []byte, maxLen int) []byte {
	if len(b) <= maxLen-3 {
		return b
	}
	b2 := append([]byte(nil), b[:maxLen-3]...)
	b2 = append(b2, '.', '.', '.')
	return b2
}

func pickNode(cs clientset.Interface) (string, error) {
	// TODO: investigate why it doesn't work on master Node.
	nodes := framework.GetReadySchedulableNodesOrDie(cs)
	if len(nodes.Items) == 0 {
		return "", fmt.Errorf("no nodes exist, can't test node proxy")
	}
	return nodes.Items[0].Name, nil
}

func nodeProxyTest(f *framework.Framework, prefix, nodeDest string) {
	node, err := pickNode(f.ClientSet)
	Expect(err).NotTo(HaveOccurred())
	// TODO: Change it to test whether all requests succeeded when requests
	// not reaching Kubelet issue is debugged.
	serviceUnavailableErrors := 0
	for i := 0; i < proxyAttempts; i++ {
		_, status, d, err := doProxy(f, prefix+node+nodeDest, i)
		if status == http.StatusServiceUnavailable {
			framework.Logf("Failed proxying node logs due to service unavailable: %v", err)
			time.Sleep(time.Second)
			serviceUnavailableErrors++
		} else {
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(d).To(BeNumerically("<", proxyHTTPCallTimeout))
		}
	}
	if serviceUnavailableErrors > 0 {
		framework.Logf("error: %d requests to proxy node logs failed", serviceUnavailableErrors)
	}
	maxFailures := int(math.Floor(0.1 * float64(proxyAttempts)))
	Expect(serviceUnavailableErrors).To(BeNumerically("<", maxFailures))
}
