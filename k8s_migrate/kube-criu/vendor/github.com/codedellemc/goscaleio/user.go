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
package goscaleio

import (
	"fmt"

	types "github.com/codedellemc/goscaleio/types/v1"
)

func (system *System) GetUser() (user []types.User, err error) {
	endpoint := system.client.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/System::%v/relationships/User", system.System.ID)

	req := system.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", system.client.Token)
	req.Header.Add("Accept", "application/json;version="+system.client.configConnect.Version)

	resp, err := system.client.retryCheckResp(&system.client.Http, req)
	if err != nil {
		return []types.User{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = system.client.decodeBody(resp, &user); err != nil {
		return []types.User{}, fmt.Errorf("error decoding instances response: %s", err)
	}

	// bs, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return types.User{}, errors.New("error reading body")
	// }
	//
	// fmt.Println(string(bs))
	// return types.User{}, nil
	return user, nil
}
