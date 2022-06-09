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
// assetfs allows packages to serve static content embedded
// with the go-bindata tool with the standard net/http package.
//
// See https://github.com/jteeuwen/go-bindata for more information
// about embedding binary data with go-bindata.
//
// Usage example, after running
//    $ go-bindata data/...
// use:
//     http.Handle("/",
//        http.FileServer(
//        &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "data"}))
package assetfs
