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

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientapi "k8s.io/client-go/tools/clientcmd/api"
	api "k8s.io/kubernetes/pkg/apis/core"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

func buildConfigFromEnvs(masterURL, kubeconfigPath string) (*restclient.Config, error) {
	if kubeconfigPath == "" && masterURL == "" {
		kubeconfig, err := restclient.InClusterConfig()
		if err != nil {
			return nil, err
		}

		return kubeconfig, nil
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientapi.Cluster{Server: masterURL}}).ClientConfig()
}

func flattenSubsets(subsets []api.EndpointSubset) []string {
	ips := []string{}
	for _, ss := range subsets {
		for _, addr := range ss.Addresses {
			ips = append(ips, fmt.Sprintf(`"%s"`, addr.IP))
		}
	}
	return ips
}

func main() {
	flag.Parse()

	glog.Info("Kubernetes Elasticsearch logging discovery")

	cc, err := buildConfigFromEnvs(os.Getenv("APISERVER_HOST"), os.Getenv("KUBE_CONFIG_FILE"))
	if err != nil {
		glog.Fatalf("Failed to make client: %v", err)
	}
	client, err := clientset.NewForConfig(cc)

	if err != nil {
		glog.Fatalf("Failed to make client: %v", err)
	}
	namespace := metav1.NamespaceSystem
	envNamespace := os.Getenv("NAMESPACE")
	if envNamespace != "" {
		if _, err := client.Core().Namespaces().Get(envNamespace, metav1.GetOptions{}); err != nil {
			glog.Fatalf("%s namespace doesn't exist: %v", envNamespace, err)
		}
		namespace = envNamespace
	}

	var elasticsearch *api.Service
	serviceName := os.Getenv("ELASTICSEARCH_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "elasticsearch-logging"
	}

	// Look for endpoints associated with the Elasticsearch logging service.
	// First wait for the service to become available.
	for t := time.Now(); time.Since(t) < 5*time.Minute; time.Sleep(10 * time.Second) {
		elasticsearch, err = client.Core().Services(namespace).Get(serviceName, metav1.GetOptions{})
		if err == nil {
			break
		}
	}
	// If we did not find an elasticsearch logging service then log a warning
	// and return without adding any unicast hosts.
	if elasticsearch == nil {
		glog.Warningf("Failed to find the elasticsearch-logging service: %v", err)
		return
	}

	var endpoints *api.Endpoints
	addrs := []string{}
	// Wait for some endpoints.
	count, _ := strconv.Atoi(os.Getenv("MINIMUM_MASTER_NODES"))
	for t := time.Now(); time.Since(t) < 5*time.Minute; time.Sleep(10 * time.Second) {
		endpoints, err = client.Core().Endpoints(namespace).Get(serviceName, metav1.GetOptions{})
		if err != nil {
			continue
		}
		addrs = flattenSubsets(endpoints.Subsets)
		glog.Infof("Found %s", addrs)
		if len(addrs) > 0 && len(addrs) == count {
			break
		}
	}
	// If there was an error finding endpoints then log a warning and quit.
	if err != nil {
		glog.Warningf("Error finding endpoints: %v", err)
		return
	}

	glog.Infof("Endpoints = %s", addrs)
	fmt.Printf("discovery.zen.ping.unicast.hosts: [%s]\n", strings.Join(addrs, ", "))
}
