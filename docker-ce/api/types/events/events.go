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
package events

const (
	// ContainerEventType is the event type that containers generate
	ContainerEventType = "container"
	// DaemonEventType is the event type that daemon generate
	DaemonEventType = "daemon"
	// ImageEventType is the event type that images generate
	ImageEventType = "image"
	// NetworkEventType is the event type that networks generate
	NetworkEventType = "network"
	// PluginEventType is the event type that plugins generate
	PluginEventType = "plugin"
	// VolumeEventType is the event type that volumes generate
	VolumeEventType = "volume"
)

// Actor describes something that generates events,
// like a container, or a network, or a volume.
// It has a defined name and a set or attributes.
// The container attributes are its labels, other actors
// can generate these attributes from other properties.
type Actor struct {
	ID         string
	Attributes map[string]string
}

// Message represents the information an event contains
type Message struct {
	// Deprecated information from JSONMessage.
	// With data only in container events.
	Status string `json:"status,omitempty"`
	ID     string `json:"id,omitempty"`
	From   string `json:"from,omitempty"`

	Type   string
	Action string
	Actor  Actor

	Time     int64 `json:"time,omitempty"`
	TimeNano int64 `json:"timeNano,omitempty"`
}
