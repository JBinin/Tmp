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
package api

// MinVersion represents Minimum REST API version supported
// Technically the first daemon API version released on Windows is v1.25 in
// engine version 1.13. However, some clients are explicitly using downlevel
// APIs (e.g. docker-compose v2.1 file format) and that is just too restrictive.
// Hence also allowing 1.24 on Windows.
const MinVersion string = "1.24"
