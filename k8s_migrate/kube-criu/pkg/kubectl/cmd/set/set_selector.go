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
/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package set

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/printers"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
)

// SelectorOptions is the start of the data required to perform the operation.  As new fields are added, add them here instead of
// referencing the cmd.Flags()
type SetSelectorOptions struct {
	// Bound
	ResourceBuilderFlags *genericclioptions.ResourceBuilderFlags
	PrintFlags           *genericclioptions.PrintFlags
	RecordFlags          *genericclioptions.RecordFlags
	dryrun               bool

	// set by args
	resources []string
	selector  *metav1.LabelSelector

	// computed
	WriteToServer  bool
	PrintObj       printers.ResourcePrinterFunc
	Recorder       genericclioptions.Recorder
	ResourceFinder genericclioptions.ResourceFinder

	// set at initialization
	genericclioptions.IOStreams
}

var (
	selectorLong = templates.LongDesc(`
		Set the selector on a resource. Note that the new selector will overwrite the old selector if the resource had one prior to the invocation
		of 'set selector'.

		A selector must begin with a letter or number, and may contain letters, numbers, hyphens, dots, and underscores, up to %[1]d characters.
		If --resource-version is specified, then updates will use this resource version, otherwise the existing resource-version will be used.
        Note: currently selectors can only be set on Service objects.`)
	selectorExample = templates.Examples(`
        # set the labels and selector before creating a deployment/service pair.
        kubectl create service clusterip my-svc --clusterip="None" -o yaml --dry-run | kubectl set selector --local -f - 'environment=qa' -o yaml | kubectl create -f -
        kubectl create deployment my-dep -o yaml --dry-run | kubectl label --local -f - environment=qa -o yaml | kubectl create -f -`)
)

func NewSelectorOptions(streams genericclioptions.IOStreams) *SetSelectorOptions {
	return &SetSelectorOptions{
		ResourceBuilderFlags: genericclioptions.NewResourceBuilderFlags().
			WithScheme(scheme.Scheme).
			WithAll(false).
			WithLocal(false).
			WithUninitialized(false).
			WithLatest(),
		PrintFlags:  genericclioptions.NewPrintFlags("selector updated").WithTypeSetter(scheme.Scheme),
		RecordFlags: genericclioptions.NewRecordFlags(),

		Recorder: genericclioptions.NoopRecorder{},

		IOStreams: streams,
	}
}

// NewCmdSelector is the "set selector" command.
func NewCmdSelector(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewSelectorOptions(streams)

	cmd := &cobra.Command{
		Use: "selector (-f FILENAME | TYPE NAME) EXPRESSIONS [--resource-version=version]",
		DisableFlagsInUseLine: true,
		Short:   i18n.T("Set the selector on a resource"),
		Long:    fmt.Sprintf(selectorLong, validation.LabelValueMaxLength),
		Example: selectorExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.RunSelector())
		},
	}

	o.ResourceBuilderFlags.AddFlags(cmd.Flags())
	o.PrintFlags.AddFlags(cmd)
	o.RecordFlags.AddFlags(cmd)

	cmd.Flags().String("resource-version", "", "If non-empty, the selectors update will only succeed if this is the current resource-version for the object. Only valid when specifying a single resource.")
	cmdutil.AddDryRunFlag(cmd)

	return cmd
}

// Complete assigns the SelectorOptions from args.
func (o *SetSelectorOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error

	o.RecordFlags.Complete(cmd)
	o.Recorder, err = o.RecordFlags.ToRecorder()
	if err != nil {
		return err
	}

	o.dryrun = cmdutil.GetDryRunFlag(cmd)

	o.resources, o.selector, err = getResourcesAndSelector(args)
	if err != nil {
		return err
	}

	o.ResourceFinder = o.ResourceBuilderFlags.ToBuilder(f, o.resources)
	o.WriteToServer = !(*o.ResourceBuilderFlags.Local || o.dryrun)

	if o.dryrun {
		o.PrintFlags.Complete("%s (dry run)")
	}
	printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	o.PrintObj = printer.PrintObj

	return err
}

// Validate basic inputs
func (o *SetSelectorOptions) Validate() error {
	if o.selector == nil {
		return fmt.Errorf("one selector is required")
	}
	return nil
}

// RunSelector executes the command.
func (o *SetSelectorOptions) RunSelector() error {
	r := o.ResourceFinder.Do()

	return r.Visit(func(info *resource.Info, err error) error {
		patch := &Patch{Info: info}
		CalculatePatch(patch, scheme.DefaultJSONEncoder(), func(obj runtime.Object) ([]byte, error) {
			selectErr := updateSelectorForObject(info.Object, *o.selector)
			if selectErr != nil {
				return nil, selectErr
			}

			// record this change (for rollout history)
			if err := o.Recorder.Record(patch.Info.Object); err != nil {
				glog.V(4).Infof("error recording current command: %v", err)
			}

			return runtime.Encode(scheme.DefaultJSONEncoder(), info.Object)
		})

		if patch.Err != nil {
			return patch.Err
		}
		if !o.WriteToServer {
			return o.PrintObj(info.Object, o.Out)
		}

		actual, err := resource.NewHelper(info.Client, info.Mapping).Patch(info.Namespace, info.Name, types.StrategicMergePatchType, patch.Patch)
		if err != nil {
			return err
		}

		return o.PrintObj(actual, o.Out)
	})
}

func updateSelectorForObject(obj runtime.Object, selector metav1.LabelSelector) error {
	copyOldSelector := func() (map[string]string, error) {
		if len(selector.MatchExpressions) > 0 {
			return nil, fmt.Errorf("match expression %v not supported on this object", selector.MatchExpressions)
		}
		dst := make(map[string]string)
		for label, value := range selector.MatchLabels {
			dst[label] = value
		}
		return dst, nil
	}
	var err error
	switch t := obj.(type) {
	case *v1.Service:
		t.Spec.Selector, err = copyOldSelector()
	default:
		err = fmt.Errorf("setting a selector is only supported for Services")
	}
	return err
}

// getResourcesAndSelector retrieves resources and the selector expression from the given args (assuming selectors the last arg)
func getResourcesAndSelector(args []string) (resources []string, selector *metav1.LabelSelector, err error) {
	if len(args) == 0 {
		return []string{}, nil, nil
	}
	resources = args[:len(args)-1]
	selector, err = metav1.ParseToLabelSelector(args[len(args)-1])
	return resources, selector, err
}
