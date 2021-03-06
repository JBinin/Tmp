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
Copyright 2014 The Kubernetes Authors.

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

package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/printers"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	scheme "k8s.io/kubernetes/pkg/api/legacyscheme"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
	"k8s.io/kubernetes/pkg/kubectl/validation"
)

var (
	convert_long = templates.LongDesc(i18n.T(`
		Convert config files between different API versions. Both YAML
		and JSON formats are accepted.

		The command takes filename, directory, or URL as input, and convert it into format
		of version specified by --output-version flag. If target version is not specified or
		not supported, convert to latest version.

		The default output will be printed to stdout in YAML format. One can use -o option
		to change to output destination.`))

	convert_example = templates.Examples(i18n.T(`
		# Convert 'pod.yaml' to latest version and print to stdout.
		kubectl convert -f pod.yaml

		# Convert the live state of the resource specified by 'pod.yaml' to the latest version
		# and print to stdout in JSON format.
		kubectl convert -f pod.yaml --local -o json

		# Convert all files under current directory to latest version and create them all.
		kubectl convert -f . | kubectl create -f -`))
)

// ConvertOptions have the data required to perform the convert operation
type ConvertOptions struct {
	PrintFlags *genericclioptions.PrintFlags
	Printer    printers.ResourcePrinter

	OutputVersion string
	Namespace     string

	builder   func() *resource.Builder
	local     bool
	validator func() (validation.Schema, error)

	resource.FilenameOptions
	genericclioptions.IOStreams
}

func NewConvertOptions(ioStreams genericclioptions.IOStreams) *ConvertOptions {
	return &ConvertOptions{
		PrintFlags: genericclioptions.NewPrintFlags("converted").WithTypeSetter(scheme.Scheme).WithDefaultOutput("yaml"),
		local:      true,
		IOStreams:  ioStreams,
	}
}

// NewCmdConvert creates a command object for the generic "convert" action, which
// translates the config file into a given version.
func NewCmdConvert(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewConvertOptions(ioStreams)

	cmd := &cobra.Command{
		Use: "convert -f FILENAME",
		DisableFlagsInUseLine: true,
		Short:   i18n.T("Convert config files between different API versions"),
		Long:    convert_long,
		Example: convert_example,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.RunConvert())
		},
	}

	cmd.Flags().BoolVar(&o.local, "local", o.local, "If true, convert will NOT try to contact api-server but run locally.")
	cmd.Flags().StringVar(&o.OutputVersion, "output-version", o.OutputVersion, i18n.T("Output the formatted object with the given group version (for ex: 'extensions/v1beta1')."))
	o.PrintFlags.AddFlags(cmd)

	cmdutil.AddValidateFlags(cmd)
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, "to need to get converted.")
	cmd.MarkFlagRequired("filename")
	return cmd
}

// Complete collects information required to run Convert command from command line.
func (o *ConvertOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) (err error) {
	o.builder = f.NewBuilder

	o.Namespace, _, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	o.validator = func() (validation.Schema, error) {
		return f.Validator(cmdutil.GetFlagBool(cmd, "validate"))
	}

	// build the printer
	o.Printer, err = o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	return nil
}

// RunConvert implements the generic Convert command
func (o *ConvertOptions) RunConvert() error {
	b := o.builder().
		WithScheme(scheme.Scheme).
		LocalParam(o.local)
	if !o.local {
		schema, err := o.validator()
		if err != nil {
			return err
		}
		b.Schema(schema)
	}

	r := b.NamespaceParam(o.Namespace).
		ContinueOnError().
		FilenameParam(false, &o.FilenameOptions).
		Flatten().
		Do()

	err := r.Err()
	if err != nil {
		return err
	}

	singleItemImplied := false
	infos, err := r.IntoSingleItemImplied(&singleItemImplied).Infos()
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		return fmt.Errorf("no objects passed to convert")
	}

	var specifiedOutputVersion schema.GroupVersion
	if len(o.OutputVersion) > 0 {
		specifiedOutputVersion, err = schema.ParseGroupVersion(o.OutputVersion)
		if err != nil {
			return err
		}
	}

	objects, err := asVersionedObject(infos, !singleItemImplied, specifiedOutputVersion, cmdutil.InternalVersionJSONEncoder())
	if err != nil {
		return err
	}

	return o.Printer.PrintObj(objects, o.Out)
}

