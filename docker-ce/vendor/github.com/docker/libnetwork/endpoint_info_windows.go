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
// +build windows

package libnetwork

import "fmt"

func (ep *endpoint) DriverInfo() (map[string]interface{}, error) {
	ep, err := ep.retrieveFromStore()
	if err != nil {
		return nil, err
	}

	var gwDriverInfo map[string]interface{}
	if sb, ok := ep.getSandbox(); ok {
		if gwep := sb.getEndpointInGWNetwork(); gwep != nil && gwep.ID() != ep.ID() {

			gwDriverInfo, err = gwep.DriverInfo()
			if err != nil {
				return nil, err
			}
		}
	}

	n, err := ep.getNetworkFromStore()
	if err != nil {
		return nil, fmt.Errorf("could not find network in store for driver info: %v", err)
	}

	driver, err := n.driver(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver info: %v", err)
	}

	epInfo, err := driver.EndpointOperInfo(n.ID(), ep.ID())
	if err != nil {
		return nil, err
	}

	if epInfo != nil {
		epInfo["GW_INFO"] = gwDriverInfo
		return epInfo, nil
	}

	return gwDriverInfo, nil
}
