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
package daemon

import (
	"fmt"
	"io"
	"runtime"

	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/ioutils"
)

// ContainerExport writes the contents of the container to the given
// writer. An error is returned if the container cannot be found.
func (daemon *Daemon) ContainerExport(name string, out io.Writer) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("the daemon on this platform does not support export of a container")
	}

	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	data, err := daemon.containerExport(container)
	if err != nil {
		return fmt.Errorf("Error exporting container %s: %v", name, err)
	}
	defer data.Close()

	// Stream the entire contents of the container (basically a volatile snapshot)
	if _, err := io.Copy(out, data); err != nil {
		return fmt.Errorf("Error exporting container %s: %v", name, err)
	}
	return nil
}

func (daemon *Daemon) containerExport(container *container.Container) (io.ReadCloser, error) {
	if err := daemon.Mount(container); err != nil {
		return nil, err
	}

	uidMaps, gidMaps := daemon.GetUIDGIDMaps()
	archive, err := archive.TarWithOptions(container.BaseFS, &archive.TarOptions{
		Compression: archive.Uncompressed,
		UIDMaps:     uidMaps,
		GIDMaps:     gidMaps,
	})
	if err != nil {
		daemon.Unmount(container)
		return nil, err
	}
	arch := ioutils.NewReadCloserWrapper(archive, func() error {
		err := archive.Close()
		daemon.Unmount(container)
		return err
	})
	daemon.LogContainerEvent(container, "export")
	return arch, err
}
