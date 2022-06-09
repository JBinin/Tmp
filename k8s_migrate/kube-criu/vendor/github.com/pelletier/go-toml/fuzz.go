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
// +build gofuzz

package toml

func Fuzz(data []byte) int {
	tree, err := LoadBytes(data)
	if err != nil {
		if tree != nil {
			panic("tree must be nil if there is an error")
		}
		return 0
	}

	str, err := tree.ToTomlString()
	if err != nil {
		if str != "" {
			panic(`str must be "" if there is an error`)
		}
		panic(err)
	}

	tree, err = Load(str)
	if err != nil {
		if tree != nil {
			panic("tree must be nil if there is an error")
		}
		return 0
	}

	return 1
}
