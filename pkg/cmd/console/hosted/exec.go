// Copyright Contributors to the Open Cluster Management project
package hosted

import (
	"fmt"

	"github.com/stolostron/cm-cli/pkg/hypershift"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	switch len(args) {
	case 1:
		o.Hosting = "local-cluster"
		o.Hosted = args[0]
	case 2:
		o.Hosting = args[1]
		o.Hosted = args[1]
	}
	return nil
}

func (o *Options) validate() error {
	if o.Hosted == "" {
		return fmt.Errorf("<hosted-cluster-name> is missing")
	}
	if o.WithCredentials {
		return fmt.Errorf("not yet implemented, comming soon")
	}
	return nil
}

func (o *Options) run() (err error) {
	// dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	// if err != nil {
	// 	return err
	// }

	err = hypershift.OpenHosted(o.CMFlags, o.Hosted, o.Timeout, o.GetOptions.PrintFlags)
	if err != nil {
		return err
	}

	// if o.WithCredentials {
	// 	secret, err := dynamicClient.Resource(helpers.GvrHC).Namespace(o.Hosting).Get(context.TODO(), o.Hosted, metav1.GetOptions{})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Printf("%s", secret)
	// 	// return cph.PrintClusterClaimCred(cc, o.GetOptions.PrintFlags, o.WithCredentials)
	// }
	return nil

}
