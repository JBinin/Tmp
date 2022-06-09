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

import "os"

// ProcessBaseLayer post-processes a base layer that has had its files extracted.
// The files should have been extracted to <path>\Files.
func ProcessBaseLayer(path string) error {
	err := processBaseImage(path)
	if err != nil {
		return &os.PathError{Op: "ProcessBaseLayer", Path: path, Err: err}
	}
	return nil
}

// ProcessUtilityVMImage post-processes a utility VM image that has had its files extracted.
// The files should have been extracted to <path>\Files.
func ProcessUtilityVMImage(path string) error {
	err := processUtilityImage(path)
	if err != nil {
		return &os.PathError{Op: "ProcessUtilityVMImage", Path: path, Err: err}
	}
	return nil
}
