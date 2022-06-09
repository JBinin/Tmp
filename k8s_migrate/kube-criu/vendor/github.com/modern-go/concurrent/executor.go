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
package concurrent

import "context"

// Executor replace go keyword to start a new goroutine
// the goroutine should cancel itself if the context passed in has been cancelled
// the goroutine started by the executor, is owned by the executor
// we can cancel all executors owned by the executor just by stop the executor itself
// however Executor interface does not Stop method, the one starting and owning executor
// should use the concrete type of executor, instead of this interface.
type Executor interface {
	// Go starts a new goroutine controlled by the context
	Go(handler func(ctx context.Context))
}
