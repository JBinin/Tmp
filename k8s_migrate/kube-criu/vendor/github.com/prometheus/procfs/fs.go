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
package procfs

import (
	"fmt"
	"os"
	"path"

	"github.com/prometheus/procfs/xfs"
)

// FS represents the pseudo-filesystem proc, which provides an interface to
// kernel data structures.
type FS string

// DefaultMountPoint is the common mount point of the proc filesystem.
const DefaultMountPoint = "/proc"

// NewFS returns a new FS mounted under the given mountPoint. It will error
// if the mount point can't be read.
func NewFS(mountPoint string) (FS, error) {
	info, err := os.Stat(mountPoint)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %s", mountPoint, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("mount point %s is not a directory", mountPoint)
	}

	return FS(mountPoint), nil
}

// Path returns the path of the given subsystem relative to the procfs root.
func (fs FS) Path(p ...string) string {
	return path.Join(append([]string{string(fs)}, p...)...)
}

// XFSStats retrieves XFS filesystem runtime statistics.
func (fs FS) XFSStats() (*xfs.Stats, error) {
	f, err := os.Open(fs.Path("fs/xfs/stat"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return xfs.ParseStats(f)
}
