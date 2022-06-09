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
package libcontainerd

// Remote on Linux defines the accesspoint to the containerd grpc API.
// Remote on Windows is largely an unimplemented interface as there is
// no remote containerd.
type Remote interface {
	// Client returns a new Client instance connected with given Backend.
	Client(Backend) (Client, error)
	// Cleanup stops containerd if it was started by libcontainerd.
	// Note this is not used on Windows as there is no remote containerd.
	Cleanup()
	// UpdateOptions allows various remote options to be updated at runtime.
	UpdateOptions(...RemoteOption) error
}

// RemoteOption allows to configure parameters of remotes.
// This is unused on Windows.
type RemoteOption interface {
	Apply(Remote) error
}
