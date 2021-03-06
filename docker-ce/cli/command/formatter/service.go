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
	"fmt"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	mounttypes "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/command/inspect"
	"github.com/docker/docker/pkg/stringid"
	units "github.com/docker/go-units"
)

const serviceInspectPrettyTemplate Format = `
ID:		{{.ID}}
Name:		{{.Name}}
{{- if .Labels }}
Labels:
{{- range $k, $v := .Labels }}
 {{ $k }}{{if $v }}={{ $v }}{{ end }}
{{- end }}{{ end }}
Service Mode:
{{- if .IsModeGlobal }}	Global
{{- else if .IsModeReplicated }}	Replicated
{{- if .ModeReplicatedReplicas }}
 Replicas:	{{ .ModeReplicatedReplicas }}
{{- end }}{{ end }}
{{- if .HasUpdateStatus }}
UpdateStatus:
 State:		{{ .UpdateStatusState }}
{{- if .HasUpdateStatusStarted }}
 Started:	{{ .UpdateStatusStarted }}
{{- end }}
{{- if .UpdateIsCompleted }}
 Completed:	{{ .UpdateStatusCompleted }}
{{- end }}
 Message:	{{ .UpdateStatusMessage }}
{{- end }}
Placement:
{{- if .TaskPlacementConstraints -}}
 Contraints:	{{ .TaskPlacementConstraints }}
{{- end }}
{{- if .HasUpdateConfig }}
UpdateConfig:
 Parallelism:	{{ .UpdateParallelism }}
{{- if .HasUpdateDelay}}
 Delay:		{{ .UpdateDelay }}
{{- end }}
 On failure:	{{ .UpdateOnFailure }}
{{- if .HasUpdateMonitor}}
 Monitoring Period: {{ .UpdateMonitor }}
{{- end }}
 Max failure ratio: {{ .UpdateMaxFailureRatio }}
{{- end }}
ContainerSpec:
 Image:		{{ .ContainerImage }}
{{- if .ContainerArgs }}
 Args:		{{ range $arg := .ContainerArgs }}{{ $arg }} {{ end }}
{{- end -}}
{{- if .ContainerEnv }}
 Env:		{{ range $env := .ContainerEnv }}{{ $env }} {{ end }}
{{- end -}}
{{- if .ContainerWorkDir }}
 Dir:		{{ .ContainerWorkDir }}
{{- end -}}
{{- if .ContainerUser }}
 User: {{ .ContainerUser }}
{{- end }}
{{- if .ContainerMounts }}
Mounts:
{{- end }}
{{- range $mount := .ContainerMounts }}
  Target = {{ $mount.Target }}
   Source = {{ $mount.Source }}
   ReadOnly = {{ $mount.ReadOnly }}
   Type = {{ $mount.Type }}
{{- end -}}
{{- if .HasResources }}
Resources:
{{- if .HasResourceReservations }}
 Reservations:
{{- if gt .ResourceReservationNanoCPUs 0.0 }}
  CPU:		{{ .ResourceReservationNanoCPUs }}
{{- end }}
{{- if .ResourceReservationMemory }}
  Memory:	{{ .ResourceReservationMemory }}
{{- end }}{{ end }}
{{- if .HasResourceLimits }}
 Limits:
{{- if gt .ResourceLimitsNanoCPUs 0.0 }}
  CPU:		{{ .ResourceLimitsNanoCPUs }}
{{- end }}
{{- if .ResourceLimitMemory }}
  Memory:	{{ .ResourceLimitMemory }}
{{- end }}{{ end }}{{ end }}
{{- if .Networks }}
Networks:
{{- range $network := .Networks }} {{ $network }}{{ end }} {{ end }}
Endpoint Mode:	{{ .EndpointMode }}
{{- if .Ports }}
Ports:
{{- range $port := .Ports }}
 PublishedPort = {{ $port.PublishedPort }}
  Protocol = {{ $port.Protocol }}
  TargetPort = {{ $port.TargetPort }}
  PublishMode = {{ $port.PublishMode }}
{{- end }} {{ end -}}
`

// NewServiceFormat returns a Format for rendering using a Context
func NewServiceFormat(source string) Format {
	switch source {
	case PrettyFormatKey:
		return serviceInspectPrettyTemplate
	default:
		return Format(strings.TrimPrefix(source, RawFormatKey))
	}
}

