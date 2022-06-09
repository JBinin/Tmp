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
package image

import (
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/layer"
)

// TypeLayers is used for RootFS.Type for filesystems organized into layers.
const TypeLayers = "layers"

// typeLayersWithBase is an older format used by Windows up to v1.12. We
// explicitly handle this as an error case to ensure that a daemon which still
// has an older image like this on disk can still start, even though the
// image itself is not usable. See https://github.com/docker/docker/pull/25806.
const typeLayersWithBase = "layers+base"

// RootFS describes images root filesystem
// This is currently a placeholder that only supports layers. In the future
// this can be made into an interface that supports different implementations.
type RootFS struct {
	Type    string         `json:"type"`
	DiffIDs []layer.DiffID `json:"diff_ids,omitempty"`
}

// NewRootFS returns empty RootFS struct
func NewRootFS() *RootFS {
	return &RootFS{Type: TypeLayers}
}

// Append appends a new diffID to rootfs
func (r *RootFS) Append(id layer.DiffID) {
	r.DiffIDs = append(r.DiffIDs, id)
}

// ChainID returns the ChainID for the top layer in RootFS.
func (r *RootFS) ChainID() layer.ChainID {
	if runtime.GOOS == "windows" && r.Type == typeLayersWithBase {
		logrus.Warnf("Layer type is unsupported on this platform. DiffIDs: '%v'", r.DiffIDs)
		return ""
	}
	return layer.CreateChainID(r.DiffIDs)
}
