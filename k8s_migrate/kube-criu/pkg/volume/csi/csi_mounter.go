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

package csi

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/golang/glog"

	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	kstrings "k8s.io/kubernetes/pkg/util/strings"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/util"
)

//TODO (vladimirvivien) move this in a central loc later
var (
	volDataKey = struct {
		specVolID,
		volHandle,
		driverName,
		nodeName,
		attachmentID string
	}{
		"specVolID",
		"volumeHandle",
		"driverName",
		"nodeName",
		"attachmentID",
	}
)

type csiMountMgr struct {
	csiClient    csiClient
	k8s          kubernetes.Interface
	plugin       *csiPlugin
	driverName   string
	volumeID     string
	specVolumeID string
	readOnly     bool
	spec         *volume.Spec
	pod          *api.Pod
	podUID       types.UID
	options      volume.VolumeOptions
	volumeInfo   map[string]string
	volume.MetricsNil
}

// volume.Volume methods
var _ volume.Volume = &csiMountMgr{}

func (c *csiMountMgr) GetPath() string {
	dir := path.Join(getTargetPath(c.podUID, c.specVolumeID, c.plugin.host), "/mount")
	glog.V(4).Info(log("mounter.GetPath generated [%s]", dir))
	return dir
}

func getTargetPath(uid types.UID, specVolumeID string, host volume.VolumeHost) string {
	specVolID := kstrings.EscapeQualifiedNameForDisk(specVolumeID)
	return host.GetPodVolumeDir(uid, kstrings.EscapeQualifiedNameForDisk(csiPluginName), specVolID)
}

// volume.Mounter methods
var _ volume.Mounter = &csiMountMgr{}

func (c *csiMountMgr) CanMount() error {
	return nil
}

func (c *csiMountMgr) SetUp(fsGroup *int64) error {
	return c.SetUpAt(c.GetPath(), fsGroup)
}

