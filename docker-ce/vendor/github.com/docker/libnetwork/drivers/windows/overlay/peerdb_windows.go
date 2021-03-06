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
package overlay

import (
	"fmt"
	"net"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libnetwork/types"

	"github.com/Microsoft/hcsshim"
)

const ovPeerTable = "overlay_peer_table"

func (d *driver) peerAdd(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, updateDb bool) error {

	logrus.Debugf("WINOVERLAY: Enter peerAdd for ca ip %s with ca mac %s", peerIP.String(), peerMac.String())

	if err := validateID(nid, eid); err != nil {
		return err
	}

	n := d.network(nid)
	if n == nil {
		return nil
	}

	if updateDb {
		logrus.Info("WINOVERLAY: peerAdd: notifying HNS of the REMOTE endpoint")

		hnsEndpoint := &hcsshim.HNSEndpoint{
			Name:             eid,
			VirtualNetwork:   n.hnsId,
			MacAddress:       peerMac.String(),
			IPAddress:        peerIP,
			IsRemoteEndpoint: true,
		}

		paPolicy, err := json.Marshal(hcsshim.PaPolicy{
			Type: "PA",
			PA:   vtep.String(),
		})

		if err != nil {
			return err
		}

		hnsEndpoint.Policies = append(hnsEndpoint.Policies, paPolicy)

		configurationb, err := json.Marshal(hnsEndpoint)
		if err != nil {
			return err
		}

		// Temp: We have to create an endpoint object to keep track of the HNS ID for
		// this endpoint so that we can retrieve it later when the endpoint is deleted.
		// This seems unnecessary when we already have dockers EID. See if we can pass
		// the global EID to HNS to use as it's ID, rather than having each HNS assign
		// it's own local ID for the endpoint

		addr, err := types.ParseCIDR(peerIP.String() + "/32")
		if err != nil {
			return err
		}

		n.removeEndpointWithAddress(addr)

		hnsresponse, err := hcsshim.HNSEndpointRequest("POST", "", string(configurationb))
		if err != nil {
			return err
		}

		ep := &endpoint{
			id:        eid,
			nid:       nid,
			addr:      addr,
			mac:       peerMac,
			profileId: hnsresponse.Id,
			remote:    true,
		}

		n.addEndpoint(ep)
	}

	return nil
}

func (d *driver) peerDelete(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, updateDb bool) error {

	logrus.Infof("WINOVERLAY: Enter peerDelete for endpoint %s and peer ip %s", eid, peerIP.String())

	if err := validateID(nid, eid); err != nil {
		return err
	}

	n := d.network(nid)
	if n == nil {
		return nil
	}

	ep := n.endpoint(eid)
	if ep == nil {
		return fmt.Errorf("could not find endpoint with id %s", eid)
	}

	if updateDb {
		_, err := hcsshim.HNSEndpointRequest("DELETE", ep.profileId, "")
		if err != nil {
			return err
		}

		n.deleteEndpoint(eid)
	}

	return nil
}
