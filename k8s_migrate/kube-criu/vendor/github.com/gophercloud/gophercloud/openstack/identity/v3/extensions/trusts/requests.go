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
package trusts

import "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"

// AuthOptsExt extends the base Identity v3 tokens AuthOpts with a TrustID.
type AuthOptsExt struct {
	tokens.AuthOptionsBuilder

	// TrustID is the ID of the trust.
	TrustID string `json:"id"`
}

// ToTokenV3CreateMap builds a create request body from the AuthOpts.
func (opts AuthOptsExt) ToTokenV3CreateMap(scope map[string]interface{}) (map[string]interface{}, error) {
	return opts.AuthOptionsBuilder.ToTokenV3CreateMap(scope)
}

// ToTokenV3ScopeMap builds a scope from AuthOpts.
func (opts AuthOptsExt) ToTokenV3ScopeMap() (map[string]interface{}, error) {
	b, err := opts.AuthOptionsBuilder.ToTokenV3ScopeMap()
	if err != nil {
		return nil, err
	}

	if opts.TrustID != "" {
		if b == nil {
			b = make(map[string]interface{})
		}
		b["OS-TRUST:trust"] = map[string]interface{}{
			"id": opts.TrustID,
		}
	}

	return b, nil
}

func (opts AuthOptsExt) CanReauth() bool {
	return opts.AuthOptionsBuilder.CanReauth()
}