func (c *csiMountMgr) SetUpAt(dir string, fsGroup *int64) error {
	glog.V(4).Infof(log("Mounter.SetUpAt(%s)", dir))

	mounted, err := isDirMounted(c.plugin, dir)
	if err != nil {
		glog.Error(log("mounter.SetUpAt failed while checking mount status for dir [%s]", dir))
		return err
	}

	if mounted {
		glog.V(4).Info(log("mounter.SetUpAt skipping mount, dir already mounted [%s]", dir))
		return nil
	}

	csiSource, err := getCSISourceFromSpec(c.spec)
	if err != nil {
		glog.Error(log("mounter.SetupAt failed to get CSI persistent source: %v", err))
		return err
	}

	csi := c.csiClient
	nodeName := string(c.plugin.host.GetNodeName())
	attachID := getAttachmentName(csiSource.VolumeHandle, csiSource.Driver, nodeName)

	ctx, cancel := context.WithTimeout(context.Background(), csiTimeout)
	defer cancel()

	// Check for STAGE_UNSTAGE_VOLUME set and populate deviceMountPath if so
	deviceMountPath := ""
	stageUnstageSet, err := hasStageUnstageCapability(ctx, csi)
	if err != nil {
		glog.Error(log("mounter.SetUpAt failed to check for STAGE_UNSTAGE_VOLUME capabilty: %v", err))
		return err
	}

	if stageUnstageSet {
		deviceMountPath, err = makeDeviceMountPath(c.plugin, c.spec)
		if err != nil {
			glog.Error(log("mounter.SetUpAt failed to make device mount path: %v", err))
			return err
		}
	}

	// search for attachment by VolumeAttachment.Spec.Source.PersistentVolumeName
	if c.volumeInfo == nil {
		attachment, err := c.k8s.StorageV1beta1().VolumeAttachments().Get(attachID, meta.GetOptions{})
		if err != nil {
			glog.Error(log("mounter.SetupAt failed while getting volume attachment [id=%v]: %v", attachID, err))
			return err
		}

		if attachment == nil {
			glog.Error(log("unable to find VolumeAttachment [id=%s]", attachID))
			return errors.New("no existing VolumeAttachment found")
		}
		c.volumeInfo = attachment.Status.AttachmentMetadata
	}

	attribs := csiSource.VolumeAttributes

	nodePublishSecrets := map[string]string{}
	if csiSource.NodePublishSecretRef != nil {
		nodePublishSecrets, err = getCredentialsFromSecret(c.k8s, csiSource.NodePublishSecretRef)
		if err != nil {
			return fmt.Errorf("fetching NodePublishSecretRef %s/%s failed: %v",
				csiSource.NodePublishSecretRef.Namespace, csiSource.NodePublishSecretRef.Name, err)
		}
	}

	// create target_dir before call to NodePublish
	if err := os.MkdirAll(dir, 0750); err != nil {
		glog.Error(log("mouter.SetUpAt failed to create dir %#v:  %v", dir, err))
		return err
	}
	glog.V(4).Info(log("created target path successfully [%s]", dir))

	//TODO (vladimirvivien) implement better AccessModes mapping between k8s and CSI
	accessMode := api.ReadWriteOnce
	if c.spec.PersistentVolume.Spec.AccessModes != nil {
		accessMode = c.spec.PersistentVolume.Spec.AccessModes[0]
	}

	fsType := csiSource.FSType
	err = csi.NodePublishVolume(
		ctx,
		c.volumeID,
		c.readOnly,
		deviceMountPath,
		dir,
		accessMode,
		c.volumeInfo,
		attribs,
		nodePublishSecrets,
		fsType,
	)

	if err != nil {
		glog.Errorf(log("mounter.SetupAt failed: %v", err))
		if removeMountDirErr := removeMountDir(c.plugin, dir); removeMountDirErr != nil {
			glog.Error(log("mounter.SetupAt failed to remove mount dir after a NodePublish() error [%s]: %v", dir, removeMountDirErr))
		}
		return err
	}

	// apply volume ownership
	if !c.readOnly && fsGroup != nil {
		err := volume.SetVolumeOwnership(c, fsGroup)
		if err != nil {
			// attempt to rollback mount.
			glog.Error(log("mounter.SetupAt failed to set fsgroup volume ownership for [%s]: %v", c.volumeID, err))
			glog.V(4).Info(log("mounter.SetupAt attempting to unpublish volume %s due to previous error", c.volumeID))
			if unpubErr := csi.NodeUnpublishVolume(ctx, c.volumeID, dir); unpubErr != nil {
				glog.Error(log(
					"mounter.SetupAt failed to unpublish volume [%s]: %v (caused by previous NodePublish error: %v)",
					c.volumeID, unpubErr, err,
				))
				return fmt.Errorf("%v (caused by %v)", unpubErr, err)
			}

			if unmountErr := removeMountDir(c.plugin, dir); unmountErr != nil {
				glog.Error(log(
					"mounter.SetupAt failed to clean mount dir [%s]: %v (caused by previous NodePublish error: %v)",
					dir, unmountErr, err,
				))
				return fmt.Errorf("%v (caused by %v)", unmountErr, err)
			}

			return err
		}
		glog.V(4).Info(log("mounter.SetupAt sets fsGroup to [%d] for %s", *fsGroup, c.volumeID))
	}

	glog.V(4).Infof(log("mounter.SetUp successfully requested NodePublish [%s]", dir))
	return nil
}

func (c *csiMountMgr) GetAttributes() volume.Attributes {
	mounter := c.plugin.host.GetMounter(c.plugin.GetPluginName())
	path := c.GetPath()
	supportSelinux, err := mounter.GetSELinuxSupport(path)
	if err != nil {
		glog.V(2).Info(log("error checking for SELinux support: %s", err))
		// Best guess
		supportSelinux = false
	}
	return volume.Attributes{
		ReadOnly:        c.readOnly,
		Managed:         !c.readOnly,
		SupportsSELinux: supportSelinux,
	}
}

