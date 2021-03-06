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

package main

import (
	goflag "flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	cadvisortest "k8s.io/kubernetes/pkg/kubelet/cadvisor/testing"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	"k8s.io/kubernetes/pkg/kubelet/dockershim"
	"k8s.io/kubernetes/pkg/kubelet/dockershim/libdocker"
	"k8s.io/kubernetes/pkg/kubemark"
	fakeiptables "k8s.io/kubernetes/pkg/util/iptables/testing"
	fakesysctl "k8s.io/kubernetes/pkg/util/sysctl/testing"
	_ "k8s.io/kubernetes/pkg/version/prometheus" // for version metric registration
	"k8s.io/kubernetes/pkg/version/verflag"
	fakeexec "k8s.io/utils/exec/testing"
)

type HollowNodeConfig struct {
	KubeconfigPath       string
	KubeletPort          int
	KubeletReadOnlyPort  int
	Morph                string
	NodeName             string
	ServerPort           int
	ContentType          string
	UseRealProxier       bool
	ProxierSyncPeriod    time.Duration
	ProxierMinSyncPeriod time.Duration
}

const (
	maxPods     = 110
	podsPerCore = 0
)

// TODO(#45650): Refactor hollow-node into hollow-kubelet and hollow-proxy
// and make the config driven.
var knownMorphs = sets.NewString("kubelet", "proxy")

func (c *HollowNodeConfig) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.KubeconfigPath, "kubeconfig", "/kubeconfig/kubeconfig", "Path to kubeconfig file.")
	fs.IntVar(&c.KubeletPort, "kubelet-port", 10250, "Port on which HollowKubelet should be listening.")
	fs.IntVar(&c.KubeletReadOnlyPort, "kubelet-read-only-port", 10255, "Read-only port on which Kubelet is listening.")
	fs.StringVar(&c.NodeName, "name", "fake-node", "Name of this Hollow Node.")
	fs.IntVar(&c.ServerPort, "api-server-port", 443, "Port on which API server is listening.")
	fs.StringVar(&c.Morph, "morph", "", fmt.Sprintf("Specifies into which Hollow component this binary should morph. Allowed values: %v", knownMorphs.List()))
	fs.StringVar(&c.ContentType, "kube-api-content-type", "application/vnd.kubernetes.protobuf", "ContentType of requests sent to apiserver.")
	fs.BoolVar(&c.UseRealProxier, "use-real-proxier", true, "Set to true if you want to use real proxier inside hollow-proxy.")
	fs.DurationVar(&c.ProxierSyncPeriod, "proxier-sync-period", 30*time.Second, "Period that proxy rules are refreshed in hollow-proxy.")
	fs.DurationVar(&c.ProxierMinSyncPeriod, "proxier-min-sync-period", 0, "Minimum period that proxy rules are refreshed in hollow-proxy.")
}

func (c *HollowNodeConfig) createClientConfigFromFile() (*restclient.Config, error) {
	clientConfig, err := clientcmd.LoadFromFile(c.KubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error while loading kubeconfig from file %v: %v", c.KubeconfigPath, err)
	}
	config, err := clientcmd.NewDefaultClientConfig(*clientConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error while creating kubeconfig: %v", err)
	}
	config.ContentType = c.ContentType
	config.QPS = 10
	config.Burst = 20
	return config, nil
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	command := newHollowNodeCommand()

	// TODO: once we switch everything over to Cobra commands, we can go back to calling
	// utilflag.InitFlags() (by removing its pflag.Parse() call). For now, we have to set the
	// normalize func and add the go flag set by hand.
	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	// utilflag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// newControllerManagerCommand creates a *cobra.Command object with default parameters
func newHollowNodeCommand() *cobra.Command {
	s := &HollowNodeConfig{}

	cmd := &cobra.Command{
		Use:  "kubemark",
		Long: "kubemark",
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()
			run(s)
		},
	}
	s.addFlags(cmd.Flags())

	return cmd
}

func run(config *HollowNodeConfig) {
	if !knownMorphs.Has(config.Morph) {
		glog.Fatalf("Unknown morph: %v. Allowed values: %v", config.Morph, knownMorphs.List())
	}

	// create a client to communicate with API server.
	clientConfig, err := config.createClientConfigFromFile()
	if err != nil {
		glog.Fatalf("Failed to create a ClientConfig: %v. Exiting.", err)
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		glog.Fatalf("Failed to create a ClientSet: %v. Exiting.", err)
	}

	if config.Morph == "kubelet" {
		cadvisorInterface := &cadvisortest.Fake{
			NodeName: config.NodeName,
		}
		containerManager := cm.NewStubContainerManager()

		fakeDockerClientConfig := &dockershim.ClientConfig{
			DockerEndpoint:    libdocker.FakeDockerEndpoint,
			EnableSleep:       true,
			WithTraceDisabled: true,
		}

		hollowKubelet := kubemark.NewHollowKubelet(
			config.NodeName,
			client,
			cadvisorInterface,
			fakeDockerClientConfig,
			config.KubeletPort,
			config.KubeletReadOnlyPort,
			containerManager,
			maxPods,
			podsPerCore,
		)
		hollowKubelet.Run()
	}

	if config.Morph == "proxy" {
		client, err := clientset.NewForConfig(clientConfig)
		if err != nil {
			glog.Fatalf("Failed to create API Server client: %v", err)
		}
		iptInterface := fakeiptables.NewFake()
		sysctl := fakesysctl.NewFake()
		execer := &fakeexec.FakeExec{}
		eventBroadcaster := record.NewBroadcaster()
		recorder := eventBroadcaster.NewRecorder(legacyscheme.Scheme, v1.EventSource{Component: "kube-proxy", Host: config.NodeName})

		hollowProxy, err := kubemark.NewHollowProxyOrDie(
			config.NodeName,
			client,
			client.CoreV1(),
			iptInterface,
			sysctl,
			execer,
			eventBroadcaster,
			recorder,
			config.UseRealProxier,
			config.ProxierSyncPeriod,
			config.ProxierMinSyncPeriod,
		)
		if err != nil {
			glog.Fatalf("Failed to create hollowProxy instance: %v", err)
		}
		hollowProxy.Run()
	}
}
