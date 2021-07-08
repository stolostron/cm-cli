// Copyright Contributors to the Open Cluster Management project
package use

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clustername is missing")
	}
	o.Cluster.Name = args[0]
	return nil
}

func (o *Options) validate() error {
	kubeClient, err := o.CMFlags.KubectlFactory.KubernetesClientSet()
	if err != nil {
		return err
	}
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	return o.validateWithClient(kubeClient, dynamicClient)
}

func (o *Options) validateWithClient(kubeClient kubernetes.Interface, dynamicClient dynamic.Interface) error {
	return nil
}

func (o *Options) run() (err error) {
	cph, err := clusterpoolhost.GetClusterPoolHost(o.Cluster.Name)
	switch {
	case clusterpoolhost.IsNotFound(err):
		err = clusterpoolhost.VerifyContext(o.Cluster.Name, o.CMFlags.DryRun, o.outputFile)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		return cph.VerifyContext(o.CMFlags.DryRun, o.outputFile)
	}
	return nil
}
