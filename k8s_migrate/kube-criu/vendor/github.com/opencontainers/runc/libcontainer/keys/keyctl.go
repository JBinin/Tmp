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

package keys

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

type KeySerial uint32

func JoinSessionKeyring(name string) (KeySerial, error) {
	sessKeyId, err := unix.KeyctlJoinSessionKeyring(name)
	if err != nil {
		return 0, fmt.Errorf("could not create session key: %v", err)
	}
	return KeySerial(sessKeyId), nil
}

// ModKeyringPerm modifies permissions on a keyring by reading the current permissions,
// anding the bits with the given mask (clearing permissions) and setting
// additional permission bits
func ModKeyringPerm(ringId KeySerial, mask, setbits uint32) error {
	dest, err := unix.KeyctlString(unix.KEYCTL_DESCRIBE, int(ringId))
	if err != nil {
		return err
	}

	res := strings.Split(dest, ";")
	if len(res) < 5 {
		return fmt.Errorf("Destination buffer for key description is too small")
	}

	// parse permissions
	perm64, err := strconv.ParseUint(res[3], 16, 32)
	if err != nil {
		return err
	}

	perm := (uint32(perm64) & mask) | setbits

	if err := unix.KeyctlSetperm(int(ringId), perm); err != nil {
		return err
	}

	return nil
}
