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
package term

import (
	"fmt"
	"strings"
)

// ASCII list the possible supported ASCII key sequence
var ASCII = []string{
	"ctrl-@",
	"ctrl-a",
	"ctrl-b",
	"ctrl-c",
	"ctrl-d",
	"ctrl-e",
	"ctrl-f",
	"ctrl-g",
	"ctrl-h",
	"ctrl-i",
	"ctrl-j",
	"ctrl-k",
	"ctrl-l",
	"ctrl-m",
	"ctrl-n",
	"ctrl-o",
	"ctrl-p",
	"ctrl-q",
	"ctrl-r",
	"ctrl-s",
	"ctrl-t",
	"ctrl-u",
	"ctrl-v",
	"ctrl-w",
	"ctrl-x",
	"ctrl-y",
	"ctrl-z",
	"ctrl-[",
	"ctrl-\\",
	"ctrl-]",
	"ctrl-^",
	"ctrl-_",
}

// ToBytes converts a string representing a suite of key-sequence to the corresponding ASCII code.
func ToBytes(keys string) ([]byte, error) {
	codes := []byte{}
next:
	for _, key := range strings.Split(keys, ",") {
		if len(key) != 1 {
			for code, ctrl := range ASCII {
				if ctrl == key {
					codes = append(codes, byte(code))
					continue next
				}
			}
			if key == "DEL" {
				codes = append(codes, 127)
			} else {
				return nil, fmt.Errorf("Unknown character: '%s'", key)
			}
		} else {
			codes = append(codes, byte(key[0]))
		}
	}
	return codes, nil
}
