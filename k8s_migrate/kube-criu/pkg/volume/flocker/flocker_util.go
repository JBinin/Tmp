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

package flocker

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	volutil "k8s.io/kubernetes/pkg/volume/util"

	flockerapi "github.com/clusterhq/flocker-go"
	"github.com/golang/glog"
)

type FlockerUtil struct{}

func (util *FlockerUtil) DeleteVolume(d *flockerVolumeDeleter) error {
	var err error

	if d.flockerClient == nil {
		d.flockerClient, err = d.plugin.newFlockerClient("")
		if err != nil {
			return err
		}
	}

	datasetUUID, err := d.GetDatasetUUID()
	if err != nil {
		return err
	}

	return d.flockerClient.DeleteDataset(datasetUUID)
}

func (util *FlockerUtil) CreateVolume(c *flockerVolumeProvisioner) (datasetUUID string, volumeSizeGiB int, labels map[string]string, err error) {

	if c.flockerClient == nil {
		c.flockerClient, err = c.plugin.newFlockerClient("")
		if err != nil {
			return
		}
	}

	nodes, err := c.flockerClient.ListNodes()
	if err != nil {
		return
	}
	if len(nodes) < 1 {
		err = fmt.Errorf("no nodes found inside the flocker cluster to provision a dataset")
		return
	}

	// select random node
	rand.Seed(time.Now().UTC().UnixNano())
	node := nodes[rand.Intn(len(nodes))]
	glog.V(2).Infof("selected flocker node with UUID '%s' to provision dataset", node.UUID)

	capacity := c.options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	requestBytes := capacity.Value()
	volumeSizeGiB, err = volutil.RoundUpToGiBInt(capacity)
	if err != nil {
		return
	}

	createOptions := &flockerapi.CreateDatasetOptions{
		MaximumSize: requestBytes,
		Metadata: map[string]string{
			"type": "k8s-dynamic-prov",
			"pvc":  c.options.PVC.Name,
		},
		Primary: node.UUID,
	}

	datasetState, err := c.flockerClient.CreateDataset(createOptions)
	if err != nil {
		return
	}
	datasetUUID = datasetState.DatasetID

	glog.V(2).Infof("successfully created Flocker dataset with UUID '%s'", datasetUUID)

	return
}
