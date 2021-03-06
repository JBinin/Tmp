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
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/registrar"
	"github.com/docker/docker/pkg/stringid"
)

var (
	validContainerNameChars   = api.RestrictedNameChars
	validContainerNamePattern = api.RestrictedNamePattern
)

func (daemon *Daemon) registerName(container *container.Container) error {
	if daemon.Exists(container.ID) {
		return fmt.Errorf("Container is already loaded")
	}
	if err := validateID(container.ID); err != nil {
		return err
	}
	if container.Name == "" {
		name, err := daemon.generateNewName(container.ID)
		if err != nil {
			return err
		}
		container.Name = name

		if err := container.ToDiskLocking(); err != nil {
			logrus.Errorf("Error saving container name to disk: %v", err)
		}
	}
	return daemon.nameIndex.Reserve(container.Name, container.ID)
}

func (daemon *Daemon) generateIDAndName(name string) (string, string, error) {
	var (
		err error
		id  = stringid.GenerateNonCryptoID()
	)

	if name == "" {
		if name, err = daemon.generateNewName(id); err != nil {
			return "", "", err
		}
		return id, name, nil
	}

	if name, err = daemon.reserveName(id, name); err != nil {
		return "", "", err
	}

	return id, name, nil
}

func (daemon *Daemon) reserveName(id, name string) (string, error) {
	if !validContainerNamePattern.MatchString(strings.TrimPrefix(name, "/")) {
		return "", fmt.Errorf("Invalid container name (%s), only %s are allowed", name, validContainerNameChars)
	}
	if name[0] != '/' {
		name = "/" + name
	}

	if err := daemon.nameIndex.Reserve(name, id); err != nil {
		if err == registrar.ErrNameReserved {
			id, err := daemon.nameIndex.Get(name)
			if err != nil {
				logrus.Errorf("got unexpected error while looking up reserved name: %v", err)
				return "", err
			}
			return "", fmt.Errorf("Conflict. The container name %q is already in use by container %q. You have to remove (or rename) that container to be able to reuse that name.", name, id)
		}
		return "", fmt.Errorf("error reserving name: %q, error: %v", name, err)
	}
	return name, nil
}

func (daemon *Daemon) releaseName(name string) {
	daemon.nameIndex.Release(name)
}

func (daemon *Daemon) generateNewName(id string) (string, error) {
	var name string
	for i := 0; i < 6; i++ {
		name = namesgenerator.GetRandomName(i)
		if name[0] != '/' {
			name = "/" + name
		}

		if err := daemon.nameIndex.Reserve(name, id); err != nil {
			if err == registrar.ErrNameReserved {
				continue
			}
			return "", err
		}
		return name, nil
	}

	name = "/" + stringid.TruncateID(id)
	if err := daemon.nameIndex.Reserve(name, id); err != nil {
		return "", err
	}
	return name, nil
}

func validateID(id string) error {
	if id == "" {
		return fmt.Errorf("Invalid empty id")
	}
	return nil
}
