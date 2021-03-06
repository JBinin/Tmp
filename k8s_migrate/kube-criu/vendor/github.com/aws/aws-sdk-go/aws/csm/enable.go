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
package csm

import (
	"fmt"
	"sync"
)

var (
	lock sync.Mutex
)

// Client side metric handler names
const (
	APICallMetricHandlerName        = "awscsm.SendAPICallMetric"
	APICallAttemptMetricHandlerName = "awscsm.SendAPICallAttemptMetric"
)

// Start will start the a long running go routine to capture
// client side metrics. Calling start multiple time will only
// start the metric listener once and will panic if a different
// client ID or port is passed in.
//
//	Example:
//		r, err := csm.Start("clientID", "127.0.0.1:8094")
//		if err != nil {
//			panic(fmt.Errorf("expected no error, but received %v", err))
//		}
//		sess := session.NewSession()
//		r.InjectHandlers(sess.Handlers)
//
//		svc := s3.New(sess)
//		out, err := svc.GetObject(&s3.GetObjectInput{
//			Bucket: aws.String("bucket"),
//			Key: aws.String("key"),
//		})
func Start(clientID string, url string) (*Reporter, error) {
	lock.Lock()
	defer lock.Unlock()

	if sender == nil {
		sender = newReporter(clientID, url)
	} else {
		if sender.clientID != clientID {
			panic(fmt.Errorf("inconsistent client IDs. %q was expected, but received %q", sender.clientID, clientID))
		}

		if sender.url != url {
			panic(fmt.Errorf("inconsistent URLs. %q was expected, but received %q", sender.url, url))
		}
	}

	if err := connect(url); err != nil {
		sender = nil
		return nil, err
	}

	return sender, nil
}

// Get will return a reporter if one exists, if one does not exist, nil will
// be returned.
func Get() *Reporter {
	lock.Lock()
	defer lock.Unlock()

	return sender
}
