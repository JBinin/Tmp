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

package factory

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	etcd2client "github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"

	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/etcd"
	"k8s.io/apiserver/pkg/storage/storagebackend"
)

func newETCD2HealthCheck(c storagebackend.Config) (func() error, error) {
	tr, err := newTransportForETCD2(c.CertFile, c.KeyFile, c.CAFile)
	if err != nil {
		return nil, err
	}

	client, err := newETCD2Client(tr, c.ServerList)
	if err != nil {
		return nil, err
	}

	members := etcd2client.NewMembersAPI(client)

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := members.List(ctx); err != nil {
			return fmt.Errorf("error listing etcd members: %v", err)
		}
		return nil
	}, nil
}

func newETCD2Storage(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {
	tr, err := newTransportForETCD2(c.CertFile, c.KeyFile, c.CAFile)
	if err != nil {
		return nil, nil, err
	}
	client, err := newETCD2Client(tr, c.ServerList)
	if err != nil {
		return nil, nil, err
	}
	s := etcd.NewEtcdStorage(client, c.Codec, c.Prefix, c.Quorum, c.DeserializationCacheSize, etcd.IdentityTransformer)
	return s, tr.CloseIdleConnections, nil
}

func newETCD2Client(tr *http.Transport, serverList []string) (etcd2client.Client, error) {
	cli, err := etcd2client.New(etcd2client.Config{
		Endpoints: serverList,
		Transport: tr,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func newTransportForETCD2(certFile, keyFile, caFile string) (*http.Transport, error) {
	info := transport.TLSInfo{
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}
	cfg, err := info.ClientConfig()
	if err != nil {
		return nil, err
	}
	// Copied from etcd.DefaultTransport declaration.
	// TODO: Determine if transport needs optimization
	tr := utilnet.SetTransportDefaults(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConnsPerHost: 500,
		TLSClientConfig:     cfg,
	})
	return tr, nil
}
