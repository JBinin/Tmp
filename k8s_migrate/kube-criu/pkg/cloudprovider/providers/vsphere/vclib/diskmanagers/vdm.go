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

package diskmanagers

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/vsphere/vclib"
)

// virtualDiskManager implements VirtualDiskProvider Interface for creating and deleting volume using VirtualDiskManager
type virtualDiskManager struct {
	diskPath      string
	volumeOptions *vclib.VolumeOptions
}

// Create implements Disk's Create interface
// Contains implementation of virtualDiskManager based Provisioning
func (diskManager virtualDiskManager) Create(ctx context.Context, datastore *vclib.Datastore) (canonicalDiskPath string, err error) {
	if diskManager.volumeOptions.SCSIControllerType == "" {
		diskManager.volumeOptions.SCSIControllerType = vclib.LSILogicControllerType
	}
	// Create virtual disk
	diskFormat := vclib.DiskFormatValidType[diskManager.volumeOptions.DiskFormat]
	// Create a virtual disk manager
	vdm := object.NewVirtualDiskManager(datastore.Client())
	// Create specification for new virtual disk
	vmDiskSpec := &types.FileBackedVirtualDiskSpec{
		VirtualDiskSpec: types.VirtualDiskSpec{
			AdapterType: diskManager.volumeOptions.SCSIControllerType,
			DiskType:    diskFormat,
		},
		CapacityKb: int64(diskManager.volumeOptions.CapacityKB),
	}
	requestTime := time.Now()
	// Create virtual disk
	task, err := vdm.CreateVirtualDisk(ctx, diskManager.diskPath, datastore.Datacenter.Datacenter, vmDiskSpec)
	if err != nil {
		vclib.RecordvSphereMetric(vclib.APICreateVolume, requestTime, err)
		glog.Errorf("Failed to create virtual disk: %s. err: %+v", diskManager.diskPath, err)
		return "", err
	}
	taskInfo, err := task.WaitForResult(ctx, nil)
	vclib.RecordvSphereMetric(vclib.APICreateVolume, requestTime, err)
	if err != nil {
		glog.Errorf("Failed to complete virtual disk creation: %s. err: %+v", diskManager.diskPath, err)
		return "", err
	}
	canonicalDiskPath = taskInfo.Result.(string)
	return canonicalDiskPath, nil
}

// Delete implements Disk's Delete interface
func (diskManager virtualDiskManager) Delete(ctx context.Context, datacenter *vclib.Datacenter) error {
	// Create a virtual disk manager
	virtualDiskManager := object.NewVirtualDiskManager(datacenter.Client())
	diskPath := vclib.RemoveStorageClusterORFolderNameFromVDiskPath(diskManager.diskPath)
	requestTime := time.Now()
	// Delete virtual disk
	task, err := virtualDiskManager.DeleteVirtualDisk(ctx, diskPath, datacenter.Datacenter)
	if err != nil {
		glog.Errorf("Failed to delete virtual disk. err: %v", err)
		vclib.RecordvSphereMetric(vclib.APIDeleteVolume, requestTime, err)
		return err
	}
	err = task.Wait(ctx)
	vclib.RecordvSphereMetric(vclib.APIDeleteVolume, requestTime, err)
	if err != nil {
		glog.Errorf("Failed to delete virtual disk. err: %v", err)
		return err
	}
	return nil
}
