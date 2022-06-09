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

// ContainerCreateWorkdir creates the working directory. This solves the
// issue arising from https://github.com/docker/docker/issues/27545,
// which was initially fixed by https://github.com/docker/docker/pull/27884. But that fix
// was too expensive in terms of performance on Windows. Instead,
// https://github.com/docker/docker/pull/28514 introduces this new functionality
// where the builder calls into the backend here to create the working directory.
func (daemon *Daemon) ContainerCreateWorkdir(cID string) error {
	container, err := daemon.GetContainer(cID)
	if err != nil {
		return err
	}
	err = daemon.Mount(container)
	if err != nil {
		return err
	}
	defer daemon.Unmount(container)
	rootUID, rootGID := daemon.GetRemappedUIDGID()
	return container.SetupWorkingDirectory(rootUID, rootGID)
}