// volume.Unmounter methods
var _ volume.Unmounter = &csiMountMgr{}

func (c *csiMountMgr) TearDown() error {
	return c.TearDownAt(c.GetPath())
}
func (c *csiMountMgr) TearDownAt(dir string) error {
	glog.V(4).Infof(log("Unmounter.TearDown(%s)", dir))

	// is dir even mounted ?
	// TODO (vladimirvivien) this check may not work for an emptyDir or local storage
	// see https://github.com/kubernetes/kubernetes/pull/56836#discussion_r155834524
	mounted, err := isDirMounted(c.plugin, dir)
	if err != nil {
		glog.Error(log("unmounter.Teardown failed while checking mount status for dir [%s]: %v", dir, err))
		return err
	}

	if !mounted {
		glog.V(4).Info(log("unmounter.Teardown skipping unmount, dir not mounted [%s]", dir))
		return nil
	}

	volID := c.volumeID
	csi := c.csiClient

	ctx, cancel := context.WithTimeout(context.Background(), csiTimeout)
	defer cancel()

	if err := csi.NodeUnpublishVolume(ctx, volID, dir); err != nil {
		glog.Errorf(log("mounter.TearDownAt failed: %v", err))
		return err
	}

	// clean mount point dir
	if err := removeMountDir(c.plugin, dir); err != nil {
		glog.Error(log("mounter.TearDownAt failed to clean mount dir [%s]: %v", dir, err))
		return err
	}
	glog.V(4).Infof(log("mounte.TearDownAt successfully unmounted dir [%s]", dir))

	return nil
}

// isDirMounted returns the !notMounted result from IsLikelyNotMountPoint check
func isDirMounted(plug *csiPlugin, dir string) (bool, error) {
	mounter := plug.host.GetMounter(plug.GetPluginName())
	notMnt, err := mounter.IsLikelyNotMountPoint(dir)
	if err != nil && !os.IsNotExist(err) {
		glog.Error(log("isDirMounted IsLikelyNotMountPoint test failed for dir [%v]", dir))
		return false, err
	}
	return !notMnt, nil
}

// removeMountDir cleans the mount dir when dir is not mounted and removed the volume data file in dir
func removeMountDir(plug *csiPlugin, mountPath string) error {
	glog.V(4).Info(log("removing mount path [%s]", mountPath))
	if pathExists, pathErr := util.PathExists(mountPath); pathErr != nil {
		glog.Error(log("failed while checking mount path stat [%s]", pathErr))
		return pathErr
	} else if !pathExists {
		glog.Warning(log("skipping mount dir removal, path does not exist [%v]", mountPath))
		return nil
	}

	mounter := plug.host.GetMounter(plug.GetPluginName())
	notMnt, err := mounter.IsLikelyNotMountPoint(mountPath)
	if err != nil {
		glog.Error(log("mount dir removal failed [%s]: %v", mountPath, err))
		return err
	}
	if notMnt {
		glog.V(4).Info(log("dir not mounted, deleting it [%s]", mountPath))
		if err := os.Remove(mountPath); err != nil && !os.IsNotExist(err) {
			glog.Error(log("failed to remove dir [%s]: %v", mountPath, err))
			return err
		}
		// remove volume data file as well
		volPath := path.Dir(mountPath)
		dataFile := path.Join(volPath, volDataFileName)
		glog.V(4).Info(log("also deleting volume info data file [%s]", dataFile))
		if err := os.Remove(dataFile); err != nil && !os.IsNotExist(err) {
			glog.Error(log("failed to delete volume data file [%s]: %v", dataFile, err))
			return err
		}
		// remove volume path
		glog.V(4).Info(log("deleting volume path [%s]", volPath))
		if err := os.Remove(volPath); err != nil && !os.IsNotExist(err) {
			glog.Error(log("failed to delete volume path [%s]: %v", volPath, err))
			return err
		}
	}
	return nil
}
