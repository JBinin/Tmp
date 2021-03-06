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

package vclib

import (
	"strings"

	"github.com/golang/glog"
)

// VolumeOptions specifies various options for a volume.
type VolumeOptions struct {
	CapacityKB             int
	Tags                   map[string]string
	Name                   string
	DiskFormat             string
	Datastore              string
	VSANStorageProfileData string
	StoragePolicyName      string
	StoragePolicyID        string
	SCSIControllerType     string
}

var (
	// DiskFormatValidType specifies the valid disk formats
	DiskFormatValidType = map[string]string{
		ThinDiskType:                              ThinDiskType,
		strings.ToLower(EagerZeroedThickDiskType): EagerZeroedThickDiskType,
		strings.ToLower(ZeroedThickDiskType):      PreallocatedDiskType,
	}
	// SCSIControllerValidType specifies the supported SCSI controllers
	SCSIControllerValidType = []string{LSILogicControllerType, LSILogicSASControllerType, PVSCSIControllerType}
)

// DiskformatValidOptions generates Valid Options for Diskformat
func DiskformatValidOptions() string {
	validopts := ""
	for diskformat := range DiskFormatValidType {
		validopts += diskformat + ", "
	}
	validopts = strings.TrimSuffix(validopts, ", ")
	return validopts
}

// CheckDiskFormatSupported checks if the diskFormat is valid
func CheckDiskFormatSupported(diskFormat string) bool {
	if DiskFormatValidType[diskFormat] == "" {
		glog.Errorf("Not a valid Disk Format. Valid options are %+q", DiskformatValidOptions())
		return false
	}
	return true
}

// SCSIControllerTypeValidOptions generates valid options for SCSIControllerType
func SCSIControllerTypeValidOptions() string {
	validopts := ""
	for _, controllerType := range SCSIControllerValidType {
		validopts += (controllerType + ", ")
	}
	validopts = strings.TrimSuffix(validopts, ", ")
	return validopts
}

// CheckControllerSupported checks if the given controller type is valid
func CheckControllerSupported(ctrlType string) bool {
	for _, c := range SCSIControllerValidType {
		if ctrlType == c {
			return true
		}
	}
	glog.Errorf("Not a valid SCSI Controller Type. Valid options are %q", SCSIControllerTypeValidOptions())
	return false
}

// VerifyVolumeOptions checks if volumeOptions.SCIControllerType is valid controller type
func (volumeOptions VolumeOptions) VerifyVolumeOptions() bool {
	// Validate only if SCSIControllerType is set by user.
	// Default value is set later in virtualDiskManager.Create and vmDiskManager.Create
	if volumeOptions.SCSIControllerType != "" {
		isValid := CheckControllerSupported(volumeOptions.SCSIControllerType)
		if !isValid {
			return false
		}
	}
	// ThinDiskType is the default, so skip the validation.
	if volumeOptions.DiskFormat != ThinDiskType {
		isValid := CheckDiskFormatSupported(volumeOptions.DiskFormat)
		if !isValid {
			return false
		}
	}
	return true
}
