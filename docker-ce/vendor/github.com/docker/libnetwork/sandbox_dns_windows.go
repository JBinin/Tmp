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

package libnetwork

import (
	"github.com/docker/libnetwork/etchosts"
)

// Stub implementations for DNS related functions

func (sb *sandbox) startResolver(bool) {
}

func (sb *sandbox) setupResolutionFiles() error {
	return nil
}

func (sb *sandbox) restorePath() {
}

func (sb *sandbox) updateHostsFile(ifaceIP string) error {
	return nil
}

func (sb *sandbox) addHostsEntries(recs []etchosts.Record) {

}

func (sb *sandbox) deleteHostsEntries(recs []etchosts.Record) {

}

func (sb *sandbox) updateDNS(ipv6Enabled bool) error {
	return nil
}
