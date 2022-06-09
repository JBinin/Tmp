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
package api

// Raw can be used to do raw queries against custom endpoints
type Raw struct {
	c *Client
}

// Raw returns a handle to query endpoints
func (c *Client) Raw() *Raw {
	return &Raw{c}
}

// Query is used to do a GET request against an endpoint
// and deserialize the response into an interface using
// standard Consul conventions.
func (raw *Raw) Query(endpoint string, out interface{}, q *QueryOptions) (*QueryMeta, error) {
	return raw.c.query(endpoint, out, q)
}

// Write is used to do a PUT request against an endpoint
// and serialize/deserialized using the standard Consul conventions.
func (raw *Raw) Write(endpoint string, in, out interface{}, q *WriteOptions) (*WriteMeta, error) {
	return raw.c.write(endpoint, in, out, q)
}
