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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	api "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	fakeclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/volume"
	volumetest "k8s.io/kubernetes/pkg/volume/testing"
	"k8s.io/kubernetes/pkg/volume/util"
)

var (
	testDriver = "test-driver"
	testVol    = "vol-123"
	testns     = "test-ns"
	testPodUID = types.UID("test-pod")
)

func TestMounterGetPath(t *testing.T) {
	plug, tmpDir := newTestPlugin(t)
	defer os.RemoveAll(tmpDir)

	// TODO (vladimirvivien) specName with slashes will not work
	testCases := []struct {
		name           string
		specVolumeName string
		path           string
	}{
		{
			name:           "simple specName",
			specVolumeName: "spec-0",
			path:           path.Join(tmpDir, fmt.Sprintf("pods/%s/volumes/kubernetes.io~csi/%s/%s", testPodUID, "spec-0", "/mount")),
		},
		{
			name:           "specName with dots",
			specVolumeName: "test.spec.1",
			path:           path.Join(tmpDir, fmt.Sprintf("pods/%s/volumes/kubernetes.io~csi/%s/%s", testPodUID, "test.spec.1", "/mount")),
		},
	}
	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		pv := makeTestPV(tc.specVolumeName, 10, testDriver, testVol)
		spec := volume.NewSpecFromPersistentVolume(pv, pv.Spec.PersistentVolumeSource.CSI.ReadOnly)
		mounter, err := plug.NewMounter(
			spec,
			&api.Pod{ObjectMeta: meta.ObjectMeta{UID: testPodUID, Namespace: testns}},
			volume.VolumeOptions{},
		)
		if err != nil {
			t.Fatalf("Failed to make a new Mounter: %v", err)
		}
		csiMounter := mounter.(*csiMountMgr)

		path := csiMounter.GetPath()

		if tc.path != path {
			t.Errorf("expecting path %s, got %s", tc.path, path)
		}
	}
}

func TestMounterSetUp(t *testing.T) {
	plug, tmpDir := newTestPlugin(t)
	defer os.RemoveAll(tmpDir)
	fakeClient := fakeclient.NewSimpleClientset()
	host := volumetest.NewFakeVolumeHostWithNodeName(
		tmpDir,
		fakeClient,
		nil,
		"fakeNode",
	)
	plug.host = host
	pv := makeTestPV("test-pv", 10, testDriver, testVol)
	pvName := pv.GetName()

	mounter, err := plug.NewMounter(
		volume.NewSpecFromPersistentVolume(pv, pv.Spec.PersistentVolumeSource.CSI.ReadOnly),
		&api.Pod{ObjectMeta: meta.ObjectMeta{UID: testPodUID, Namespace: testns}},
		volume.VolumeOptions{},
	)
	if err != nil {
		t.Fatalf("Failed to make a new Mounter: %v", err)
	}

	if mounter == nil {
		t.Fatal("failed to create CSI mounter")
	}

	csiMounter := mounter.(*csiMountMgr)
	csiMounter.csiClient = setupClient(t, true)

	attachID := getAttachmentName(csiMounter.volumeID, csiMounter.driverName, string(plug.host.GetNodeName()))

	attachment := &storage.VolumeAttachment{
		ObjectMeta: meta.ObjectMeta{
			Name: attachID,
		},
		Spec: storage.VolumeAttachmentSpec{
			NodeName: "test-node",
			Attacher: csiPluginName,
			Source: storage.VolumeAttachmentSource{
				PersistentVolumeName: &pvName,
			},
		},
		Status: storage.VolumeAttachmentStatus{
			Attached:    false,
			AttachError: nil,
			DetachError: nil,
		},
	}
	_, err = csiMounter.k8s.StorageV1beta1().VolumeAttachments().Create(attachment)
	if err != nil {
		t.Fatalf("failed to setup VolumeAttachment: %v", err)
	}

	// Mounter.SetUp()
	fsGroup := int64(2000)
	if err := csiMounter.SetUp(&fsGroup); err != nil {
		t.Fatalf("mounter.Setup failed: %v", err)
	}

	//Test the default value of file system type is not overridden
	if len(csiMounter.spec.PersistentVolume.Spec.CSI.FSType) != 0 {
		t.Errorf("default value of file system type was overridden by type %s", csiMounter.spec.PersistentVolume.Spec.CSI.FSType)
	}

	path := csiMounter.GetPath()
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("SetUp() failed, volume path not created: %s", path)
		} else {
			t.Errorf("SetUp() failed: %v", err)
		}
	}

	// ensure call went all the way
	pubs := csiMounter.csiClient.(*fakeCsiDriverClient).nodeClient.GetNodePublishedVolumes()
	if pubs[csiMounter.volumeID] != csiMounter.GetPath() {
		t.Error("csi server may not have received NodePublishVolume call")
	}
}

