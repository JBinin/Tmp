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
package discovery

import (
	"github.com/go-check/check"
)

func (s *DiscoverySuite) TestGeneratorNotGenerate(c *check.C) {
	ips := Generate("127.0.0.1")
	c.Assert(len(ips), check.Equals, 1)
	c.Assert(ips[0], check.Equals, "127.0.0.1")
}

func (s *DiscoverySuite) TestGeneratorWithPortNotGenerate(c *check.C) {
	ips := Generate("127.0.0.1:8080")
	c.Assert(len(ips), check.Equals, 1)
	c.Assert(ips[0], check.Equals, "127.0.0.1:8080")
}

func (s *DiscoverySuite) TestGeneratorMatchFailedNotGenerate(c *check.C) {
	ips := Generate("127.0.0.[1]")
	c.Assert(len(ips), check.Equals, 1)
	c.Assert(ips[0], check.Equals, "127.0.0.[1]")
}

func (s *DiscoverySuite) TestGeneratorWithPort(c *check.C) {
	ips := Generate("127.0.0.[1:11]:2375")
	c.Assert(len(ips), check.Equals, 11)
	c.Assert(ips[0], check.Equals, "127.0.0.1:2375")
	c.Assert(ips[1], check.Equals, "127.0.0.2:2375")
	c.Assert(ips[2], check.Equals, "127.0.0.3:2375")
	c.Assert(ips[3], check.Equals, "127.0.0.4:2375")
	c.Assert(ips[4], check.Equals, "127.0.0.5:2375")
	c.Assert(ips[5], check.Equals, "127.0.0.6:2375")
	c.Assert(ips[6], check.Equals, "127.0.0.7:2375")
	c.Assert(ips[7], check.Equals, "127.0.0.8:2375")
	c.Assert(ips[8], check.Equals, "127.0.0.9:2375")
	c.Assert(ips[9], check.Equals, "127.0.0.10:2375")
	c.Assert(ips[10], check.Equals, "127.0.0.11:2375")
}

func (s *DiscoverySuite) TestGenerateWithMalformedInputAtRangeStart(c *check.C) {
	malformedInput := "127.0.0.[x:11]:2375"
	ips := Generate(malformedInput)
	c.Assert(len(ips), check.Equals, 1)
	c.Assert(ips[0], check.Equals, malformedInput)
}

func (s *DiscoverySuite) TestGenerateWithMalformedInputAtRangeEnd(c *check.C) {
	malformedInput := "127.0.0.[1:x]:2375"
	ips := Generate(malformedInput)
	c.Assert(len(ips), check.Equals, 1)
	c.Assert(ips[0], check.Equals, malformedInput)
}
