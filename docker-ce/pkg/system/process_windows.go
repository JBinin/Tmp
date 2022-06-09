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
package system

// IsProcessAlive returns true if process with a given pid is running.
func IsProcessAlive(pid int) bool {
	// TODO Windows containerd. Not sure this is needed
	//	p, err := os.FindProcess(pid)
	//	if err == nil {
	//		return true
	//	}
	return false
}

// KillProcess force-stops a process.
func KillProcess(pid int) {
	// TODO Windows containerd. Not sure this is needed
	//	p, err := os.FindProcess(pid)
	//	if err == nil {
	//		p.Kill()
	//	}
}
