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
// Package cache provides facilities to speed up access to the storage
// backend.
package cache

import (
	"fmt"

	"github.com/docker/distribution"
)

// BlobDescriptorCacheProvider provides repository scoped
// BlobDescriptorService cache instances and a global descriptor cache.
type BlobDescriptorCacheProvider interface {
	distribution.BlobDescriptorService

	RepositoryScoped(repo string) (distribution.BlobDescriptorService, error)
}

// ValidateDescriptor provides a helper function to ensure that caches have
// common criteria for admitting descriptors.
func ValidateDescriptor(desc distribution.Descriptor) error {
	if err := desc.Digest.Validate(); err != nil {
		return err
	}

	if desc.Size < 0 {
		return fmt.Errorf("cache: invalid length in descriptor: %v < 0", desc.Size)
	}

	if desc.MediaType == "" {
		return fmt.Errorf("cache: empty mediatype on descriptor: %v", desc)
	}

	return nil
}
