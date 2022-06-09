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
package volume

import (
	// TODO return types need to be refactored into pkg
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Backend is the methods that need to be implemented to provide
// volume specific functionality
type Backend interface {
	Volumes(filter string) ([]*types.Volume, []string, error)
	VolumeInspect(name string) (*types.Volume, error)
	VolumeCreate(name, driverName string, opts, labels map[string]string) (*types.Volume, error)
	VolumeRm(name string, force bool) error
	VolumesPrune(pruneFilters filters.Args) (*types.VolumesPruneReport, error)
}
