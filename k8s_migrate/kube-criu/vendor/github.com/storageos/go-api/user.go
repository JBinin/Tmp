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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/storageos/go-api/types"
)

var (

	// UserAPIPrefix is a partial path to the HTTP endpoint.
	UserAPIPrefix = "users"

	// ErrNoSuchUser is the error returned when the user does not exist.
	ErrNoSuchUser = errors.New("no such user")
)

// UserList returns the list of available users.
func (c *Client) UserList(opts types.ListOptions) ([]*types.User, error) {
	listOpts := doOptions{
		fieldSelector: opts.FieldSelector,
		labelSelector: opts.LabelSelector,
		namespace:     opts.Namespace,
		context:       opts.Context,
	}

	if opts.LabelSelector != "" {
		query := url.Values{}
		query.Add("labelSelector", opts.LabelSelector)
		listOpts.values = query
	}

	resp, err := c.do("GET", UserAPIPrefix, listOpts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	users := make([]*types.User, 0)
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// User returns a user by its username/id.
func (c *Client) User(username string) (*types.User, error) {
	path := fmt.Sprintf("%s/%s", UserAPIPrefix, username)
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrNoSuchUser
		}
		return nil, err
	}
	defer resp.Body.Close()

	var user *types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

// UserCreate creates a user on the server.
func (c *Client) UserCreate(opts types.UserCreateOptions) error {
	_, err := c.do("POST", UserAPIPrefix, doOptions{
		data:    opts,
		context: opts.Context,
	})
	return err
}

// UserUpdate updates a user on the server.
func (c *Client) UserUpdate(user *types.User, ctx context.Context) error {
	var ref string
	switch {
	case user.UUID != "":
		ref = user.UUID
	case user.Username != "":
		ref = user.Username
	default:
		return ErrNoSuchUser
	}

	path := fmt.Sprintf("%s/%s", UserAPIPrefix, ref)
	resp, err := c.do("POST", path, doOptions{
		data:    user,
		context: ctx,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return ErrNoSuchUser
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}

// UserDelete removes a user by its reference.
func (c *Client) UserDelete(opts types.DeleteOptions) error {
	resp, err := c.do("DELETE", UserAPIPrefix+"/"+opts.Name, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok {
			if e.Status == http.StatusNotFound {
				return ErrNoSuchUser
			}
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}
