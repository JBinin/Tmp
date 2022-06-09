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
package main

import "errors"

var (
	errorLoadingDeps         = errors.New("error loading dependencies")
	errorLoadingPackages     = errors.New("error loading packages")
	errorCopyingSourceCode   = errors.New("error copying source code")
	errorNoPackagesUpdatable = errors.New("no packages can be updated")
)

type errPackageNotFound struct {
	path string
}

func (e errPackageNotFound) Error() string {
	return "Package (" + e.path + ") not found"
}
