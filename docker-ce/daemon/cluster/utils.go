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
package cluster

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/ioutils"
)

func loadPersistentState(root string) (*nodeStartConfig, error) {
	dt, err := ioutil.ReadFile(filepath.Join(root, stateFile))
	if err != nil {
		return nil, err
	}
	// missing certificate means no actual state to restore from
	if _, err := os.Stat(filepath.Join(root, "certificates/swarm-node.crt")); err != nil {
		if os.IsNotExist(err) {
			clearPersistentState(root)
		}
		return nil, err
	}
	var st nodeStartConfig
	if err := json.Unmarshal(dt, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

func savePersistentState(root string, config nodeStartConfig) error {
	dt, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutils.AtomicWriteFile(filepath.Join(root, stateFile), dt, 0600)
}

func clearPersistentState(root string) error {
	// todo: backup this data instead of removing?
	if err := os.RemoveAll(root); err != nil {
		return err
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	return nil
}

func removingManagerCausesLossOfQuorum(reachable, unreachable int) bool {
	return reachable-2 <= unreachable
}

func isLastManager(reachable, unreachable int) bool {
	return reachable == 1 && unreachable == 0
}