// ServiceInspectWrite renders the context for a list of services
func ServiceInspectWrite(ctx Context, refs []string, getRef inspect.GetRefFunc) error {
	if ctx.Format != serviceInspectPrettyTemplate {
		return inspect.Inspect(ctx.Output, refs, string(ctx.Format), getRef)
	}
	render := func(format func(subContext subContext) error) error {
		for _, ref := range refs {
			serviceI, _, err := getRef(ref)
			if err != nil {
				return err
			}
			service, ok := serviceI.(swarm.Service)
			if !ok {
				return fmt.Errorf("got wrong object to inspect")
			}
			if err := format(&serviceInspectContext{Service: service}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&serviceInspectContext{}, render)
}

type serviceInspectContext struct {
	swarm.Service
	subContext
}

func (ctx *serviceInspectContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(ctx)
}

func (ctx *serviceInspectContext) ID() string {
	return ctx.Service.ID
}

func (ctx *serviceInspectContext) Name() string {
	return ctx.Service.Spec.Name
}

func (ctx *serviceInspectContext) Labels() map[string]string {
	return ctx.Service.Spec.Labels
}

func (ctx *serviceInspectContext) IsModeGlobal() bool {
	return ctx.Service.Spec.Mode.Global != nil
}

func (ctx *serviceInspectContext) IsModeReplicated() bool {
	return ctx.Service.Spec.Mode.Replicated != nil
}

func (ctx *serviceInspectContext) ModeReplicatedReplicas() *uint64 {
	return ctx.Service.Spec.Mode.Replicated.Replicas
}

func (ctx *serviceInspectContext) HasUpdateStatus() bool {
	return ctx.Service.UpdateStatus != nil && ctx.Service.UpdateStatus.State != ""
}

func (ctx *serviceInspectContext) UpdateStatusState() swarm.UpdateState {
	return ctx.Service.UpdateStatus.State
}

func (ctx *serviceInspectContext) HasUpdateStatusStarted() bool {
	return ctx.Service.UpdateStatus.StartedAt != nil
}

func (ctx *serviceInspectContext) UpdateStatusStarted() string {
	return units.HumanDuration(time.Since(*ctx.Service.UpdateStatus.StartedAt))
}

func (ctx *serviceInspectContext) UpdateIsCompleted() bool {
	return ctx.Service.UpdateStatus.State == swarm.UpdateStateCompleted && ctx.Service.UpdateStatus.CompletedAt != nil
}

func (ctx *serviceInspectContext) UpdateStatusCompleted() string {
	return units.HumanDuration(time.Since(*ctx.Service.UpdateStatus.CompletedAt))
}

func (ctx *serviceInspectContext) UpdateStatusMessage() string {
	return ctx.Service.UpdateStatus.Message
}

func (ctx *serviceInspectContext) TaskPlacementConstraints() []string {
	if ctx.Service.Spec.TaskTemplate.Placement != nil {
		return ctx.Service.Spec.TaskTemplate.Placement.Constraints
	}
	return nil
}

func (ctx *serviceInspectContext) HasUpdateConfig() bool {
	return ctx.Service.Spec.UpdateConfig != nil
}

func (ctx *serviceInspectContext) UpdateParallelism() uint64 {
	return ctx.Service.Spec.UpdateConfig.Parallelism
}

func (ctx *serviceInspectContext) HasUpdateDelay() bool {
	return ctx.Service.Spec.UpdateConfig.Delay.Nanoseconds() > 0
}

func (ctx *serviceInspectContext) UpdateDelay() time.Duration {
	return ctx.Service.Spec.UpdateConfig.Delay
}

func (ctx *serviceInspectContext) UpdateOnFailure() string {
	return ctx.Service.Spec.UpdateConfig.FailureAction
}

func (ctx *serviceInspectContext) HasUpdateMonitor() bool {
	return ctx.Service.Spec.UpdateConfig.Monitor.Nanoseconds() > 0
}

func (ctx *serviceInspectContext) UpdateMonitor() time.Duration {
	return ctx.Service.Spec.UpdateConfig.Monitor
}

func (ctx *serviceInspectContext) UpdateMaxFailureRatio() float32 {
	return ctx.Service.Spec.UpdateConfig.MaxFailureRatio
}

func (ctx *serviceInspectContext) ContainerImage() string {
	return ctx.Service.Spec.TaskTemplate.ContainerSpec.Image
}

func (ctx *serviceInspectContext) ContainerArgs() []string {
	return ctx.Service.Spec.TaskTemplate.ContainerSpec.Args
}

func (ctx *serviceInspectContext) ContainerEnv() []string {
	return ctx.Service.Spec.TaskTemplate.ContainerSpec.Env
}

func (ctx *serviceInspectContext) ContainerWorkDir() string {
	return ctx.Service.Spec.TaskTemplate.ContainerSpec.Dir
}

func (ctx *serviceInspectContext) ContainerUser() string {
	return ctx.Service.Spec.TaskTemplate.ContainerSpec.User
}

func (ctx *serviceInspectContext) ContainerMounts() []mounttypes.Mount {
	return ctx.Service.Spec.TaskTemplate.ContainerSpec.Mounts
}

func (ctx *serviceInspectContext) HasResources() bool {
	return ctx.Service.Spec.TaskTemplate.Resources != nil
}

func (ctx *serviceInspectContext) HasResourceReservations() bool {
	if ctx.Service.Spec.TaskTemplate.Resources == nil || ctx.Service.Spec.TaskTemplate.Resources.Reservations == nil {
		return false
	}
	return ctx.Service.Spec.TaskTemplate.Resources.Reservations.NanoCPUs > 0 || ctx.Service.Spec.TaskTemplate.Resources.Reservations.MemoryBytes > 0
}

func (ctx *serviceInspectContext) ResourceReservationNanoCPUs() float64 {
	if ctx.Service.Spec.TaskTemplate.Resources.Reservations.NanoCPUs == 0 {
		return float64(0)
	}
	return float64(ctx.Service.Spec.TaskTemplate.Resources.Reservations.NanoCPUs) / 1e9
}

func (ctx *serviceInspectContext) ResourceReservationMemory() string {
	if ctx.Service.Spec.TaskTemplate.Resources.Reservations.MemoryBytes == 0 {
		return ""
	}
	return units.BytesSize(float64(ctx.Service.Spec.TaskTemplate.Resources.Reservations.MemoryBytes))
}

func (ctx *serviceInspectContext) HasResourceLimits() bool {
	if ctx.Service.Spec.TaskTemplate.Resources == nil || ctx.Service.Spec.TaskTemplate.Resources.Limits == nil {
		return false
	}
	return ctx.Service.Spec.TaskTemplate.Resources.Limits.NanoCPUs > 0 || ctx.Service.Spec.TaskTemplate.Resources.Limits.MemoryBytes > 0
}

func (ctx *serviceInspectContext) ResourceLimitsNanoCPUs() float64 {
	return float64(ctx.Service.Spec.TaskTemplate.Resources.Limits.NanoCPUs) / 1e9
}

func (ctx *serviceInspectContext) ResourceLimitMemory() string {
	if ctx.Service.Spec.TaskTemplate.Resources.Limits.MemoryBytes == 0 {
		return ""
	}
	return units.BytesSize(float64(ctx.Service.Spec.TaskTemplate.Resources.Limits.MemoryBytes))
}

func (ctx *serviceInspectContext) Networks() []string {
	var out []string
	for _, n := range ctx.Service.Spec.Networks {
		out = append(out, n.Target)
	}
	return out
}

func (ctx *serviceInspectContext) EndpointMode() string {
	if ctx.Service.Spec.EndpointSpec == nil {
		return ""
	}

	return string(ctx.Service.Spec.EndpointSpec.Mode)
}

func (ctx *serviceInspectContext) Ports() []swarm.PortConfig {
	return ctx.Service.Endpoint.Ports
}

const (
	defaultServiceTableFormat = "table {{.ID}}\t{{.Name}}\t{{.Mode}}\t{{.Replicas}}\t{{.Image}}"

	serviceIDHeader = "ID"
	modeHeader      = "MODE"
	replicasHeader  = "REPLICAS"
)

// NewServiceListFormat returns a Format for rendering using a service Context
func NewServiceListFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultQuietFormat
		}
		return defaultServiceTableFormat
	case RawFormatKey:
		if quiet {
			return `id: {{.ID}}`
		}
		return `id: {{.ID}}\nname: {{.Name}}\nmode: {{.Mode}}\nreplicas: {{.Replicas}}\nimage: {{.Image}}\n`
	}
	return Format(source)
}

