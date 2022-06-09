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
package formatter

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/stringutils"
)

const (
	defaultPluginTableFormat = "table {{.ID}}\t{{.Name}}\t{{.Description}}\t{{.Enabled}}"

	pluginIDHeader    = "ID"
	descriptionHeader = "DESCRIPTION"
	enabledHeader     = "ENABLED"
)

// NewPluginFormat returns a Format for rendering using a plugin Context
func NewPluginFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultQuietFormat
		}
		return defaultPluginTableFormat
	case RawFormatKey:
		if quiet {
			return `plugin_id: {{.ID}}`
		}
		return `plugin_id: {{.ID}}\nname: {{.Name}}\ndescription: {{.Description}}\nenabled: {{.Enabled}}\n`
	}
	return Format(source)
}

// PluginWrite writes the context
func PluginWrite(ctx Context, plugins []*types.Plugin) error {
	render := func(format func(subContext subContext) error) error {
		for _, plugin := range plugins {
			pluginCtx := &pluginContext{trunc: ctx.Trunc, p: *plugin}
			if err := format(pluginCtx); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&pluginContext{}, render)
}

type pluginContext struct {
	HeaderContext
	trunc bool
	p     types.Plugin
}

func (c *pluginContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *pluginContext) ID() string {
	c.AddHeader(pluginIDHeader)
	if c.trunc {
		return stringid.TruncateID(c.p.ID)
	}
	return c.p.ID
}

func (c *pluginContext) Name() string {
	c.AddHeader(nameHeader)
	return c.p.Name
}

func (c *pluginContext) Description() string {
	c.AddHeader(descriptionHeader)
	desc := strings.Replace(c.p.Config.Description, "\n", "", -1)
	desc = strings.Replace(desc, "\r", "", -1)
	if c.trunc {
		desc = stringutils.Ellipsis(desc, 45)
	}

	return desc
}

func (c *pluginContext) Enabled() bool {
	c.AddHeader(enabledHeader)
	return c.p.Enabled
}

func (c *pluginContext) PluginReference() string {
	c.AddHeader(imageHeader)
	return c.p.PluginReference
}
