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
package types

// VersionInfo describes version and runtime info.
type VersionInfo struct {
	Name          string `json:"name"`
	BuildDate     string `json:"buildDate"`
	Revision      string `json:"revision"`
	Version       string `json:"version"`
	APIVersion    string `json:"apiVersion"`
	GoVersion     string `json:"goVersion"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	KernelVersion string `json:"kernelVersion"`
	Experimental  bool   `json:"experimental"`
}

type VersionResponse struct {
	Client *VersionInfo
	Server *VersionInfo
}

// ServerOK returns true when the client could connect to the docker server
// and parse the information received. It returns false otherwise.
func (v VersionResponse) ServerOK() bool {
	return v.Server != nil
}
