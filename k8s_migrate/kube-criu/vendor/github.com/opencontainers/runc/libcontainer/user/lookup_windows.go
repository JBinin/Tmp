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

package user

import (
	"fmt"
	"os/user"
)

func lookupUser(username string) (User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return User{}, err
	}
	return userFromOS(u)
}

func lookupUid(uid int) (User, error) {
	u, err := user.LookupId(fmt.Sprintf("%d", uid))
	if err != nil {
		return User{}, err
	}
	return userFromOS(u)
}

func lookupGroup(groupname string) (Group, error) {
	g, err := user.LookupGroup(groupname)
	if err != nil {
		return Group{}, err
	}
	return groupFromOS(g)
}

func lookupGid(gid int) (Group, error) {
	g, err := user.LookupGroupId(fmt.Sprintf("%d", gid))
	if err != nil {
		return Group{}, err
	}
	return groupFromOS(g)
}
