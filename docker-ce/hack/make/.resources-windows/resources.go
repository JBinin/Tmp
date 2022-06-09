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
/*

Package winresources is used to embed Windows resources into docker.exe.
These resources are used to provide

    * Version information
    * An icon
    * A Windows manifest declaring Windows version support

The resource object files are generated in hack/make/.go-autogen from
source files in hack/make/.resources-windows. This occurs automatically
when you run hack/make.sh.

These object files are picked up automatically by go build when this package
is included.

*/
package winresources
