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

package https

import (
	"io/ioutil"
	"net/http"

	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubernetes/cmd/kubeadm/app/discovery/file"
)

// RetrieveValidatedConfigInfo connects to the API Server and makes sure it can talk
// securely to the API Server using the provided CA cert and
// optionally refreshes the cluster-info information from the cluster-info ConfigMap
func RetrieveValidatedConfigInfo(httpsURL, clustername string) (*clientcmdapi.Config, error) {
	client := &http.Client{Transport: netutil.SetOldTransportDefaults(&http.Transport{})}
	response, err := client.Get(httpsURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	kubeconfig, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return nil, err
	}
	return file.ValidateConfigInfo(config, clustername)
}
