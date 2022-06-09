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

package daemon

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	rootKeyFile   = "/proc/sys/kernel/keys/root_maxkeys"
	rootBytesFile = "/proc/sys/kernel/keys/root_maxbytes"
	rootKeyLimit  = 1000000
	// it is standard configuration to allocate 25 bytes per key
	rootKeyByteMultiplier = 25
)

// ModifyRootKeyLimit checks to see if the root key limit is set to
// at least 1000000 and changes it to that limit along with the maxbytes
// allocated to the keys at a 25 to 1 multiplier.
func ModifyRootKeyLimit() error {
	value, err := readRootKeyLimit(rootKeyFile)
	if err != nil {
		return err
	}
	if value < rootKeyLimit {
		return setRootKeyLimit(rootKeyLimit)
	}
	return nil
}

func setRootKeyLimit(limit int) error {
	keys, err := os.OpenFile(rootKeyFile, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer keys.Close()
	if _, err := fmt.Fprintf(keys, "%d", limit); err != nil {
		return err
	}
	bytes, err := os.OpenFile(rootBytesFile, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer bytes.Close()
	_, err = fmt.Fprintf(bytes, "%d", limit*rootKeyByteMultiplier)
	return err
}

func readRootKeyLimit(path string) (int, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(strings.Trim(string(data), "\n"))
}
