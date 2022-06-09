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
package types

// ErrCapacityStatsUnchanged can be used when comparing stats
const ErrCapacityStatsUnchanged = "no changes"

// CapacityStats is used to report capacity statistics on pools and controllers.
type CapacityStats struct {

	// TotalCapacityBytes is the object's total capacity in bytes.
	TotalCapacityBytes uint64 `json:"totalCapacityBytes"`

	// AvailableCapacityBytes is the object's available capacity in bytes.
	AvailableCapacityBytes uint64 `json:"availableCapacityBytes"`

	// ProvisionedCapacityBytes is the object's provisioned capacity in bytes.
	ProvisionedCapacityBytes uint64 `json:"provisionedCapacityBytes"`
}

// IsEqual checks if capacity values are the same
func (c CapacityStats) IsEqual(n CapacityStats) bool {
	if c == n {
		return true
	}
	return false
}
