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
/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package authenticatorfactory

type RequestHeaderConfig struct {
	// UsernameHeaders are the headers to check (in order, case-insensitively) for an identity. The first header with a value wins.
	UsernameHeaders []string
	// GroupHeaders are the headers to check (case-insensitively) for a group names.  All values will be used.
	GroupHeaders []string
	// ExtraHeaderPrefixes are the head prefixes to check (case-insentively) for filling in
	// the user.Info.Extra.  All values of all matching headers will be added.
	ExtraHeaderPrefixes []string
	// ClientCA points to CA bundle file which is used verify the identity of the front proxy
	ClientCA string
	// AllowedClientNames is a list of common names that may be presented by the authenticating front proxy.  Empty means: accept any.
	AllowedClientNames []string
}
