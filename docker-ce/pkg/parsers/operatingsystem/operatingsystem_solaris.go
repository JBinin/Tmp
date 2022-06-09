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
// +build solaris,cgo

package operatingsystem

/*
#include <zone.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"io/ioutil"
)

var etcOsRelease = "/etc/release"

// GetOperatingSystem gets the name of the current operating system.
func GetOperatingSystem() (string, error) {
	b, err := ioutil.ReadFile(etcOsRelease)
	if err != nil {
		return "", err
	}
	if i := bytes.Index(b, []byte("\n")); i >= 0 {
		b = bytes.Trim(b[:i], " ")
		return string(b), nil
	}
	return "", errors.New("release not found")
}

// IsContainerized returns true if we are running inside a container.
func IsContainerized() (bool, error) {
	if C.getzoneid() != 0 {
		return true, nil
	}
	return false, nil
}
