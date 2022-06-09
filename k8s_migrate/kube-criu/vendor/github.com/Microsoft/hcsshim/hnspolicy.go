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
package hcsshim

// Type of Request Support in ModifySystem
type PolicyType string

// RequestType const
const (
	Nat                  PolicyType = "NAT"
	ACL                  PolicyType = "ACL"
	PA                   PolicyType = "PA"
	VLAN                 PolicyType = "VLAN"
	VSID                 PolicyType = "VSID"
	VNet                 PolicyType = "VNET"
	L2Driver             PolicyType = "L2Driver"
	Isolation            PolicyType = "Isolation"
	QOS                  PolicyType = "QOS"
	OutboundNat          PolicyType = "OutBoundNAT"
	ExternalLoadBalancer PolicyType = "ELB"
	Route                PolicyType = "ROUTE"
)

type NatPolicy struct {
	Type         PolicyType `json:"Type"`
	Protocol     string
	InternalPort uint16
	ExternalPort uint16
}

type QosPolicy struct {
	Type                            PolicyType `json:"Type"`
	MaximumOutgoingBandwidthInBytes uint64
}

type IsolationPolicy struct {
	Type               PolicyType `json:"Type"`
	VLAN               uint
	VSID               uint
	InDefaultIsolation bool
}

type VlanPolicy struct {
	Type PolicyType `json:"Type"`
	VLAN uint
}

type VsidPolicy struct {
	Type PolicyType `json:"Type"`
	VSID uint
}

type PaPolicy struct {
	Type PolicyType `json:"Type"`
	PA   string     `json:"PA"`
}

type OutboundNatPolicy struct {
	Policy
	VIP        string   `json:"VIP,omitempty"`
	Exceptions []string `json:"ExceptionList,omitempty"`
}

type ActionType string
type DirectionType string
type RuleType string

const (
	Allow ActionType = "Allow"
	Block ActionType = "Block"

	In  DirectionType = "In"
	Out DirectionType = "Out"

	Host   RuleType = "Host"
	Switch RuleType = "Switch"
)

type ACLPolicy struct {
	Type            PolicyType `json:"Type"`
	Protocol        uint16
	InternalPort    uint16
	Action          ActionType
	Direction       DirectionType
	LocalAddresses  string
	RemoteAddresses string
	LocalPort       uint16
	RemotePort      uint16
	RuleType        RuleType `json:"RuleType,omitempty"`
	Priority        uint16
	ServiceName     string
}

type Policy struct {
	Type PolicyType `json:"Type"`
}