// asVersionedObject converts a list of infos into a single object - either a List containing
// the objects as children, or if only a single Object is present, as that object. The provided
// version will be preferred as the conversion target, but the Object's mapping version will be
// used if that version is not present.
func asVersionedObject(infos []*resource.Info, forceList bool, specifiedOutputVersion schema.GroupVersion, encoder runtime.Encoder) (runtime.Object, error) {
	objects, err := asVersionedObjects(infos, specifiedOutputVersion, encoder)
	if err != nil {
		return nil, err
	}

	var object runtime.Object
	if len(objects) == 1 && !forceList {
		object = objects[0]
	} else {
		object = &api.List{Items: objects}
		targetVersions := []schema.GroupVersion{}
		if !specifiedOutputVersion.Empty() {
			targetVersions = append(targetVersions, specifiedOutputVersion)
		}
		targetVersions = append(targetVersions, schema.GroupVersion{Group: "", Version: "v1"})

		converted, err := tryConvert(scheme.Scheme, object, targetVersions...)
		if err != nil {
			return nil, err
		}
		object = converted
	}

	actualVersion := object.GetObjectKind().GroupVersionKind()
	if actualVersion.Version != specifiedOutputVersion.Version {
		defaultVersionInfo := ""
		if len(actualVersion.Version) > 0 {
			defaultVersionInfo = fmt.Sprintf("Defaulting to %q", actualVersion.Version)
		}
		glog.V(1).Infof("info: the output version specified is invalid. %s\n", defaultVersionInfo)
	}
	return object, nil
}

// asVersionedObjects converts a list of infos into versioned objects. The provided
// version will be preferred as the conversion target, but the Object's mapping version will be
// used if that version is not present.
func asVersionedObjects(infos []*resource.Info, specifiedOutputVersion schema.GroupVersion, encoder runtime.Encoder) ([]runtime.Object, error) {
	objects := []runtime.Object{}
	for _, info := range infos {
		if info.Object == nil {
			continue
		}

		targetVersions := []schema.GroupVersion{}
		// objects that are not part of api.Scheme must be converted to JSON
		// TODO: convert to map[string]interface{}, attach to runtime.Unknown?
		if !specifiedOutputVersion.Empty() {
			if _, _, err := scheme.Scheme.ObjectKinds(info.Object); runtime.IsNotRegisteredError(err) {
				// TODO: ideally this would encode to version, but we don't expose multiple codecs here.
				data, err := runtime.Encode(encoder, info.Object)
				if err != nil {
					return nil, err
				}
				// TODO: Set ContentEncoding and ContentType.
				objects = append(objects, &runtime.Unknown{Raw: data})
				continue
			}
			targetVersions = append(targetVersions, specifiedOutputVersion)
		} else {
			gvks, _, err := scheme.Scheme.ObjectKinds(info.Object)
			if err == nil {
				for _, gvk := range gvks {
					for _, version := range scheme.Scheme.PrioritizedVersionsForGroup(gvk.Group) {
						targetVersions = append(targetVersions, version)
					}
				}
			}
		}

		converted, err := tryConvert(scheme.Scheme, info.Object, targetVersions...)
		if err != nil {
			return nil, err
		}
		objects = append(objects, converted)
	}
	return objects, nil
}

// tryConvert attempts to convert the given object to the provided versions in order. This function assumes
// the object is in internal version.
func tryConvert(converter runtime.ObjectConvertor, object runtime.Object, versions ...schema.GroupVersion) (runtime.Object, error) {
	var last error
	for _, version := range versions {
		if version.Empty() {
			return object, nil
		}
		obj, err := converter.ConvertToVersion(object, version)
		if err != nil {
			last = err
			continue
		}
		return obj, nil
	}
	return nil, last
}
