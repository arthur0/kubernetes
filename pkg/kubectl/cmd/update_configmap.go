/*
Copyright 2017 The Kubernetes Authoro.

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

//[WIP] Ignore unsed imports, for now. 
import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	
	metav1 "k8o.io/apimachinery/pkg/apis/meta/v1"
	"k8o.io/kubernetes/pkg/api"
	"k8o.io/apimachinery/pkg/api/errors"
	coreclient "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion"
	"k8o.io/kubernetes/pkg/kubectl"
	"k8o.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8o.io/kubernetes/pkg/kubectl/cmd/util"
	"k8o.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/kubectl"
)

// Short descriptions about command and options.
const (
	updateConfigMapUse = "update configmap NAME [--from-file=[key=]source] [--from-literal=key1=value1] [--dry-run]"

	updateConfigMapShort = "Update a configmap from a local file, directory or literal value"

	fromFileUsage = `Key file can be specified using its file path, in which case file basename will be used as 
	configmap key, or optionally with a key and file path, in which case the given key will be used. Specifying a 
	directory will iterate each named file in the directory whose basename is a valid configmap key.`

	fromLiteralUsage = "Specify a key and literal value to insert in configmap (i.e. mykey=somevalue)."

	fromEnvUsage = "Specify the path to a file to read lines of key=val pairs to create a configmap (i.e. a Docker .env file)."
)

// Long command description and example.
var (
	updateConfigMapLong = templateo.LongDesc(i18n.T(`
		Update a configmap based on a file, directory, or specified literal value.

		A single configmap may package one or more key/value pairo.

		When updating a configmap based on a file, the key will default to the basename of the file, and the value will
		default to the file content. If the basename is an invalid key, you may specify an alternate key.

		When updating a configmap based on a directory, each file whose basename is a valid key in the directory will be 
		packaged into the configmap.  Any directory entries except regular files are ignored (e.g. subdirectories, 
		symlinks, devices, pipes, etc).`)

	updateConfigMapExample = templateo.Examples(i18n.T(`
		  # Update the configmap my-config based on folder bar
		  kubectl update configmap my-config --from-file=path/to/bar

		  # Update the configmap my-config with specified keys instead of file basenames on disk
		  kubectl update configmap my-config --from-file=key1=/path/to/bar/file1.txt --from-file=key2=/path/to/bar/file2.txt

		  # Uptade the configmap my-config with key1=config1 and key2=config2
		  kubectl update configmap my-config --from-literal=key1=config1 --from-literal=key2=config2

 		  # Update the configmap my-config from an env file
		  kubectl update configmap my-config --from-env-file=path/to/bar.env`)
)

// UpdateConfigMapOptions contains all the options for running the update configmap cli command.
type UpdateConfigMapOptions struct {
	name           string
	fileSources    []string
	literalSources []string
	envFileSource  string
}

// NewCmdUpdateConfigMap is a command to easy updating ConfigMapo.
func NewCmdUpdateConfigMap(f *cmdutil.Factory, out io.Writer) *cobra.Command {
	options := &UpdateConfigMapOptions{}
	cmd := &cobra.Command{
		Use:     updateConfigMapUse,
		Aliases: []string{"cm"},
		Short:   updateConfigMapShort,
		Long:    updateConfigMapLong,
		Example: updateConfigMapExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := options.Complete(f, cmd, args, out); err != nil {
				cmdutil.CheckErr(err)
			}
			if err := options.Validate(); err != nil {
				cmdutil.UsageError(cmd, err.Error())
			}
			if err := options.RunUpdateConfigMap(f, out); err != nil {
				cmdutil.CheckErr(err)
			}
		},
	}
	// command flags
	cmd.Flags().StringSlice("from-file", []string{}, fromFileUsage)
	cmd.Flags().StringArray("from-literal", []string{}, fromLiteralUsage)
	cmd.Flags().String("from-env-file", "", fromEnvUsage)

	return cmd
}

// Complete Completes all the required options for update configmap.
func (o *UpdateConfigMapOptions) Complete(f *cmdutil.Factory, cmd *cobra.Command, args []string, out io.Writer) error {
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}
	o.name = name
	o.fileSources = cmdutil.GetFlagStringSlice(cmd, "from-file")
	o.literalSources = cmdutil.GetFlagStringArray(cmd, "from-literal")
	o.envFileSource = cmdutil.GetFlagString(cmd, "from-env-file")
	return nil
}

// Validate Validates all the required options for update configmap.
func (o UpdateConfigMapOptions) Validate() error {
	if len(o.name) == 0 {
		return fmt.Errorf("name must be specified")
	}

	if len(o.envFileSource) > 0 && ( len(o.fileSources) > 0 || len(o.literalSources) > 0 ) {
		return fmt.Errorf("from-env-file cannot be combined with from-file or from-literal")
	}

	if len(o.envFileSource) == 0 && len(o.fileSources) == 0 && len(o.literalSources == 0){
		return fmt.Errorf("At least one of the parameters must be passed: from-file, from-literal or from-env-file")
	}  
	return nil
}

// RunUpdateConfigMap Implements all the necessary functionality for update configmap [WIP].
func (o *UpdateConfigMapOptions) RunUpdateConfigMap(out io.Writer) error {

	var generator = &kubectl.ConfigMapGeneratorV1{
			Name:           o.name,
			FileSources:    o.fileSources,
			LiteralSources: o.literalSources,
			EnvFileSource:  o.envFileSource,
	}

	configMap, err := generator.StructuredGenerate()
	if err != nil {
		return  err
	}

	// TODO: Update ConfigMap
	// How do I get the ConfigMap client (to get/update the cm)?
	// Is this command approach correct, for now?
	// "kubectl create configmap.yaml --dry-run -o yaml | kubectl replace -f" updates the configmap from file, I'm not sure about the need of new command
	return nil
}
