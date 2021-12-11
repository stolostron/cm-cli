// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"
	"path/filepath"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost/scenario"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	scenarioDirectory = "create/clusterpool"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "common/values-template.yaml")

var example = `
# Create a clusterpool
%[1]s create cp --values values.yaml [--cph <clusterpoolhost_name>]

# Create a cluster with cluster name overwrite by args
%[1]s create cp [<clusterpool_name>] --values values.yaml [--cph <clusterpoolhost_name>]
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "clusterpool",
		Aliases:      []string{"cp"},
		Short:        "Create a clusterpool",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
	}

	cmd.SetUsageTemplate(clusteradmhelpers.UsageTempate(cmd, scenario.GetScenarioResourcesReader(), valuesTemplatePath))
	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().StringVar(&o.valuesPath, "values", "", "The files containing the values")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().StringVar(&o.clusterSetName, "cluster-set", "", "The clusterset to which the clusterpool should be place")
	return cmd
}
