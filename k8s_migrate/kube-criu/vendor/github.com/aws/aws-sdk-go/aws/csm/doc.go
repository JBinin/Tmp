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
// Package csm provides Client Side Monitoring (CSM) which enables sending metrics
// via UDP connection. Using the Start function will enable the reporting of
// metrics on a given port. If Start is called, with different parameters, again,
// a panic will occur.
//
// Pause can be called to pause any metrics publishing on a given port. Sessions
// that have had their handlers modified via InjectHandlers may still be used.
// However, the handlers will act as a no-op meaning no metrics will be published.
//
//	Example:
//		r, err := csm.Start("clientID", ":31000")
//		if err != nil {
//			panic(fmt.Errorf("failed starting CSM:  %v", err))
//		}
//
//		sess, err := session.NewSession(&aws.Config{})
//		if err != nil {
//			panic(fmt.Errorf("failed loading session: %v", err))
//		}
//
//		r.InjectHandlers(&sess.Handlers)
//
//		client := s3.New(sess)
//		resp, err := client.GetObject(&s3.GetObjectInput{
//			Bucket: aws.String("bucket"),
//			Key: aws.String("key"),
//		})
//
//		// Will pause monitoring
//		r.Pause()
//		resp, err = client.GetObject(&s3.GetObjectInput{
//			Bucket: aws.String("bucket"),
//			Key: aws.String("key"),
//		})
//
//		// Resume monitoring
//		r.Continue()
//
// Start returns a Reporter that is used to enable or disable monitoring. If
// access to the Reporter is required later, calling Get will return the Reporter
// singleton.
//
//	Example:
//		r := csm.Get()
//		r.Continue()
package csm
