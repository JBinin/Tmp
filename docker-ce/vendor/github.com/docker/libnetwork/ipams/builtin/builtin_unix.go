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
// +build linux freebsd solaris darwin

package builtin

import (
	"errors"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/ipam"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/ipamutils"
)

// Init registers the built-in ipam service with libnetwork
func Init(ic ipamapi.Callback, l, g interface{}) error {
	var (
		ok                bool
		localDs, globalDs datastore.DataStore
	)

	if l != nil {
		if localDs, ok = l.(datastore.DataStore); !ok {
			return errors.New("incorrect local datastore passed to built-in ipam init")
		}
	}

	if g != nil {
		if globalDs, ok = g.(datastore.DataStore); !ok {
			return errors.New("incorrect global datastore passed to built-in ipam init")
		}
	}

	ipamutils.InitNetworks()

	a, err := ipam.NewAllocator(localDs, globalDs)
	if err != nil {
		return err
	}

	cps := &ipamapi.Capability{RequiresRequestReplay: true}

	return ic.RegisterIpamDriverWithCapabilities(ipamapi.DefaultIPAM, a, cps)
}
