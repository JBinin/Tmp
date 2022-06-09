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

// Status can be used to query the Status endpoints
type Status struct {
	c *Client
}

// Status returns a handle to the status endpoints
func (c *Client) Status() *Status {
	return &Status{c}
}

// Leader is used to query for a known leader
func (s *Status) Leader() (string, error) {
	r := s.c.newRequest("GET", "/v1/status/leader")
	_, resp, err := requireOK(s.c.doRequest(r))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var leader string
	if err := decodeBody(resp, &leader); err != nil {
		return "", err
	}
	return leader, nil
}

// Peers is used to query for a known raft peers
func (s *Status) Peers() ([]string, error) {
	r := s.c.newRequest("GET", "/v1/status/peers")
	_, resp, err := requireOK(s.c.doRequest(r))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var peers []string
	if err := decodeBody(resp, &peers); err != nil {
		return nil, err
	}
	return peers, nil
}
