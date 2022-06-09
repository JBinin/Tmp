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

import "sync"

type volumeDriverRegistry struct {
	nameToInitFunc     map[string]func(map[string]string) (VolumeDriver, error)
	nameToVolumeDriver map[string]VolumeDriver
	lock               *sync.RWMutex
	isShutdown         bool
}

func newVolumeDriverRegistry(nameToInitFunc map[string]func(map[string]string) (VolumeDriver, error)) *volumeDriverRegistry {
	return &volumeDriverRegistry{
		nameToInitFunc,
		make(map[string]VolumeDriver),
		&sync.RWMutex{},
		false,
	}
}

func (v *volumeDriverRegistry) Get(name string) (VolumeDriver, error) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	if v.isShutdown {
		return nil, ErrAlreadyShutdown
	}
	volumeDriver, ok := v.nameToVolumeDriver[name]
	if !ok {
		return nil, ErrDriverNotFound
	}
	return volumeDriver, nil
}

func (v *volumeDriverRegistry) Add(name string, init func(map[string]string) (VolumeDriver, error)) error {
	v.nameToInitFunc[name] = init

	return nil
}

func (v *volumeDriverRegistry) Register(name string, params map[string]string) error {
	initFunc, ok := v.nameToInitFunc[name]
	if !ok {
		return ErrNotSupported
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.isShutdown {
		return ErrAlreadyShutdown
	}
	if _, ok := v.nameToVolumeDriver[name]; ok {
		return ErrExist
	}
	volumeDriver, err := initFunc(params)
	if err != nil {
		return err
	}
	v.nameToVolumeDriver[name] = volumeDriver
	return nil
}

func (v *volumeDriverRegistry) Shutdown() error {
	v.lock.Lock()
	if v.isShutdown {
		return ErrAlreadyShutdown
	}
	for _, volumeDriver := range v.nameToVolumeDriver {
		volumeDriver.Shutdown()
	}
	v.isShutdown = true
	return nil
}
