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

// StatusKind indicates the severity of a status
type StatusKind int32

const (
	// StatusSeverityLow indicates an OK status
	StatusSeverityLow StatusKind = iota
	// StatusSeverityMedium indicates a status which is in transition from OK to BAD or vice versa
	StatusSeverityMedium
	// StatusSeverityHigh indicates a BAD status
	StatusSeverityHigh
)

var statusToStatusKind = map[Status]StatusKind{
	Status_STATUS_NONE:                  StatusSeverityHigh,
	Status_STATUS_INIT:                  StatusSeverityMedium,
	Status_STATUS_OK:                    StatusSeverityLow,
	Status_STATUS_OFFLINE:               StatusSeverityHigh,
	Status_STATUS_ERROR:                 StatusSeverityHigh,
	Status_STATUS_NOT_IN_QUORUM:         StatusSeverityHigh,
	Status_STATUS_DECOMMISSION:          StatusSeverityHigh,
	Status_STATUS_MAINTENANCE:           StatusSeverityHigh,
	Status_STATUS_STORAGE_DOWN:          StatusSeverityHigh,
	Status_STATUS_STORAGE_DEGRADED:      StatusSeverityHigh,
	Status_STATUS_NEEDS_REBOOT:          StatusSeverityHigh,
	Status_STATUS_STORAGE_REBALANCE:     StatusSeverityMedium,
	Status_STATUS_STORAGE_DRIVE_REPLACE: StatusSeverityMedium,
	// Add statuses before MAX
	Status_STATUS_MAX: StatusSeverityHigh,
}

// StatusSimpleValueOf returns the string format of Status
func StatusSimpleValueOf(s string) (Status, error) {
	obj, err := simpleValueOf("status", Status_value, s)
	return Status(obj), err
}

// SimpleString returns the string format of Status
func (x Status) SimpleString() string {
	return simpleString("status", Status_name, int32(x))
}

// StatusKind returns the king of status
func (x Status) StatusKind() StatusKind {
	statusType, _ := statusToStatusKind[x]
	return statusType
}

// StatusKindMapLength used only for unit testing
func StatusKindMapLength() int {
	return len(statusToStatusKind)
}
