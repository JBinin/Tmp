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
package storageos

import (
	"encoding/json"
	"errors"
)

var (
	// LoginAPIPrefix is a partial path to the HTTP endpoint.
	LoginAPIPrefix = "auth/login"
	ErrLoginFailed = errors.New("Failed to get token from API endpoint")
)

// Login attemps to get a token from the API
func (c *Client) Login() (token string, err error) {
	resp, err := c.do("POST", LoginAPIPrefix, doOptions{data: struct {
		User string `json:"username"`
		Pass string `json:"password"`
	}{c.username, c.secret}})

	if err != nil {
		if _, ok := err.(*Error); ok {
			return "", ErrLoginFailed
		}

		return "", err
	}

	if resp.StatusCode != 200 {
		return "", ErrLoginFailed
	}

	unmarsh := struct {
		Token string `json:"token"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&unmarsh); err != nil {
		return "", err
	}

	if unmarsh.Token == "" {
		return "", ErrLoginFailed
	}

	return unmarsh.Token, nil
}
