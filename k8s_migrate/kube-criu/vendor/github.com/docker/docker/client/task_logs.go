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

import (
	"io"
	"net/url"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	timetypes "github.com/docker/docker/api/types/time"
)

// TaskLogs returns the logs generated by a task in an io.ReadCloser.
// It's up to the caller to close the stream.
func (cli *Client) TaskLogs(ctx context.Context, taskID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	query := url.Values{}
	if options.ShowStdout {
		query.Set("stdout", "1")
	}

	if options.ShowStderr {
		query.Set("stderr", "1")
	}

	if options.Since != "" {
		ts, err := timetypes.GetTimestamp(options.Since, time.Now())
		if err != nil {
			return nil, err
		}
		query.Set("since", ts)
	}

	if options.Timestamps {
		query.Set("timestamps", "1")
	}

	if options.Details {
		query.Set("details", "1")
	}

	if options.Follow {
		query.Set("follow", "1")
	}
	query.Set("tail", options.Tail)

	resp, err := cli.get(ctx, "/tasks/"+taskID+"/logs", query, nil)
	if err != nil {
		return nil, err
	}
	return resp.body, nil
}
