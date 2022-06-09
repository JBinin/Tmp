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
Package volumeattach provides the ability to attach and detach volumes
from servers.

Example to Attach a Volume

	serverID := "7ac8686c-de71-4acb-9600-ec18b1a1ed6d"
	volumeID := "87463836-f0e2-4029-abf6-20c8892a3103"

	createOpts := volumeattach.CreateOpts{
		Device:   "/dev/vdc",
		VolumeID: volumeID,
	}

	result, err := volumeattach.Create(computeClient, serverID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Detach a Volume

	serverID := "7ac8686c-de71-4acb-9600-ec18b1a1ed6d"
	attachmentID := "ed081613-1c9b-4231-aa5e-ebfd4d87f983"

	err := volumeattach.Delete(computeClient, serverID, attachmentID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package volumeattach
