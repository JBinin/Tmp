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

package common

import (
	. "github.com/onsi/ginkgo"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/test/e2e/framework"
)

var _ = Describe("[sig-network] Networking", func() {
	f := framework.NewDefaultFramework("pod-network-test")

	Describe("Granular Checks: Pods", func() {

		// Try to hit all endpoints through a test container, retry 5 times,
		// expect exactly one unique hostname. Each of these endpoints reports
		// its own hostname.
		/*
			Release : v1.9
			Testname: Networking, intra pod http
			Description: Create a hostexec pod that is capable of curl to netcat commands. Create a test Pod that will act as a webserver front end exposing ports 8080 for tcp and 8081 for udp. The netserver service proxies are created on specified number of nodes.
			The kubectl exec on the webserver container MUST reach a http port on the each of service proxy endpoints in the cluster and the request MUST be successful. Container will execute curl command to reach the service port within specified max retry limit and MUST result in reporting unique hostnames.
		*/
		framework.ConformanceIt("should function for intra-pod communication: http [NodeConformance]", func() {
			config := framework.NewCoreNetworkingTestConfig(f)
			for _, endpointPod := range config.EndpointPods {
				config.DialFromTestContainer("http", endpointPod.Status.PodIP, framework.EndpointHttpPort, config.MaxTries, 0, sets.NewString(endpointPod.Name))
			}
		})

		/*
			Release : v1.9
			Testname: Networking, intra pod udp
			Description: Create a hostexec pod that is capable of curl to netcat commands. Create a test Pod that will act as a webserver front end exposing ports 8080 for tcp and 8081 for udp. The netserver service proxies are created on specified number of nodes.
			The kubectl exec on the webserver container MUST reach a udp port on the each of service proxy endpoints in the cluster and the request MUST be successful. Container will execute curl command to reach the service port within specified max retry limit and MUST result in reporting unique hostnames.
		*/
		framework.ConformanceIt("should function for intra-pod communication: udp [NodeConformance]", func() {
			config := framework.NewCoreNetworkingTestConfig(f)
			for _, endpointPod := range config.EndpointPods {
				config.DialFromTestContainer("udp", endpointPod.Status.PodIP, framework.EndpointUdpPort, config.MaxTries, 0, sets.NewString(endpointPod.Name))
			}
		})

		/*
			Release : v1.9
			Testname: Networking, intra pod http, from node
			Description: Create a hostexec pod that is capable of curl to netcat commands. Create a test Pod that will act as a webserver front end exposing ports 8080 for tcp and 8081 for udp. The netserver service proxies are created on specified number of nodes.
			The kubectl exec on the webserver container MUST reach a http port on the each of service proxy endpoints in the cluster using a http post(protocol=tcp)  and the request MUST be successful. Container will execute curl command to reach the service port within specified max retry limit and MUST result in reporting unique hostnames.
		*/
		framework.ConformanceIt("should function for node-pod communication: http [NodeConformance]", func() {
			config := framework.NewCoreNetworkingTestConfig(f)
			for _, endpointPod := range config.EndpointPods {
				config.DialFromNode("http", endpointPod.Status.PodIP, framework.EndpointHttpPort, config.MaxTries, 0, sets.NewString(endpointPod.Name))
			}
		})

		/*
			Release : v1.9
			Testname: Networking, intra pod http, from node
			Description: Create a hostexec pod that is capable of curl to netcat commands. Create a test Pod that will act as a webserver front end exposing ports 8080 for tcp and 8081 for udp. The netserver service proxies are created on specified number of nodes.
			The kubectl exec on the webserver container MUST reach a http port on the each of service proxy endpoints in the cluster using a http post(protocol=udp)  and the request MUST be successful. Container will execute curl command to reach the service port within specified max retry limit and MUST result in reporting unique hostnames.
		*/
		framework.ConformanceIt("should function for node-pod communication: udp [NodeConformance]", func() {
			config := framework.NewCoreNetworkingTestConfig(f)
			for _, endpointPod := range config.EndpointPods {
				config.DialFromNode("udp", endpointPod.Status.PodIP, framework.EndpointUdpPort, config.MaxTries, 0, sets.NewString(endpointPod.Name))
			}
		})
	})
})
