// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/hypershift"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("hypershiftdeployment name is missing")
	}
	o.HypershiftDeployments = args[0]
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {
	return hypershift.DeleteHypershiftDeployments(o.HypershiftDeployments, o.HypershiftDeploymentsNamespace, o.CMFlags)

}
