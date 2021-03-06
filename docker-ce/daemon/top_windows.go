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
package daemon

import (
	"errors"
	"fmt"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/go-units"
)

// ContainerTop handles `docker top` client requests.
// Future considerations:
// -- Windows users are far more familiar with CPU% total.
//    Further, users on Windows rarely see user/kernel CPU stats split.
//    The kernel returns everything in terms of 100ns. To obtain
//    CPU%, we could do something like docker stats does which takes two
//    samples, subtract the difference and do the maths. Unfortunately this
//    would slow the stat call down and require two kernel calls. So instead,
//    we do something similar to linux and display the CPU as combined HH:MM:SS.mmm.
// -- Perhaps we could add an argument to display "raw" stats
// -- "Memory" is an extremely overloaded term in Windows. Hence we do what
//    task manager does and use the private working set as the memory counter.
//    We could return more info for those who really understand how memory
//    management works in Windows if we introduced a "raw" stats (above).
func (daemon *Daemon) ContainerTop(name string, psArgs string) (*containertypes.ContainerTopOKBody, error) {
	// It's not at all an equivalent to linux 'ps' on Windows
	if psArgs != "" {
		return nil, errors.New("Windows does not support arguments to top")
	}

	container, err := daemon.GetContainer(name)
	if err != nil {
		return nil, err
	}

	s, err := daemon.containerd.Summary(container.ID)
	if err != nil {
		return nil, err
	}
	procList := &containertypes.ContainerTopOKBody{}
	procList.Titles = []string{"Name", "PID", "CPU", "Private Working Set"}

	for _, j := range s {
		d := time.Duration((j.KernelTime100ns + j.UserTime100ns) * 100) // Combined time in nanoseconds
		procList.Processes = append(procList.Processes, []string{
			j.ImageName,
			fmt.Sprint(j.ProcessId),
			fmt.Sprintf("%02d:%02d:%02d.%03d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60, int(d.Nanoseconds()/1000000)%1000),
			units.HumanSize(float64(j.MemoryWorkingSetPrivateBytes))})
	}
	return procList, nil
}
