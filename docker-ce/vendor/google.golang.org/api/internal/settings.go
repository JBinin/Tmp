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
// Package internal supports the options and transport packages.
package internal

import (
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

// DialSettings holds information needed to establish a connection with a
// Google API service.
type DialSettings struct {
	Endpoint                   string
	Scopes                     []string
	ServiceAccountJSONFilename string // if set, TokenSource is ignored.
	TokenSource                oauth2.TokenSource
	UserAgent                  string
	APIKey                     string
	HTTPClient                 *http.Client
	GRPCDialOpts               []grpc.DialOption
	GRPCConn                   *grpc.ClientConn
}
