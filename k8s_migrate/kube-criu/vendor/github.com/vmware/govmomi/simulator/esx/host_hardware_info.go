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
Copyright (c) 2017 VMware, Inc. All Rights Reserved.

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

package esx

import (
	"time"

	"github.com/vmware/govmomi/vim25/types"
)

// HostHardwareInfo is the default template for the HostSystem hardware property.
// Capture method:
//   govc object.collect -s -dump HostSystem:ha-host hardware
var HostHardwareInfo = &types.HostHardwareInfo{
	SystemInfo: types.HostSystemInfo{
		Vendor: "VMware, Inc.",
		Model:  "VMware Virtual Platform",
		Uuid:   "e88d4d56-9f1e-3ea1-71fa-13a8e1a7fd70",
		OtherIdentifyingInfo: []types.HostSystemIdentificationInfo{
			{
				IdentifierValue: " No Asset Tag",
				IdentifierType: &types.ElementDescription{
					Description: types.Description{
						Label:   "Asset Tag",
						Summary: "Asset tag of the system",
					},
					Key: "AssetTag",
				},
			},
			{
				IdentifierValue: "[MS_VM_CERT/SHA1/27d66596a61c48dd3dc7216fd715126e33f59ae7]",
				IdentifierType: &types.ElementDescription{
					Description: types.Description{
						Label:   "OEM specific string",
						Summary: "OEM specific string",
					},
					Key: "OemSpecificString",
				},
			},
			{
				IdentifierValue: "Welcome to the Virtual Machine",
				IdentifierType: &types.ElementDescription{
					Description: types.Description{
						Label:   "OEM specific string",
						Summary: "OEM specific string",
					},
					Key: "OemSpecificString",
				},
			},
			{
				IdentifierValue: "VMware-56 4d 8d e8 1e 9f a1 3e-71 fa 13 a8 e1 a7 fd 70",
				IdentifierType: &types.ElementDescription{
					Description: types.Description{
						Label:   "Service tag",
						Summary: "Service tag of the system",
					},
					Key: "ServiceTag",
				},
			},
		},
	},
	CpuPowerManagementInfo: &types.HostCpuPowerManagementInfo{
		CurrentPolicy:   "Balanced",
		HardwareSupport: "",
	},
	CpuInfo: types.HostCpuInfo{
		NumCpuPackages: 2,
		NumCpuCores:    2,
		NumCpuThreads:  2,
		Hz:             3591345000,
	},
	CpuPkg: []types.HostCpuPackage{
		{
			Index:       0,
			Vendor:      "intel",
			Hz:          3591345000,
			BusHz:       115849838,
			Description: "Intel(R) Xeon(R) CPU E5-1620 0 @ 3.60GHz",
			ThreadId:    []int16{0},
			CpuFeature: []types.HostCpuIdInfo{
				{
					Level:  0,
					Vendor: "",
					Eax:    "0000:0000:0000:0000:0000:0000:0000:1101",
					Ebx:    "0111:0101:0110:1110:0110:0101:0100:0111",
					Ecx:    "0110:1100:0110:0101:0111:0100:0110:1110",
					Edx:    "0100:1001:0110:0101:0110:1110:0110:1001",
				},
				{
					Level:  1,
					Vendor: "",
					Eax:    "0000:0000:0000:0010:0000:0110:1101:0111",
					Ebx:    "0000:0000:0000:0001:0000:1000:0000:0000",
					Ecx:    "1001:0111:1011:1010:0010:0010:0010:1011",
					Edx:    "0000:1111:1010:1011:1111:1011:1111:1111",
				},
				{
					Level:  -2147483648,
					Vendor: "",
					Eax:    "1000:0000:0000:0000:0000:0000:0000:1000",
					Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ecx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Edx:    "0000:0000:0000:0000:0000:0000:0000:0000",
				},
				{
					Level:  -2147483647,
					Vendor: "",
					Eax:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ecx:    "0000:0000:0000:0000:0000:0000:0000:0001",
					Edx:    "0010:1000:0001:0000:0000:1000:0000:0000",
				},
				{
					Level:  -2147483640,
					Vendor: "",
					Eax:    "0000:0000:0000:0000:0011:0000:0010:1010",
					Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ecx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Edx:    "0000:0000:0000:0000:0000:0000:0000:0000",
				},
			},
		},
		{
			Index:       1,
			Vendor:      "intel",
			Hz:          3591345000,
			BusHz:       115849838,
			Description: "Intel(R) Xeon(R) CPU E5-1620 0 @ 3.60GHz",
			ThreadId:    []int16{1},
			CpuFeature: []types.HostCpuIdInfo{
				{
					Level:  0,
					Vendor: "",
					Eax:    "0000:0000:0000:0000:0000:0000:0000:1101",
					Ebx:    "0111:0101:0110:1110:0110:0101:0100:0111",
					Ecx:    "0110:1100:0110:0101:0111:0100:0110:1110",
					Edx:    "0100:1001:0110:0101:0110:1110:0110:1001",
				},
				{
					Level:  1,
					Vendor: "",
					Eax:    "0000:0000:0000:0010:0000:0110:1101:0111",
					Ebx:    "0000:0010:0000:0001:0000:1000:0000:0000",
					Ecx:    "1001:0111:1011:1010:0010:0010:0010:1011",
					Edx:    "0000:1111:1010:1011:1111:1011:1111:1111",
				},
				{
					Level:  -2147483648,
					Vendor: "",
					Eax:    "1000:0000:0000:0000:0000:0000:0000:1000",
					Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ecx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Edx:    "0000:0000:0000:0000:0000:0000:0000:0000",
				},
				{
					Level:  -2147483647,
					Vendor: "",
					Eax:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ecx:    "0000:0000:0000:0000:0000:0000:0000:0001",
					Edx:    "0010:1000:0001:0000:0000:1000:0000:0000",
				},
				{
					Level:  -2147483640,
					Vendor: "",
					Eax:    "0000:0000:0000:0000:0011:0000:0010:1010",
					Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Ecx:    "0000:0000:0000:0000:0000:0000:0000:0000",
					Edx:    "0000:0000:0000:0000:0000:0000:0000:0000",
				},
			},
		},
	},
	MemorySize: 4294430720,
	NumaInfo: &types.HostNumaInfo{
		Type:     "NUMA",
		NumNodes: 1,
		NumaNode: []types.HostNumaNode{
			{
				TypeId:            0x0,
				CpuID:             []int16{1, 0},
				MemoryRangeBegin:  4294967296,
				MemoryRangeLength: 1073741824,
			},
		},
	},
	SmcPresent: types.NewBool(false),
	PciDevice: []types.HostPciDevice{
		{
			Id:           "0000:00:00.0",
			ClassId:      1536,
			Bus:          0x0,
			Slot:         0x0,
			Function:     0x0,
			VendorId:     -32634,
			SubVendorId:  5549,
			VendorName:   "Intel Corporation",
			DeviceId:     29072,
			SubDeviceId:  6518,
			ParentBridge: "",
			DeviceName:   "Virtual Machine Chipset",
		},
		{
			Id:           "0000:00:01.0",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x1,
			Function:     0x0,
			VendorId:     -32634,
			SubVendorId:  0,
			VendorName:   "Intel Corporation",
			DeviceId:     29073,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "440BX/ZX/DX - 82443BX/ZX/DX AGP bridge",
		},
		{
			Id:           "0000:00:07.0",
			ClassId:      1537,
			Bus:          0x0,
			Slot:         0x7,
			Function:     0x0,
			VendorId:     -32634,
			SubVendorId:  5549,
			VendorName:   "Intel Corporation",
			DeviceId:     28944,
			SubDeviceId:  6518,
			ParentBridge: "",
			DeviceName:   "Virtual Machine Chipset",
		},
		{
			Id:           "0000:00:07.1",
			ClassId:      257,
			Bus:          0x0,
			Slot:         0x7,
			Function:     0x1,
			VendorId:     -32634,
			SubVendorId:  5549,
			VendorName:   "Intel Corporation",
			DeviceId:     28945,
			SubDeviceId:  6518,
			ParentBridge: "",
			DeviceName:   "PIIX4 for 430TX/440BX/MX IDE Controller",
		},
		{
			Id:           "0000:00:07.3",
			ClassId:      1664,
			Bus:          0x0,
			Slot:         0x7,
			Function:     0x3,
			VendorId:     -32634,
			SubVendorId:  5549,
			VendorName:   "Intel Corporation",
			DeviceId:     28947,
			SubDeviceId:  6518,
			ParentBridge: "",
			DeviceName:   "Virtual Machine Chipset",
		},
		{
			Id:           "0000:00:07.7",
			ClassId:      2176,
			Bus:          0x0,
			Slot:         0x7,
			Function:     0x7,
			VendorId:     5549,
			SubVendorId:  5549,
			VendorName:   "VMware",
			DeviceId:     1856,
			SubDeviceId:  1856,
			ParentBridge: "",
			DeviceName:   "Virtual Machine Communication Interface",
		},
		{
			Id:           "0000:00:0f.0",
			ClassId:      768,
			Bus:          0x0,
			Slot:         0xf,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  5549,
			VendorName:   "VMware",
			DeviceId:     1029,
			SubDeviceId:  1029,
			ParentBridge: "",
			DeviceName:   "SVGA II Adapter",
		},
		{
			Id:           "0000:00:11.0",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x11,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1936,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI bridge",
		},
		{
			Id:           "0000:00:15.0",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:15.1",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x1,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:15.2",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x2,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:15.3",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x3,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:15.4",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x4,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:15.5",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x5,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:15.6",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x6,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:15.7",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x15,
			Function:     0x7,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.0",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.1",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x1,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.2",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x2,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.3",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x3,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.4",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x4,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.5",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x5,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.6",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x6,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:16.7",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x16,
			Function:     0x7,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.0",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.1",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x1,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.2",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x2,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.3",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x3,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.4",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x4,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.5",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x5,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.6",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x6,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:17.7",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x17,
			Function:     0x7,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.0",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.1",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x1,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.2",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x2,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.3",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x3,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.4",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x4,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.5",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x5,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.6",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x6,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:00:18.7",
			ClassId:      1540,
			Bus:          0x0,
			Slot:         0x18,
			Function:     0x7,
			VendorId:     5549,
			SubVendorId:  0,
			VendorName:   "VMware",
			DeviceId:     1952,
			SubDeviceId:  0,
			ParentBridge: "",
			DeviceName:   "PCI Express Root Port",
		},
		{
			Id:           "0000:03:00.0",
			ClassId:      263,
			Bus:          0x3,
			Slot:         0x0,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  5549,
			VendorName:   "VMware",
			DeviceId:     1984,
			SubDeviceId:  1984,
			ParentBridge: "0000:00:15.0",
			DeviceName:   "PVSCSI SCSI Controller",
		},
		{
			Id:           "0000:0b:00.0",
			ClassId:      512,
			Bus:          0xb,
			Slot:         0x0,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  5549,
			VendorName:   "VMware Inc.",
			DeviceId:     1968,
			SubDeviceId:  1968,
			ParentBridge: "0000:00:16.0",
			DeviceName:   "vmxnet3 Virtual Ethernet Controller",
		},
		{
			Id:           "0000:13:00.0",
			ClassId:      512,
			Bus:          0x13,
			Slot:         0x0,
			Function:     0x0,
			VendorId:     5549,
			SubVendorId:  5549,
			VendorName:   "VMware Inc.",
			DeviceId:     1968,
			SubDeviceId:  1968,
			ParentBridge: "0000:00:17.0",
			DeviceName:   "vmxnet3 Virtual Ethernet Controller",
		},
	},
	CpuFeature: []types.HostCpuIdInfo{
		{
			Level:  0,
			Vendor: "",
			Eax:    "0000:0000:0000:0000:0000:0000:0000:1101",
			Ebx:    "0111:0101:0110:1110:0110:0101:0100:0111",
			Ecx:    "0110:1100:0110:0101:0111:0100:0110:1110",
			Edx:    "0100:1001:0110:0101:0110:1110:0110:1001",
		},
		{
			Level:  1,
			Vendor: "",
			Eax:    "0000:0000:0000:0010:0000:0110:1101:0111",
			Ebx:    "0000:0000:0000:0001:0000:1000:0000:0000",
			Ecx:    "1001:0111:1011:1010:0010:0010:0010:1011",
			Edx:    "0000:1111:1010:1011:1111:1011:1111:1111",
		},
		{
			Level:  -2147483648,
			Vendor: "",
			Eax:    "1000:0000:0000:0000:0000:0000:0000:1000",
			Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
			Ecx:    "0000:0000:0000:0000:0000:0000:0000:0000",
			Edx:    "0000:0000:0000:0000:0000:0000:0000:0000",
		},
		{
			Level:  -2147483647,
			Vendor: "",
			Eax:    "0000:0000:0000:0000:0000:0000:0000:0000",
			Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
			Ecx:    "0000:0000:0000:0000:0000:0000:0000:0001",
			Edx:    "0010:1000:0001:0000:0000:1000:0000:0000",
		},
		{
			Level:  -2147483640,
			Vendor: "",
			Eax:    "0000:0000:0000:0000:0011:0000:0010:1010",
			Ebx:    "0000:0000:0000:0000:0000:0000:0000:0000",
			Ecx:    "0000:0000:0000:0000:0000:0000:0000:0000",
			Edx:    "0000:0000:0000:0000:0000:0000:0000:0000",
		},
	},
	BiosInfo: &types.HostBIOSInfo{
		BiosVersion:          "6.00",
		ReleaseDate:          nil,
		Vendor:               "",
		MajorRelease:         0,
		MinorRelease:         0,
		FirmwareMajorRelease: 0,
		FirmwareMinorRelease: 0,
	},
	ReliableMemoryInfo: &types.HostReliableMemoryInfo{},
}

func init() {
	date, _ := time.Parse("2006-01-02", "2015-07-02")

	HostHardwareInfo.BiosInfo.ReleaseDate = &date
}
