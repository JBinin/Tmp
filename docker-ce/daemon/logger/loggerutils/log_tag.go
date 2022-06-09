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
package loggerutils

import (
	"bytes"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/pkg/templates"
)

// DefaultTemplate defines the defaults template logger should use.
const DefaultTemplate = "{{.ID}}"

// ParseLogTag generates a context aware tag for consistency across different
// log drivers based on the context of the running container.
func ParseLogTag(info logger.Info, defaultTemplate string) (string, error) {
	tagTemplate := info.Config["tag"]
	if tagTemplate == "" {
		tagTemplate = defaultTemplate
	}

	tmpl, err := templates.NewParse("log-tag", tagTemplate)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, &info); err != nil {
		return "", err
	}

	return buf.String(), nil
}
