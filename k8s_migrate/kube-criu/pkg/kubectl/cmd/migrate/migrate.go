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

package migrate

import (
	"fmt"
	"os"
	"io/ioutil"
	"os/exec"
	"time"
	"gopkg.in/yaml.v2"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"

)

const (
	PodSucceeded     = "Succeeded"
	PodHalfFailed    = "HalfFailed"
	PodFailed        = "Failed"
)

type MigrateOptions struct {

	DynamicClient dynamic.Interface
	Mapper        meta.RESTMapper
	Result        *resource.Result

	FilenameOptions resource.FilenameOptions
	Node            string

	genericclioptions.IOStreams
}

var (
	migrateLong = templates.LongDesc(i18n.T(`
		Migrate a pod to a new node.`))

	migrateExample = templates.Examples(i18n.T(`
		# Migrate a pod using the data in pod.json.
		kubectl migrate -f pod.yaml --node=nodelables`))
)

func NewMigrateOptions(ioStreams genericclioptions.IOStreams) *MigrateOptions {
	return &MigrateOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdMigrate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewMigrateOptions(ioStreams)

	cmd := &cobra.Command{
		Use: "migrate -f FILENAME --storage filePath --node nodeLable",
		DisableFlagsInUseLine: true,
		Short:   i18n.T("migrate a pod to a new node."),
		Long:    migrateLong,
		Example: migrateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.ValidateArgs(cmd, args))
			cmdutil.CheckErr(o.RunMigrate(f, cmd, args))
		},
	}

	// bind flag structs
	usage := "to migrate a pod"
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	//cmd.MarkFlagRequired("filename")
	cmdutil.AddValidateFlags(cmd)
	cmdutil.AddApplyAnnotationFlags(cmd)
	cmd.Flags().StringVarP(&o.Node, "node", "l", o.Node, "The node to which the pod is to be migrated, uses label for node(e.g. -n NODENAME)")
	return cmd
}

func (o *MigrateOptions) ValidateArgs(cmd *cobra.Command, args []string) error {
	if len(o.FilenameOptions.Filenames) < 1 {
		return cmdutil.UsageErrorf(cmd, "File for podcheckpoint is null!!")
	}
	return nil
}

func (o *MigrateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	var err error

	o.Mapper, err = f.ToRESTMapper()
	if err != nil {
		return err
	}
	o.DynamicClient, err = f.DynamicClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *MigrateOptions) RunMigrate(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var (
		err     error
		conf 	map[string]interface{}
		flag    bool
	)
 
	home := os.Getenv("HOME")
	filenames := o.FilenameOptions.Filenames
	for _, s := range filenames {
		getConf(&conf, s)
	}

	metadata := conf["metadata"].(map[interface{}]interface{})              
    podcheckpointName := metadata["name"].(string)
        
    spec := conf["spec"].(map[interface{}]interface{})
    podName := spec["podName"].(string)

	err = o.RunCreateCheckpoint(f, cmd)
	if err != nil {
		glog.Infof(err.Error())
		return err
	}

	filename := home + "/" + "podcheckpoint-" + podcheckpointName + ".yaml" 
	var iter int
	iter = 0
	for {
		status, err := getPodCheckpointStatus(filename, podcheckpointName) 
		iter++
		if err == nil {
			flag = false
			switch status {
			case PodSucceeded:
				glog.Infof("Create Checkpoint Succeeded!!!")
				flag = true
				break
			case PodFailed:
				glog.Infof("Create Checkpoint Failed!!!")
				return nil
			case PodHalfFailed:
				glog.Infof("Create Checkpoint HalfFailed!!!")
				return nil
			default:
				if iter == 1{
					glog.Infof("Waitting for Create Checkpoint....")
				}		
				flag = false
				time.Sleep(time.Millisecond * 2)
			}
			if flag == true {
				break
			}
		} else {
			glog.Infof("getPodCheckpointStatus Failed!!!")
		}	
	}

	err = clear(filename)	
	err = o.RunCreateNewPod(home, podcheckpointName, podName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (o *MigrateOptions) RunCreateCheckpoint(f cmdutil.Factory, cmd *cobra.Command) error {
	createCmd := "kubectl create -f " + o.FilenameOptions.Filenames[0]
	c := exec.Command("bash", "-c", createCmd)
	if err := c.Run(); err != nil {
		glog.Infof("Error: ", err)
		return err
	}
	return nil
}

func (o *MigrateOptions) RunCreateNewPod(home string, podcheckpointName string, podName string) error {
	var err error
    inFilepath := home + "/" + podName + "-" + "podcheckpoint" + ".yaml"
	outFilepath := home + "/" + podName + "-" + "podcheckpoint-new" + ".yaml"

    getConfCmd := "kubectl get pod " + podName + " -o yaml > " + inFilepath
    c := exec.Command("bash", "-c", getConfCmd)
    if err = c.Run(); err != nil {
        return err
    }

    deleteCmd := "kubectl delete pod " + podName
	c = exec.Command("bash", "-c", deleteCmd)
	if err = c.Run(); err != nil {
		return err
	}

	HandleYamlFile(inFilepath, outFilepath, o.Node, podcheckpointName)

	createCmd := "kubectl create -f " + outFilepath
	c = exec.Command("bash", "-c", createCmd)
	if err = c.Run(); err != nil {
		return err
	}

	err = clear(inFilepath)
	err = clear(outFilepath)
	return nil
}

func getConf(conf *map[string]interface{}, filepath string) *map[string]interface{} {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		glog.Infof(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		glog.Infof(err.Error())
	}
	return conf
}

func HandleYamlFile(inFile string, outFile string, nodeName string, podcheckpointName string) error {
    // Read buffer from jsonFile
    byteValue, err := ioutil.ReadFile(inFile)
    if err != nil {
        return err
    }
    // We have known the outer json object is a map, so we define  result as map.
    // otherwise, result could be defined as slice if outer is an array
    var result map[string]interface{}
    err = yaml.Unmarshal(byteValue, &result)
    if err != nil {
        return err
    }
    // handle peers
    spec := result["spec"].(map[interface{}]interface{})
    spec["nodeName"] = nodeName

    metadata := result["metadata"].(map[interface{}]interface{})
    annots := metadata["annotations"]
    if annots == nil {
    	var m map[string]string
    	m = make(map[string]string)
    	m["podCheckpoint"] = podcheckpointName
    	metadata["annotations"] = m
    } else {
    	annot := annots.(map[interface{}]interface{})
    	annot["podCheckpoint"] = podcheckpointName
    }


    // Convert golang object back to byte
    byteValue, err = yaml.Marshal(result)
    if err != nil {
        return err
    }
    // Write back to file
    err = ioutil.WriteFile(outFile, byteValue, 0644)
    if err != nil{
        return err
    }
    return err
}

func getPodCheckpointStatus(filename string, podcheckpointName string) (string, error) {
	getConfCmd := "kubectl get podcheckpoint " + podcheckpointName + " -o yaml > " + filename
    c := exec.Command("bash", "-c", getConfCmd)
    if err := c.Run(); err != nil {
        return "", err
    }

    var conf map[string]interface{}
    getConf(&conf, filename)
    status := conf["status"].(map[interface{}]interface{})
    podcheckpointStatus := status["phase"].(string)

    return podcheckpointStatus, nil
}

func clear(filepath string) error {
	c := exec.Command("rm", "-rf", filepath)
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}