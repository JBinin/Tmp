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
package libnetwork

import (
	"fmt"
	"strconv"

	"github.com/docker/libnetwork/drivers/bridge"
)

const libnGWNetwork = "docker_gwbridge"

func getPlatformOption() EndpointOption {
	return nil
}

func (c *controller) createGWNetwork() (Network, error) {
	netOption := map[string]string{
		bridge.BridgeName:         libnGWNetwork,
		bridge.EnableICC:          strconv.FormatBool(false),
		bridge.EnableIPMasquerade: strconv.FormatBool(true),
	}

	n, err := c.NewNetwork("bridge", libnGWNetwork, "",
		NetworkOptionDriverOpts(netOption),
		NetworkOptionEnableIPv6(false),
	)

	if err != nil {
		return nil, fmt.Errorf("error creating external connectivity network: %v", err)
	}
	return n, err
}
