// Copyright Contributors to the Open Cluster Management project
package version

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/spf13/cobra"
	cmcli "github.com/stolostron/cm-cli"
	"github.com/stolostron/cm-cli/pkg/helpers"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	return nil
}

func (o *Options) validate() error {
	return nil
}
func (o *Options) run() (err error) {
	fmt.Printf("client\t\t\tversion\t:%s\n", cmcli.GetVersion())
	isSupported, err := helpers.IsSupported(o.CMFlags)
	if err != nil {
		return err
	}
	if isSupported {
		kubeClient, err := o.CMFlags.KubectlFactory.KubernetesClientSet()
		if err != nil {
			return err
		}
		dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
		if err != nil {
			return err
		}
		return o.runWithClient(kubeClient, dynamicClient)
	}
	return nil
}

func (o *Options) runWithClient(kubeClient kubernetes.Interface, dynamicClient dynamic.Interface) (err error) {
	var version, snapshot, server string
	switch {
	case helpers.IsRHACM(o.CMFlags):
		server = helpers.RHACM
		version, snapshot, err = helpers.GetACMVersion(o.CMFlags, kubeClient, dynamicClient)
	case helpers.IsMCE(o.CMFlags):
		server = helpers.MCE
		version, snapshot, err = helpers.GetMCEVersion(o.CMFlags, kubeClient, dynamicClient)
	}
	if version != "" {
		fmt.Printf("server %s release\tversion\t:%s\n", server, version)
	}
	if snapshot != "" {
		fmt.Printf("server %s image\ttag\t:%s\n", server, snapshot)
	}
	if err != nil {
		return err
	}
	return nil
}
