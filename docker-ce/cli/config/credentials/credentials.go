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
package credentials

import (
	"github.com/docker/docker/api/types"
)

// Store is the interface that any credentials store must implement.
type Store interface {
	// Erase removes credentials from the store for a given server.
	Erase(serverAddress string) error
	// Get retrieves credentials from the store for a given server.
	Get(serverAddress string) (types.AuthConfig, error)
	// GetAll retrieves all the credentials from the store.
	GetAll() (map[string]types.AuthConfig, error)
	// Store saves credentials in the store.
	Store(authConfig types.AuthConfig) error
}
