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
package daemon

import (
	// Importing packages here only to make sure their init gets called and
	// therefore they register themselves to the logdriver factory.
	_ "github.com/docker/docker/daemon/logger/awslogs"
	_ "github.com/docker/docker/daemon/logger/etwlogs"
	_ "github.com/docker/docker/daemon/logger/fluentd"
	_ "github.com/docker/docker/daemon/logger/jsonfilelog"
	_ "github.com/docker/docker/daemon/logger/logentries"
	_ "github.com/docker/docker/daemon/logger/splunk"
	_ "github.com/docker/docker/daemon/logger/syslog"
)
