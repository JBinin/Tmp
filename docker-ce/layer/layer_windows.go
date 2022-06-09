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
package layer

import "errors"

// GetLayerPath returns the path to a layer
func GetLayerPath(s Store, layer ChainID) (string, error) {
	ls, ok := s.(*layerStore)
	if !ok {
		return "", errors.New("unsupported layer store")
	}
	ls.layerL.Lock()
	defer ls.layerL.Unlock()

	rl, ok := ls.layerMap[layer]
	if !ok {
		return "", ErrLayerDoesNotExist
	}

	path, err := ls.driver.Get(rl.cacheID, "")
	if err != nil {
		return "", err
	}

	if err := ls.driver.Put(rl.cacheID); err != nil {
		return "", err
	}

	return path, nil
}

func (ls *layerStore) mountID(name string) string {
	// windows has issues if container ID doesn't match mount ID
	return name
}
