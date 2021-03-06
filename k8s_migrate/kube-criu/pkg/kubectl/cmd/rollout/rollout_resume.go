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

package rollout

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/printers"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/kubectl/cmd/set"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/polymorphichelpers"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
)

// ResumeOptions is the start of the data required to perform the operation.  As new fields are added, add them here instead of
// referencing the cmd.Flags()
type ResumeOptions struct {
	PrintFlags *genericclioptions.PrintFlags
	ToPrinter  func(string) (printers.ResourcePrinter, error)

	Resources []string

	Builder          func() *resource.Builder
	Resumer          polymorphichelpers.ObjectResumerFunc
	Namespace        string
	EnforceNamespace bool

	resource.FilenameOptions
	genericclioptions.IOStreams
}

var (
	resume_long = templates.LongDesc(`
		Resume a paused resource

		Paused resources will not be reconciled by a controller. By resuming a
		resource, we allow it to be reconciled again.
		Currently only deployments support being resumed.`)

	resume_example = templates.Examples(`
		# Resume an already paused deployment
		kubectl rollout resume deployment/nginx`)
)

func NewRolloutResumeOptions(streams genericclioptions.IOStreams) *ResumeOptions {
	return &ResumeOptions{
		PrintFlags: genericclioptions.NewPrintFlags("resumed").WithTypeSetter(scheme.Scheme),
		IOStreams:  streams,
	}
}

func NewCmdRolloutResume(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewRolloutResumeOptions(streams)

	validArgs := []string{"deployment"}

	cmd := &cobra.Command{
		Use: "resume RESOURCE",
		DisableFlagsInUseLine: true,
		Short:   i18n.T("Resume a paused resource"),
		Long:    resume_long,
		Example: resume_example,
		Run: func(cmd *cobra.Command, args []string) {
			allErrs := []error{}
			err := o.Complete(f, cmd, args)
			if err != nil {
				allErrs = append(allErrs, err)
			}
			err = o.RunResume()
			if err != nil {
				allErrs = append(allErrs, err)
			}
			cmdutil.CheckErr(utilerrors.Flatten(utilerrors.NewAggregate(allErrs)))
		},
		ValidArgs: validArgs,
	}

	usage := "identifying the resource to get from a server."
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	o.PrintFlags.AddFlags(cmd)
	return cmd
}

func (o *ResumeOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(args) == 0 && cmdutil.IsFilenameSliceEmpty(o.Filenames) {
		return cmdutil.UsageErrorf(cmd, "%s", cmd.Use)
	}

	o.Resources = args

	o.Resumer = polymorphichelpers.ObjectResumerFn

	var err error
	o.Namespace, o.EnforceNamespace, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	o.ToPrinter = func(operation string) (printers.ResourcePrinter, error) {
		o.PrintFlags.NamePrintFlags.Operation = operation
		return o.PrintFlags.ToPrinter()
	}

	o.Builder = f.NewBuilder

	return nil
}

func (o ResumeOptions) RunResume() error {
	r := o.Builder().
		WithScheme(legacyscheme.Scheme).
		NamespaceParam(o.Namespace).DefaultNamespace().
		FilenameParam(o.EnforceNamespace, &o.FilenameOptions).
		ResourceTypeOrNameArgs(true, o.Resources...).
		ContinueOnError().
		Latest().
		Flatten().
		Do()
	if err := r.Err(); err != nil {
		return err
	}

	allErrs := []error{}
	infos, err := r.Infos()
	if err != nil {
		// restore previous command behavior where
		// an error caused by retrieving infos due to
		// at least a single broken object did not result
		// in an immediate return, but rather an overall
		// aggregation of errors.
		allErrs = append(allErrs, err)
	}

	for _, patch := range set.CalculatePatches(infos, cmdutil.InternalVersionJSONEncoder(), set.PatchFn(o.Resumer)) {
		info := patch.Info

		if patch.Err != nil {
			resourceString := info.Mapping.Resource.Resource
			if len(info.Mapping.Resource.Group) > 0 {
				resourceString = resourceString + "." + info.Mapping.Resource.Group
			}
			allErrs = append(allErrs, fmt.Errorf("error: %s %q %v", resourceString, info.Name, patch.Err))
			continue
		}

		if string(patch.Patch) == "{}" || len(patch.Patch) == 0 {
			printer, err := o.ToPrinter("already resumed")
			if err != nil {
				allErrs = append(allErrs, err)
				continue
			}
			if err = printer.PrintObj(cmdutil.AsDefaultVersionedOrOriginal(info.Object, info.Mapping), o.Out); err != nil {
				allErrs = append(allErrs, err)
			}
			continue
		}

		obj, err := resource.NewHelper(info.Client, info.Mapping).Patch(info.Namespace, info.Name, types.StrategicMergePatchType, patch.Patch)
		if err != nil {
			allErrs = append(allErrs, fmt.Errorf("failed to patch: %v", err))
			continue
		}

		info.Refresh(obj, true)
		printer, err := o.ToPrinter("resumed")
		if err != nil {
			allErrs = append(allErrs, err)
			continue
		}
		if err = printer.PrintObj(cmdutil.AsDefaultVersionedOrOriginal(info.Object, info.Mapping), o.Out); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return utilerrors.NewAggregate(allErrs)
}
