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
package volume

import (
	"fmt"
	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/api/client"
	"github.com/libopenstorage/openstorage/volume"
)

// VolumeDriver returns a REST wrapper for the VolumeDriver interface.
func VolumeDriver(c *client.Client) volume.VolumeDriver {
	return newVolumeClient(c)
}

// NewAuthDriverClient returns a new REST client of the supplied version for specified driver.
// host: REST endpoint [http://<ip>:<port> OR unix://<path-to-unix-socket>]. default: [unix:///var/lib/osd/<driverName>.sock]
// version: Volume API version
func NewAuthDriverClient(host, driverName, version, authstring, accesstoken, userAgent string) (*client.Client, error) {
	if driverName == "" {
		return nil, fmt.Errorf("Driver Name cannot be empty")
	}
	if host == "" {
		host = client.GetUnixServerPath(driverName, volume.DriverAPIBase)
	}
	if version == "" {
		// Set the default version
		version = volume.APIVersion
	}
	return client.NewAuthClient(host, version, authstring, accesstoken, userAgent)
}

// NewDriverClient returns a new REST client of the supplied version for specified driver.
// host: REST endpoint [http://<ip>:<port> OR unix://<path-to-unix-socket>]. default: [unix:///var/lib/osd/<driverName>.sock]
// version: Volume API version
func NewDriverClient(host, driverName, version, userAgent string) (*client.Client, error) {
	if driverName == "" {
		return nil, fmt.Errorf("Driver Name cannot be empty")
	}
	if host == "" {
		host = client.GetUnixServerPath(driverName, volume.DriverAPIBase)
	}
	if version == "" {
		// Set the default version
		version = volume.APIVersion
	}
	return client.NewClient(host, version, userAgent)
}

// GetSupportedDriverVersions returns a list of supported versions
// for the provided driver. It uses the given server endpoint or the
// standard unix domain socket
func GetSupportedDriverVersions(driverName, host string) ([]string, error) {
	// Get a client handler
	if host == "" {
		host = client.GetUnixServerPath(driverName, volume.DriverAPIBase)
	}

	client, err := client.NewClient(host, "", "")
	if err != nil {
		return []string{}, err
	}
	versions, err := client.Versions(api.OsdVolumePath)
	if err != nil {
		return []string{}, err
	}
	return versions, nil
}
