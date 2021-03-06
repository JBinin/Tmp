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
package osl

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// NeighOption is a function option type to set interface options
type NeighOption func(nh *neigh)

type neigh struct {
	dstIP    net.IP
	dstMac   net.HardwareAddr
	linkName string
	linkDst  string
	family   int
}

func (n *networkNamespace) findNeighbor(dstIP net.IP, dstMac net.HardwareAddr) *neigh {
	n.Lock()
	defer n.Unlock()

	for _, nh := range n.neighbors {
		if nh.dstIP.Equal(dstIP) && bytes.Equal(nh.dstMac, dstMac) {
			return nh
		}
	}

	return nil
}

func (n *networkNamespace) DeleteNeighbor(dstIP net.IP, dstMac net.HardwareAddr, osDelete bool) error {
	var (
		iface netlink.Link
		err   error
	)

	nh := n.findNeighbor(dstIP, dstMac)
	if nh == nil {
		return fmt.Errorf("could not find the neighbor entry to delete")
	}

	if osDelete {
		n.Lock()
		nlh := n.nlHandle
		n.Unlock()

		if nh.linkDst != "" {
			iface, err = nlh.LinkByName(nh.linkDst)
			if err != nil {
				return fmt.Errorf("could not find interface with destination name %s: %v",
					nh.linkDst, err)
			}
		}

		nlnh := &netlink.Neigh{
			IP:     dstIP,
			State:  netlink.NUD_PERMANENT,
			Family: nh.family,
		}

		if nlnh.Family > 0 {
			nlnh.HardwareAddr = dstMac
			nlnh.Flags = netlink.NTF_SELF
		}

		if nh.linkDst != "" {
			nlnh.LinkIndex = iface.Attrs().Index
		}

		if err := nlh.NeighDel(nlnh); err != nil {
			return fmt.Errorf("could not delete neighbor entry: %v", err)
		}
	}

	n.Lock()
	for i, nh := range n.neighbors {
		if nh.dstIP.Equal(dstIP) && bytes.Equal(nh.dstMac, dstMac) {
			n.neighbors = append(n.neighbors[:i], n.neighbors[i+1:]...)
			break
		}
	}
	n.Unlock()

	return nil
}

func (n *networkNamespace) AddNeighbor(dstIP net.IP, dstMac net.HardwareAddr, options ...NeighOption) error {
	var (
		iface netlink.Link
		err   error
	)

	nh := n.findNeighbor(dstIP, dstMac)
	if nh != nil {
		logrus.Debugf("Neighbor entry already present for IP %v, mac %v", dstIP, dstMac)
		// If it exists silently return
		return nil
	}

	nh = &neigh{
		dstIP:  dstIP,
		dstMac: dstMac,
	}

	nh.processNeighOptions(options...)

	if nh.linkName != "" {
		nh.linkDst = n.findDst(nh.linkName, false)
		if nh.linkDst == "" {
			return fmt.Errorf("could not find the interface with name %s", nh.linkName)
		}
	}

	n.Lock()
	nlh := n.nlHandle
	n.Unlock()

	if nh.linkDst != "" {
		iface, err = nlh.LinkByName(nh.linkDst)
		if err != nil {
			return fmt.Errorf("could not find interface with destination name %s: %v",
				nh.linkDst, err)
		}
	}

	nlnh := &netlink.Neigh{
		IP:           dstIP,
		HardwareAddr: dstMac,
		State:        netlink.NUD_PERMANENT,
		Family:       nh.family,
	}

	if nlnh.Family > 0 {
		nlnh.Flags = netlink.NTF_SELF
	}

	if nh.linkDst != "" {
		nlnh.LinkIndex = iface.Attrs().Index
	}

	if err := nlh.NeighSet(nlnh); err != nil {
		return fmt.Errorf("could not add neighbor entry: %v", err)
	}

	n.neighbors = append(n.neighbors, nh)

	return nil
}
