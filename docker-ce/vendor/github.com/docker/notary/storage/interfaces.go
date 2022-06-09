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
package storage

// NoSizeLimit is represented as -1 for arguments to GetMeta
const NoSizeLimit int64 = -1

// MetadataStore must be implemented by anything that intends to interact
// with a store of TUF files
type MetadataStore interface {
	GetSized(name string, size int64) ([]byte, error)
	Set(name string, blob []byte) error
	SetMulti(map[string][]byte) error
	RemoveAll() error
	Remove(name string) error
}

// PublicKeyStore must be implemented by a key service
type PublicKeyStore interface {
	GetKey(role string) ([]byte, error)
	RotateKey(role string) ([]byte, error)
}

// RemoteStore is similar to LocalStore with the added expectation that it should
// provide a way to download targets once located
type RemoteStore interface {
	MetadataStore
	PublicKeyStore
}

// Bootstrapper is a thing that can set itself up
type Bootstrapper interface {
	// Bootstrap instructs a configured Bootstrapper to perform
	// its setup operations.
	Bootstrap() error
}
