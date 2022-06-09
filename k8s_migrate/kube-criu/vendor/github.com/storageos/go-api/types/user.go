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
package types

import (
	"context"
	"encoding/json"
	"strings"
)

type User struct {
	UUID     string   `json:"id"`
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
	Password string   `json:"password,omitempty"`
	Role     string   `json:"role"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		UUID     string `json:"id"`
		Username string `json:"username"`
		Groups   string `json:"groups"`
		Password string `json:"password,omitempty"`
		Role     string `json:"role"`
	}{
		UUID:     u.UUID,
		Username: u.Username,
		Groups:   strings.Join(u.Groups, ","),
		Password: u.Password,
		Role:     u.Role,
	})

}

func (u *User) UnmarshalJSON(data []byte) error {
	temp := &struct {
		UUID     string `json:"id"`
		Username string `json:"username"`
		Groups   string `json:"groups"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{}

	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	u.UUID = temp.UUID
	u.Username = temp.Username
	u.Password = temp.Password
	u.Role = temp.Role
	u.Groups = strings.Split(temp.Groups, ",")

	return nil
}

type UserCreateOptions struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
	Password string   `json:"password"`
	Role     string   `json:"role"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}

func (u UserCreateOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Username string `json:"username"`
		Groups   string `json:"groups"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{
		Username: u.Username,
		Groups:   strings.Join(u.Groups, ","),
		Password: u.Password,
		Role:     u.Role,
	})

}
