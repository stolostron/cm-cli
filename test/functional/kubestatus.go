// Copyright Contributors to the Open Cluster Management project

package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {
	root := &cobra.Command{
		Use:   "status",
		Short: "CLI for Red Hat Advanced Cluster Management",
		//This remove the auto-generated tag in the cobra doc
		DisableAutoGenTag: true,
	}

	flags := root.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

	matchVersionKubeConfigFlags.AddFlags(flags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	// From this point and forward we get warnings on flags that contain "_" separators
	root.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	setStatusCmd := NewSetStatusCmd(f, streams)
	root.AddCommand(setStatusCmd)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}

}

type setStatusCmd struct {
	KubectlFactory cmdutil.Factory
	streams        genericclioptions.IOStreams
	group          string
	version        string
	resource       string
	name           string
	namespace      string
	status         string
}

func NewSetStatusCmd(kubectlFactory cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &setStatusCmd{
		KubectlFactory: kubectlFactory,
		streams:        streams,
	}
	cmd := &cobra.Command{
		Use: "update",

		RunE: func(c *cobra.Command, args []string) error {
			o.name = args[0]
			return o.run()
		},
	}

	cmd.Flags().StringVar(&o.group, "group", "", "the group of the resource")
	cmd.Flags().StringVar(&o.resource, "resource", "", "the name of the resource")
	cmd.Flags().StringVar(&o.version, "version", "", "the version of the resource")
	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", "", "the version of the resource")
	cmd.Flags().StringVar(&o.status, "status", "", "the content status json")

	return cmd
}

func (ss *setStatusCmd) run() error {

	d, err := ss.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	gvr := schema.GroupVersionResource{Group: ss.group, Version: ss.version, Resource: ss.resource}
	obj, err := d.Resource(gvr).Namespace(ss.namespace).Get(context.TODO(), ss.name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	status := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(ss.status), &status)
	if err != nil {
		return err
	}

	obj.Object["status"] = status

	if _, err := d.Resource(gvr).Namespace(ss.namespace).UpdateStatus(context.TODO(), obj, metav1.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}
