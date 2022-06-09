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
package dump

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	"k8s.io/client-go/dynamic"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
)

type DumpOptions struct {
	DynamicClient dynamic.Interface
	Mapper        meta.RESTMapper
	Result        *resource.Result

	PrintFlags  *genericclioptions.PrintFlags
	RecordFlags *genericclioptions.RecordFlags

	DryRun bool

	FilenameOptions resource.FilenameOptions
	Selector        string
	Raw             string

	Recorder genericclioptions.Recorder
	PrintObj func(obj kruntime.Object) error

	genericclioptions.IOStreams
}

var (
	createLong = templates.LongDesc(i18n.T(`
		Dump a pod using podname.
		`))

	createExample = templates.Examples(i18n.T(`
		# Dump a pod using podname.
		kubectl dump pod podname
		`))
)

func NewDumpOptions(ioStreams genericclioptions.IOStreams) *DumpOptions {
	return &DumpOptions{
		PrintFlags:  genericclioptions.NewPrintFlags("dumped").WithTypeSetter(scheme.Scheme),
		RecordFlags: genericclioptions.NewRecordFlags(),

		Recorder: genericclioptions.NoopRecorder{},

		IOStreams: ioStreams,
	}
}

func NewCmdDump(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewDumpOptions(ioStreams)

	cmd := &cobra.Command{
		Use: "dump -f FILENAME",
		DisableFlagsInUseLine: true,
		Short:   i18n.T("dump a pod realtime state to checkpoint."),
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, args, cmd))
			cmdutil.CheckErr(o.ValidateArgs(cmd, args))
			cmdutil.CheckErr(o.RunDump())
		},
	}

	// bind flag structs
	o.RecordFlags.AddFlags(cmd)

	usage := "to dump a checkpoint"
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	//cmd.MarkFlagRequired("filename")
	cmdutil.AddValidateFlags(cmd)
	cmdutil.AddApplyAnnotationFlags(cmd)
	cmd.Flags().StringVar(&o.Raw, "raw", o.Raw, "Raw URI to POST to the server.  Uses the transport specified by the kubeconfig file.")

	o.PrintFlags.AddFlags(cmd)
	return cmd
}

//验证参数是否合法
func (o *DumpOptions) ValidateArgs(cmd *cobra.Command, args []string) error {

	return nil
}

//完成参数的配置
func (o *DumpOptions) Complete(f cmdutil.Factory, args []string, cmd *cobra.Command) error {
	cmdNamespace, enforceNamespace, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	includeUninitialized := cmdutil.ShouldIncludeUninitialized(cmd, false)
	r := f.NewBuilder().
		Unstructured().
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, &o.FilenameOptions).
		IncludeUninitialized(includeUninitialized).
		ResourceTypeOrNameArgs(false, args...).RequireObject(false).
		Flatten().
		Do()
	err = r.Err()
	if err != nil {
		return err
	}
	o.Result = r

	o.Mapper, err = f.ToRESTMapper()
	if err != nil {
		return err
	}

	o.DynamicClient, err = f.DynamicClient()
	if err != nil {
		return err
	}
    
    printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	o.PrintObj = func(obj kruntime.Object) error {
		return printer.PrintObj(obj, o.Out)
	}

	return nil
}

func (o *DumpOptions) RunDump() error {
	return o.DumpResult(o.Result)
}

func (o *DumpOptions) DumpResult(r *resource.Result) error {
	// raw only makes sense for a single file resource multiple objects aren't likely to do what you want.
	// the validator enforces this, so
	err := r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		if err := o.Recorder.Record(info.Object); err != nil {
			glog.V(4).Infof("error recording current command: %v", err)
		}
		if _, err := DumpPod(info); err != nil {
			return cmdutil.AddSourceToErr("creating", info.Source, err)
		}

		return o.PrintObj(info.Object)
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *DumpOptions) raw(f cmdutil.Factory) error {
	restClient, err := f.RESTClient()
	if err != nil {
		return err
	}

	var data io.ReadCloser
	if o.FilenameOptions.Filenames[0] == "-" {
		data = os.Stdin
	} else {
		data, err = os.Open(o.FilenameOptions.Filenames[0])
		if err != nil {
			return err
		}
	}
	// TODO post content with stream.  Right now it ignores body content
	result := restClient.Post().RequestURI(o.Raw).Body(data).Do()
	if err := result.Error(); err != nil {
		return err
	}
	body, err := result.Raw()
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "%v", string(body))
	return nil
}

func DumpPod(info *resource.Info) (kruntime.Object, error) {
	options := &metav1.DumpOptions{}
	dumpResponse, err := resource.NewHelper(info.Client, info.Mapping).DumpResource(info.Namespace, info.Name, options)
	if err != nil {
		return dumpResponse, err
	}

	return dumpResponse, err
}