func TestUnmounterTeardown(t *testing.T) {
	plug, tmpDir := newTestPlugin(t)
	defer os.RemoveAll(tmpDir)
	pv := makeTestPV("test-pv", 10, testDriver, testVol)

	// save the data file prior to unmount
	dir := path.Join(getTargetPath(testPodUID, pv.ObjectMeta.Name, plug.host), "/mount")
	if err := os.MkdirAll(dir, 0755); err != nil && !os.IsNotExist(err) {
		t.Errorf("failed to create dir [%s]: %v", dir, err)
	}

	// do a fake local mount
	diskMounter := util.NewSafeFormatAndMountFromHost(plug.GetPluginName(), plug.host)
	if err := diskMounter.FormatAndMount("/fake/device", dir, "testfs", nil); err != nil {
		t.Errorf("failed to mount dir [%s]: %v", dir, err)
	}

	if err := saveVolumeData(
		path.Dir(dir),
		volDataFileName,
		map[string]string{
			volDataKey.specVolID:  pv.ObjectMeta.Name,
			volDataKey.driverName: testDriver,
			volDataKey.volHandle:  testVol,
		},
	); err != nil {
		t.Fatalf("failed to save volume data: %v", err)
	}

	unmounter, err := plug.NewUnmounter(pv.ObjectMeta.Name, testPodUID)
	if err != nil {
		t.Fatalf("failed to make a new Unmounter: %v", err)
	}

	csiUnmounter := unmounter.(*csiMountMgr)
	csiUnmounter.csiClient = setupClient(t, true)
	err = csiUnmounter.TearDownAt(dir)
	if err != nil {
		t.Fatal(err)
	}

	// ensure csi client call
	pubs := csiUnmounter.csiClient.(*fakeCsiDriverClient).nodeClient.GetNodePublishedVolumes()
	if _, ok := pubs[csiUnmounter.volumeID]; ok {
		t.Error("csi server may not have received NodeUnpublishVolume call")
	}

}

func TestSaveVolumeData(t *testing.T) {
	plug, tmpDir := newTestPlugin(t)
	defer os.RemoveAll(tmpDir)
	testCases := []struct {
		name       string
		data       map[string]string
		shouldFail bool
	}{
		{name: "test with data ok", data: map[string]string{"key0": "val0", "_key1": "val1", "key2": "val2"}},
		{name: "test with data ok 2 ", data: map[string]string{"_key0_": "val0", "&key1": "val1", "key2": "val2"}},
	}

	for i, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		specVolID := fmt.Sprintf("spec-volid-%d", i)
		mountDir := path.Join(getTargetPath(testPodUID, specVolID, plug.host), "/mount")
		if err := os.MkdirAll(mountDir, 0755); err != nil && !os.IsNotExist(err) {
			t.Errorf("failed to create dir [%s]: %v", mountDir, err)
		}

		err := saveVolumeData(path.Dir(mountDir), volDataFileName, tc.data)

		if !tc.shouldFail && err != nil {
			t.Errorf("unexpected failure: %v", err)
		}
		// did file get created
		dataDir := getTargetPath(testPodUID, specVolID, plug.host)
		file := path.Join(dataDir, volDataFileName)
		if _, err := os.Stat(file); err != nil {
			t.Errorf("failed to create data dir: %v", err)
		}

		// validate content
		data, err := ioutil.ReadFile(file)
		if !tc.shouldFail && err != nil {
			t.Errorf("failed to read data file: %v", err)
		}

		jsonData := new(bytes.Buffer)
		if err := json.NewEncoder(jsonData).Encode(tc.data); err != nil {
			t.Errorf("failed to encode json: %v", err)
		}
		if string(data) != jsonData.String() {
			t.Errorf("expecting encoded data %v, got %v", string(data), jsonData)
		}
	}
}
