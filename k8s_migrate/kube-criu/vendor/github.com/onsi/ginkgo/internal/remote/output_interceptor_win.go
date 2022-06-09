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
// +build windows

package remote

import (
	"errors"
)

func NewOutputInterceptor() OutputInterceptor {
	return &outputInterceptor{}
}

type outputInterceptor struct {
	intercepting bool
}

func (interceptor *outputInterceptor) StartInterceptingOutput() error {
	if interceptor.intercepting {
		return errors.New("Already intercepting output!")
	}
	interceptor.intercepting = true

	// not working on windows...

	return nil
}

func (interceptor *outputInterceptor) StopInterceptingAndReturnOutput() (string, error) {
	// not working on windows...
	interceptor.intercepting = false

	return "", nil
}
