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
package hcsshim

import (
	"time"

	"github.com/sirupsen/logrus"
)

func processAsyncHcsResult(err error, resultp *uint16, callbackNumber uintptr, expectedNotification hcsNotification, timeout *time.Duration) error {
	err = processHcsResult(err, resultp)
	if IsPending(err) {
		return waitForNotification(callbackNumber, expectedNotification, timeout)
	}

	return err
}

func waitForNotification(callbackNumber uintptr, expectedNotification hcsNotification, timeout *time.Duration) error {
	callbackMapLock.RLock()
	channels := callbackMap[callbackNumber].channels
	callbackMapLock.RUnlock()

	expectedChannel := channels[expectedNotification]
	if expectedChannel == nil {
		logrus.Errorf("unknown notification type in waitForNotification %x", expectedNotification)
		return ErrInvalidNotificationType
	}

	var c <-chan time.Time
	if timeout != nil {
		timer := time.NewTimer(*timeout)
		c = timer.C
		defer timer.Stop()
	}

	select {
	case err, ok := <-expectedChannel:
		if !ok {
			return ErrHandleClose
		}
		return err
	case err, ok := <-channels[hcsNotificationSystemExited]:
		if !ok {
			return ErrHandleClose
		}
		// If the expected notification is hcsNotificationSystemExited which of the two selects
		// chosen is random. Return the raw error if hcsNotificationSystemExited is expected
		if channels[hcsNotificationSystemExited] == expectedChannel {
			return err
		}
		return ErrUnexpectedContainerExit
	case _, ok := <-channels[hcsNotificationServiceDisconnect]:
		if !ok {
			return ErrHandleClose
		}
		// hcsNotificationServiceDisconnect should never be an expected notification
		// it does not need the same handling as hcsNotificationSystemExited
		return ErrUnexpectedProcessAbort
	case <-c:
		return ErrTimeout
	}
	return nil
}
