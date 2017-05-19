/*
Copyright 2017 The Kubernetes Authors.

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

package config

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/util/i18n"
	"fmt"
)


// renameContextOptions contains the assignable options from the args.
type renameContextOptions struct {
	configAccess clientcmd.ConfigAccess
	contextName  string
	newName      string
	out          io.Writer
}

const (
	renameContextUse = "rename-context CONTEXT_NAME NEW_NAME"

	renameContextShort = "Renames a context from the kubeconfig file."

	renameContextLong =
		`Renames a context from the kubeconfig file.

		CONTEXT_NAME is the context name that you wish change.

		NEW_NAME is the new name you wish to set.`

	renameContextExample =
		`# Rename the context 'old-name' to 'new-name' in your kubeconfig file
		kubectl config rename-context old-name new-name`
)


// NewCmdConfigRenameContext creates a command object for the "rename-context" action
func NewCmdConfigRenameContext(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &renameContextOptions{configAccess: configAccess}
	cmd := &cobra.Command{
		Use:     renameContextUse,
		Short:   i18n.T(renameContextShort),
		Long:    templates.LongDesc(renameContextLong),
		Example: templates.Examples(renameContextExample),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.validate())
			cmdutil.CheckErr(options.run())
			fmt.Fprintf(out, "Context %q was renamed to %q.\n", options.contextName, options.newName)
		},
	}
	return cmd
}

func (o *renameContextOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) != 2 {
		cmd.Help()
		return fmt.Errorf("Unexpected args: %v", args)
	}

	o.contextName = args[0]
	o.newName = args[1]
	return nil
}

func (o renameContextOptions) validate() error {
	if len(o.newName) == 0 {
		return errors.New("You must specify a new non-empty context name")
	}
	return nil
}

func (o renameContextOptions) run() error {
	config, err := o.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	configFile := o.configAccess.GetDefaultFilename()
	if o.configAccess.IsExplicitFile() {
		configFile = o.configAccess.GetExplicitFile()
	}

	context, exists := config.Contexts[o.contextName]
	if !exists {
		return fmt.Errorf("cannot rename the context %s, it's not in %s", o.contextName, configFile)
	}

	_, newExists := config.Contexts[o.newName]
	if newExists{
		return fmt.Errorf("cannot rename the context %s, the context %s already exists in %s",  o.contextName, o.newName, configFile)
	}

	config.Contexts[o.newName] = context
	delete(config.Contexts, o.contextName)

	if config.CurrentContext == o.contextName{
		config.CurrentContext = o.newName
	}

	return nil;
}