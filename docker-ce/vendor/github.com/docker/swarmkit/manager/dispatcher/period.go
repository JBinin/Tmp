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
package dispatcher

import (
	"math/rand"
	"time"
)

type periodChooser struct {
	period  time.Duration
	epsilon time.Duration
	rand    *rand.Rand
}

func newPeriodChooser(period, eps time.Duration) *periodChooser {
	return &periodChooser{
		period:  period,
		epsilon: eps,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (pc *periodChooser) Choose() time.Duration {
	var adj int64
	if pc.epsilon > 0 {
		adj = rand.Int63n(int64(2*pc.epsilon)) - int64(pc.epsilon)
	}
	return pc.period + time.Duration(adj)
}