// ServiceListInfo stores the information about mode and replicas to be used by template
type ServiceListInfo struct {
	Mode     string
	Replicas string
}

// ServiceListWrite writes the context
func ServiceListWrite(ctx Context, services []swarm.Service, info map[string]ServiceListInfo) error {
	render := func(format func(subContext subContext) error) error {
		for _, service := range services {
			serviceCtx := &serviceContext{service: service, mode: info[service.ID].Mode, replicas: info[service.ID].Replicas}
			if err := format(serviceCtx); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&serviceContext{}, render)
}

type serviceContext struct {
	HeaderContext
	service  swarm.Service
	mode     string
	replicas string
}

func (c *serviceContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *serviceContext) ID() string {
	c.AddHeader(serviceIDHeader)
	return stringid.TruncateID(c.service.ID)
}

func (c *serviceContext) Name() string {
	c.AddHeader(nameHeader)
	return c.service.Spec.Name
}

func (c *serviceContext) Mode() string {
	c.AddHeader(modeHeader)
	return c.mode
}

func (c *serviceContext) Replicas() string {
	c.AddHeader(replicasHeader)
	return c.replicas
}

func (c *serviceContext) Image() string {
	c.AddHeader(imageHeader)
	image := c.service.Spec.TaskTemplate.ContainerSpec.Image
	if ref, err := reference.ParseNormalizedNamed(image); err == nil {
		// update image string for display, (strips any digest)
		if nt, ok := ref.(reference.NamedTagged); ok {
			if namedTagged, err := reference.WithTag(reference.TrimNamed(nt), nt.Tag()); err == nil {
				image = reference.FamiliarString(namedTagged)
			}
		}
	}

	return image
}
