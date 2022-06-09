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
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/image/cache"
)

// MakeImageCache creates a stateful image cache.
func (daemon *Daemon) MakeImageCache(sourceRefs []string) builder.ImageCache {
	if len(sourceRefs) == 0 {
		return cache.NewLocal(daemon.imageStore)
	}

	cache := cache.New(daemon.imageStore)

	for _, ref := range sourceRefs {
		img, err := daemon.GetImage(ref)
		if err != nil {
			logrus.Warnf("Could not look up %s for cache resolution, skipping: %+v", ref, err)
			continue
		}
		cache.Populate(img)
	}

	return cache
}
