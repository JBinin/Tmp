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
package nodes

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/pkg/discovery"
)

// Discovery is exported
type Discovery struct {
	entries discovery.Entries
}

func init() {
	Init()
}

// Init is exported
func Init() {
	discovery.Register("nodes", &Discovery{})
}

// Initialize is exported
func (s *Discovery) Initialize(uris string, _ time.Duration, _ time.Duration, _ map[string]string) error {
	for _, input := range strings.Split(uris, ",") {
		for _, ip := range discovery.Generate(input) {
			entry, err := discovery.NewEntry(ip)
			if err != nil {
				return fmt.Errorf("%s, please check you are using the correct discovery (missing token:// ?)", err.Error())
			}
			s.entries = append(s.entries, entry)
		}
	}

	return nil
}

// Watch is exported
func (s *Discovery) Watch(stopCh <-chan struct{}) (<-chan discovery.Entries, <-chan error) {
	ch := make(chan discovery.Entries)
	go func() {
		defer close(ch)
		ch <- s.entries
		<-stopCh
	}()
	return ch, nil
}

// Register is exported
func (s *Discovery) Register(addr string) error {
	return discovery.ErrNotImplemented
}
