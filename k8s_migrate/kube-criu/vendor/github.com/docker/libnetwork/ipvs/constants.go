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
// +build linux

package ipvs

const (
	genlCtrlID = 0x10
)

// GENL control commands
const (
	genlCtrlCmdUnspec uint8 = iota
	genlCtrlCmdNewFamily
	genlCtrlCmdDelFamily
	genlCtrlCmdGetFamily
)

// GENL family attributes
const (
	genlCtrlAttrUnspec int = iota
	genlCtrlAttrFamilyID
	genlCtrlAttrFamilyName
)

// IPVS genl commands
const (
	ipvsCmdUnspec uint8 = iota
	ipvsCmdNewService
	ipvsCmdSetService
	ipvsCmdDelService
	ipvsCmdGetService
	ipvsCmdNewDest
	ipvsCmdSetDest
	ipvsCmdDelDest
	ipvsCmdGetDest
	ipvsCmdNewDaemon
	ipvsCmdDelDaemon
	ipvsCmdGetDaemon
	ipvsCmdSetConfig
	ipvsCmdGetConfig
	ipvsCmdSetInfo
	ipvsCmdGetInfo
	ipvsCmdZero
	ipvsCmdFlush
)

// Attributes used in the first level of commands
const (
	ipvsCmdAttrUnspec int = iota
	ipvsCmdAttrService
	ipvsCmdAttrDest
	ipvsCmdAttrDaemon
	ipvsCmdAttrTimeoutTCP
	ipvsCmdAttrTimeoutTCPFin
	ipvsCmdAttrTimeoutUDP
)

// Attributes used to describe a service. Used inside nested attribute
// ipvsCmdAttrService
const (
	ipvsSvcAttrUnspec int = iota
	ipvsSvcAttrAddressFamily
	ipvsSvcAttrProtocol
	ipvsSvcAttrAddress
	ipvsSvcAttrPort
	ipvsSvcAttrFWMark
	ipvsSvcAttrSchedName
	ipvsSvcAttrFlags
	ipvsSvcAttrTimeout
	ipvsSvcAttrNetmask
	ipvsSvcAttrStats
	ipvsSvcAttrPEName
)

// Attributes used to describe a destination (real server). Used
// inside nested attribute ipvsCmdAttrDest.
const (
	ipvsDestAttrUnspec int = iota
	ipvsDestAttrAddress
	ipvsDestAttrPort
	ipvsDestAttrForwardingMethod
	ipvsDestAttrWeight
	ipvsDestAttrUpperThreshold
	ipvsDestAttrLowerThreshold
	ipvsDestAttrActiveConnections
	ipvsDestAttrInactiveConnections
	ipvsDestAttrPersistentConnections
	ipvsDestAttrStats
	ipvsDestAttrAddressFamily
)

// IPVS Svc Statistics constancs

const (
	ipvsSvcStatsUnspec int = iota
	ipvsSvcStatsConns
	ipvsSvcStatsPktsIn
	ipvsSvcStatsPktsOut
	ipvsSvcStatsBytesIn
	ipvsSvcStatsBytesOut
	ipvsSvcStatsCPS
	ipvsSvcStatsPPSIn
	ipvsSvcStatsPPSOut
	ipvsSvcStatsBPSIn
	ipvsSvcStatsBPSOut
)

// Destination forwarding methods
const (
	// ConnectionFlagFwdmask indicates the mask in the connection
	// flags which is used by forwarding method bits.
	ConnectionFlagFwdMask = 0x0007

	// ConnectionFlagMasq is used for masquerade forwarding method.
	ConnectionFlagMasq = 0x0000

	// ConnectionFlagLocalNode is used for local node forwarding
	// method.
	ConnectionFlagLocalNode = 0x0001

	// ConnectionFlagTunnel is used for tunnel mode forwarding
	// method.
	ConnectionFlagTunnel = 0x0002

	// ConnectionFlagDirectRoute is used for direct routing
	// forwarding method.
	ConnectionFlagDirectRoute = 0x0003
)

const (
	// RoundRobin distributes jobs equally amongst the available
	// real servers.
	RoundRobin = "rr"

	// LeastConnection assigns more jobs to real servers with
	// fewer active jobs.
	LeastConnection = "lc"

	// DestinationHashing assigns jobs to servers through looking
	// up a statically assigned hash table by their destination IP
	// addresses.
	DestinationHashing = "dh"

	// SourceHashing assigns jobs to servers through looking up
	// a statically assigned hash table by their source IP
	// addresses.
	SourceHashing = "sh"
)
