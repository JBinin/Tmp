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
package dns

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ClientConfig wraps the contents of the /etc/resolv.conf file.
type ClientConfig struct {
	Servers  []string // servers to use
	Search   []string // suffixes to append to local name
	Port     string   // what port to use
	Ndots    int      // number of dots in name to trigger absolute lookup
	Timeout  int      // seconds before giving up on packet
	Attempts int      // lost packets before giving up on server, not used in the package dns
}

// ClientConfigFromFile parses a resolv.conf(5) like file and returns
// a *ClientConfig.
func ClientConfigFromFile(resolvconf string) (*ClientConfig, error) {
	file, err := os.Open(resolvconf)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	c := new(ClientConfig)
	scanner := bufio.NewScanner(file)
	c.Servers = make([]string, 0)
	c.Search = make([]string, 0)
	c.Port = "53"
	c.Ndots = 1
	c.Timeout = 5
	c.Attempts = 2

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()
		f := strings.Fields(line)
		if len(f) < 1 {
			continue
		}
		switch f[0] {
		case "nameserver": // add one name server
			if len(f) > 1 {
				// One more check: make sure server name is
				// just an IP address.  Otherwise we need DNS
				// to look it up.
				name := f[1]
				c.Servers = append(c.Servers, name)
			}

		case "domain": // set search path to just this domain
			if len(f) > 1 {
				c.Search = make([]string, 1)
				c.Search[0] = f[1]
			} else {
				c.Search = make([]string, 0)
			}

		case "search": // set search path to given servers
			c.Search = make([]string, len(f)-1)
			for i := 0; i < len(c.Search); i++ {
				c.Search[i] = f[i+1]
			}

		case "options": // magic options
			for i := 1; i < len(f); i++ {
				s := f[i]
				switch {
				case len(s) >= 6 && s[:6] == "ndots:":
					n, _ := strconv.Atoi(s[6:])
					if n < 1 {
						n = 1
					}
					c.Ndots = n
				case len(s) >= 8 && s[:8] == "timeout:":
					n, _ := strconv.Atoi(s[8:])
					if n < 1 {
						n = 1
					}
					c.Timeout = n
				case len(s) >= 8 && s[:9] == "attempts:":
					n, _ := strconv.Atoi(s[9:])
					if n < 1 {
						n = 1
					}
					c.Attempts = n
				case s == "rotate":
					/* not imp */
				}
			}
		}
	}
	return c, nil
}
