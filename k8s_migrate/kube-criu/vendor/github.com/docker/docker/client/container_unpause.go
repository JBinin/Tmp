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
package client

import "golang.org/x/net/context"

// ContainerUnpause resumes the process execution within a container
func (cli *Client) ContainerUnpause(ctx context.Context, containerID string) error {
	resp, err := cli.post(ctx, "/containers/"+containerID+"/unpause", nil, nil, nil)
	ensureReaderClosed(resp)
	return err
}
