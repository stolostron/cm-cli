// Copyright Contributors to the Open Cluster Management project
package acm

import (
	"fmt"
	"path/filepath"

	"github.com/stolostron/cm-cli/pkg/cmd/attach/cluster/scenario"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# install Advanced cluster management
%[1]s install acm --namespace <namespace> --channel <channel> [--manual-approval]
`

const (
	scenarioDirectory = "install"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")
var valuesDefaultPath = filepath.Join(scenarioDirectory, "values-default.yaml")

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cmd := &cobra.Command{
		Use:          "acm",
		Short:        "install acm",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			//TODO check if OCP installed
			// isSupported, err := helpers.IsSupported(o.CMFlags)
			// if err != nil {
			// 	return err
			// }
			// if !isSupported {
			// 	return fmt.Errorf("this command '%s %s' is only available on %s or %s",
			// 		helpers.GetExampleHeader(),
			// 		strings.Join(os.Args[1:], " "),
			// 		helpers.RHACM,
			// 		helpers.MCE)
			// }
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.complete(c, args); err != nil {
				return err
			}
			if err := o.validate(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.SetUsageTemplate(clusteradmhelpers.UsageTempate(cmd, scenario.GetScenarioResourcesReader(), valuesTemplatePath))
	cmd.Flags().StringVar(&o.channel, "channel", "", "The channel to use")
	cmd.Flags().StringVar(&o.namespace, "namespace", "open-cluster-management", "The namespace where to install ACM")
	cmd.Flags().StringVar(&o.operatorGroup, "operatorGroup", "open-cluster-management-group", "The operator group")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().BoolVar(&o.wait, "wait", false, "Wait until ACM installed is completed")
	cmd.Flags().BoolVar(&o.manualApproval, "manual-approval", false, "Set for manual approval otherwize automatic")
	cmd.Flags().IntVar(&o.timeout, "timeout", 180, "Timeout to get the klusterlet agent or addons ready in seconds")
	return cmd
}
