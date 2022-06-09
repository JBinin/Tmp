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
package overlay

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/libnetwork/osl"
)

func validateID(nid, eid string) error {
	if nid == "" {
		return fmt.Errorf("invalid network id")
	}

	if eid == "" {
		return fmt.Errorf("invalid endpoint id")
	}

	return nil
}

func createVxlan(name string, vni uint32, mtu int) error {
	defer osl.InitOSContext()()

	// Get default interface to plumb the vxlan on
	routeCmd := "/usr/sbin/ipadm show-addr -p -o addrobj " +
		"`/usr/sbin/route get default | /usr/bin/grep interface | " +
		"/usr/bin/awk '{print $2}'`"
	out, err := exec.Command("/usr/bin/bash", "-c", routeCmd).Output()
	if err != nil {
		return fmt.Errorf("cannot get default route: %v", err)
	}

	defaultInterface := strings.SplitN(string(out), "/", 2)
	propList := fmt.Sprintf("interface=%s,vni=%d", defaultInterface[0], vni)

	out, err = exec.Command("/usr/sbin/dladm", "create-vxlan", "-t", "-p", propList,
		name).Output()
	if err != nil {
		return fmt.Errorf("error creating vxlan interface: %v %s", err, out)
	}

	return nil
}

func deleteInterfaceBySubnet(brPrefix string, s *subnet) error {
	return nil

}

func deleteInterface(name string) error {
	defer osl.InitOSContext()()

	out, err := exec.Command("/usr/sbin/dladm", "delete-vxlan", name).Output()
	if err != nil {
		return fmt.Errorf("error creating vxlan interface: %v %s", err, out)
	}

	return nil
}
