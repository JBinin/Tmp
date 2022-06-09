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
package builders

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

// Swarm creates a swarm with default values.
// Any number of swarm function builder can be pass to augment it.
func Swarm(swarmBuilders ...func(*swarm.Swarm)) *swarm.Swarm {
	t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	swarm := &swarm.Swarm{
		ClusterInfo: swarm.ClusterInfo{
			ID: "swarm",
			Meta: swarm.Meta{
				CreatedAt: t1,
			},
			Spec: swarm.Spec{},
		},
		JoinTokens: swarm.JoinTokens{
			Worker:  "worker-join-token",
			Manager: "manager-join-token",
		},
	}

	for _, builder := range swarmBuilders {
		builder(swarm)
	}

	return swarm
}

// Autolock set the swarm into autolock mode
func Autolock() func(*swarm.Swarm) {
	return func(swarm *swarm.Swarm) {
		swarm.Spec.EncryptionConfig.AutoLockManagers = true
	}
}
